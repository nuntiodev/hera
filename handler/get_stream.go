package handler

import (
	"context"
	"errors"
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
	Server     go_block.UserService_GetStreamServer
	UsedAt     time.Time
}

var (
	mu          sync.Mutex
	connections = map[string]*StreamConn{}
)

func getStream(sessionId string) (*StreamConn, error) {
	if sessionId != "" {
		mu.Lock()
		defer mu.Unlock()
		val, ok := connections[sessionId]
		if ok && val != nil && val.Connection != nil && val.Server != nil {
			return val, nil
		}
	}
	return nil, errors.New("no stream with that session id")
}

func addStream(sessionId string, stream *mongo.ChangeStream, server go_block.UserService_GetStreamServer) {
	if stream != nil && sessionId != "" {
		mu.Lock()
		defer mu.Unlock()
		connections[sessionId] = &StreamConn{
			Connection: stream,
			UsedAt:     time.Now(),
			Server:     server,
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

func (h *defaultHandler) handleStream(ctx context.Context, stream *mongo.ChangeStream, req *go_block.UserRequest, server go_block.UserService_GetStreamServer) error {
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
		if req.EncryptionKey != "" && userResp.EncryptedAt.IsValid() {
			if err := h.crypto.DecryptUser(req.EncryptionKey, userResp); err != nil {
				return err
			}
		}
		fmt.Println("ENC KEY")
		fmt.Println(req.EncryptionKey, userResp.EncryptedAt.IsValid())
		fmt.Println("ENC KEY")
		streamResp := &go_block.UserStream{
			StreamType: streamType,
			User:       userResp,
		}
		h.zapLog.Debug(fmt.Sprintf("streaming new user info: %s", streamResp.String()))
		mu.Lock()
		connections[req.SessionId].UsedAt = time.Now()
		mu.Unlock()
		g.Go(func() error {
			if err := server.Send(streamResp); err != nil {
				return err
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	if err := stream.Err(); err != nil {
		h.zapLog.Debug(err.Error())
		return err
	}
	// setup new mongo stream
	h.zapLog.Debug("returning from broken stream...")
	users, err := h.repository.Users(context.Background(), req.Namespace)
	if err != nil {
		return err
	}
	newMongoStream, err := users.GetStream(ctx, req.User)
	if err != nil {
		return err
	}
	return h.handleStream(ctx, newMongoStream, req, server)
}

func (h *defaultHandler) GetStream(req *go_block.UserRequest, server go_block.UserService_GetStreamServer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if conn, err := getStream(req.SessionId); err == nil {
		h.zapLog.Debug("stream is already running")
		return h.handleStream(ctx, conn.Connection, req, conn.Server)
	}
	h.zapLog.Debug("initializing stream")
	users, err := h.repository.Users(context.Background(), req.Namespace)
	if err != nil {
		return err
	}
	stream, err := users.GetStream(context.Background(), req.User)
	if err != nil {
		return err
	}
	// add new connection
	addStream(req.SessionId, stream, server)
	defer removeConnection(ctx, req.SessionId)
	h.zapLog.Debug(fmt.Sprintf("adding new connection with a total count of: %d", len(connections)))
	// stream
	return h.handleStream(ctx, stream, req, server)
}
