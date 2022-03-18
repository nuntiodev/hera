package handler

import (
	"context"
	"fmt"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/repository/user_repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

const (
	mongoInsert = "insert"
	mongoUpdate = "update"
	mongoDelete = "delete"
)

type ChangeID struct {
	Data string `bson:"_data"`
}

type documentKey struct {
	ID string `bson:"_id"`
}

type namespace struct {
	Db   string `bson:"db"`
	Coll string `bson:"coll"`
}

type UpdateDescription struct {
	UpdatedFields user_repository.User `bson:"updatedFields"`
}

type ChangeEvent struct {
	MongoID           ChangeID             `bson:"_id"`
	ClusterTime       primitive.Timestamp  `bson:"clusterTime"`
	OperationType     string               `bson:"operationType"`
	FullDocument      user_repository.User `bson:"fullDocument"`
	UpdateDescription UpdateDescription    `bson:"updateDescription"`
	DocumentKey       documentKey          `bson:"documentKey"`
	Ns                namespace            `bson:"ns"`
}

type StreamConn struct {
	ProjectId  string
	Connection *mongo.ChangeStream
	UsedAt     time.Time
}

var (
	mu           sync.Mutex
	connections  = map[string]*StreamConn{}
	lastTimeUsed = time.Minute * 5
)

func cleanupConnections() {
	ctx := context.Background()
	for k, v := range connections {
		if v != nil && time.Now().Sub(v.UsedAt) > lastTimeUsed {
			removeConnection(ctx, k)
		}
	}
}

func addStream(sessionId string, stream *mongo.ChangeStream) {
	if stream != nil && sessionId != "" {
		mu.Lock()
		defer mu.Unlock()
		connections[sessionId] = &StreamConn{
			Connection: stream,
			UsedAt:     time.Now(),
		}
	}
}

func removeConnection(ctx context.Context, sessionID string) {
	if sessionID == "" {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if val, ok := connections[sessionID]; ok {
		if val != nil {
			if val.Connection != nil {
				val.Connection.Close(ctx)
			}
		}
		delete(connections, sessionID)
	}
}

func (h *defaultHandler) handleStream(ctx context.Context, stream *mongo.ChangeStream, server go_block.UserService_GetStreamServer, encryptionKey, sessionId string) error {
	defer removeConnection(ctx, sessionId)
	g := new(errgroup.Group)
	for stream.Next(ctx) {
		var changeEvent ChangeEvent
		var streamType go_block.StreamType
		if err := stream.Decode(&changeEvent); err != nil {
			return err
		}
		userResp := &go_block.User{}
		switch changeEvent.OperationType {
		case mongoInsert:
			userResp = user_repository.UserToProtoUser(&changeEvent.FullDocument)
			streamType = go_block.StreamType_CREATE
		case mongoUpdate:
			userResp = user_repository.UserToProtoUser(&changeEvent.UpdateDescription.UpdatedFields)
			streamType = go_block.StreamType_UPDATE
		case mongoDelete:
			streamType = go_block.StreamType_DELETE
		}
		userResp.Id = changeEvent.DocumentKey.ID
		if encryptionKey != "" && userResp.EncryptedAt.IsValid() {
			if err := h.crypto.DecryptUser(encryptionKey, userResp); err != nil {
				h.zapLog.Debug(err.Error())
				return err
			}
		}
		fmt.Println(encryptionKey, userResp.EncryptedAt.IsValid())
		streamResp := &go_block.UserStream{
			StreamType: streamType,
			User:       userResp,
		}
		h.zapLog.Debug(fmt.Sprintf("streaming new user info: %s", streamResp.String()))
		connections[sessionId].UsedAt = time.Now()
		g.Go(func() error {
			if err := server.Send(streamResp); err != nil {
				h.zapLog.Debug(err.Error())
				stream.Close(ctx)
				return err
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func (h *defaultHandler) GetStream(req *go_block.UserRequest, server go_block.UserService_GetStreamServer) error {
	ctx, _ := context.WithCancel(context.Background())
	users, err := h.repository.Users(context.Background(), req.Namespace)
	if err != nil {
		return err
	}
	// remove old connection if exist
	removeConnection(ctx, req.SessionId)
	stream, err := users.GetStream(context.Background(), req.User)
	if err != nil {
		return err
	}
	// add new connection
	addStream(req.SessionId, stream)
	h.zapLog.Debug(fmt.Sprintf("adding new connection with a total count of: %d", len(connections)))
	// stream
	return h.handleStream(ctx, stream, server, req.EncryptionKey, req.SessionId)
}
