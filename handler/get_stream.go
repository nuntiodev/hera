package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/repository/user_repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

const (
	mongoInsert = "insert"
	mongoUpdate = "update"
	mongoDelete = "delete"
)

var (
	maxConnectionsPerClient = 15
	maxStreamAge            = time.Minute * 4
	clientConnections       = map[string]int{}
	sessionConnections      = map[string]*mongo.ChangeStream{}
	mu                      sync.Mutex
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

func (h *defaultHandler) handleStream(ctx context.Context, stream *mongo.ChangeStream, req *go_block.UserRequest, server go_block.UserService_GetStreamServer) error {
	lastUsedAt := time.Now()
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
		streamResp := &go_block.UserStream{
			StreamType: streamType,
			User:       userResp,
		}
		h.zapLog.Debug(fmt.Sprintf("streaming new user info"))
		if err := server.Send(streamResp); err != nil {
			return err
		}
		if time.Now().Sub(lastUsedAt) > maxStreamAge {
			h.zapLog.Debug("breaking out of stream due to inactivity")
			break
		}
	}
	if err := stream.Err(); err != nil {
		h.zapLog.Debug(err.Error())
		return err
	}
	return nil
}

func (h *defaultHandler) GetStream(req *go_block.UserRequest, server go_block.UserService_GetStreamServer) error {
	h.zapLog.Debug("initializing stream")
	ctx, cancel := context.WithTimeout(context.Background(), maxStreamAge+time.Second*5)
	defer cancel()
	// add connection to client list
	if len(clientConnections) > maxConnectionsPerClient {
		return errors.New(fmt.Sprintf("Max stream connections per client reached %d. Streams are expensive so remember to clean the up properly.", maxConnectionsPerClient))
	}
	// only allow single connection per session
	clientConnections[req.Namespace] += 1
	if val, ok := sessionConnections[req.SessionId]; ok && req.SessionId != "" {
		if err := val.Close(ctx); err != nil {
			h.zapLog.Debug(err.Error())
		}
	}
	// create stream
	users, err := h.repository.Users(ctx, req.Namespace)
	if err != nil {
		return err
	}
	stream, err := users.GetStream(ctx, req.User)
	if err != nil {
		return err
	}
	// remember to close connections
	defer func() {
		if err := stream.Close(ctx); err != nil {
			h.zapLog.Debug(err.Error())
		}
	}()
	// stream
	return h.handleStream(ctx, stream, req, server)
}
