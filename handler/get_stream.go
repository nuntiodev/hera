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
	"strings"
	"sync"
	"time"
)

const (
	mongoInsert = "insert"
	mongoUpdate = "update"
	mongoDelete = "delete"
)

var (
	maxStreams           = 100
	maxStreamConnections = 5
	clientConnections    = map[string]int{}
	sessionConnections   = map[string]*Connection{}
	mu                   sync.Mutex
)

type Connection struct {
	changeStream *mongo.ChangeStream
	cancel       context.CancelFunc
	usedAt       time.Time
	namespace    string
}

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

// todo: automatically delete connections
func (h *defaultHandler) handleStream(ctx context.Context, stream *mongo.ChangeStream, req *go_block.UserRequest, server go_block.UserService_GetStreamServer) error {
	defer h.removeConnection(context.Background(), req.SessionId, req.Namespace)
	var g errgroup.Group
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
		h.zapLog.Debug(fmt.Sprintf("streaming new user: %s", streamType.String()))
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
	return nil
}

func (h *defaultHandler) GetStream(req *go_block.UserRequest, server go_block.UserService_GetStreamServer) error {
	h.zapLog.Debug("initializing stream")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := h.validateMaxStreams(ctx, req.SessionId, req.Namespace); err != nil {
		return err
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
	mu.Lock()
	sessionConnections[req.SessionId] = &Connection{
		namespace:    req.Namespace,
		changeStream: stream,
		cancel:       cancel,
		usedAt:       time.Now(),
	}
	mu.Unlock()
	// stream
	return h.handleStream(ctx, stream, req, server)
}

func (h *defaultHandler) cleanupConnections() {
	mu.Lock()
	defer mu.Unlock()
	for k, v := range sessionConnections {
		if time.Now().Sub(v.usedAt) < time.Second*60 {
			h.removeConnection(context.Background(), k, v.namespace)
		}
	}
}

func (h *defaultHandler) removeConnection(ctx context.Context, sessionId, ns string) {
	h.zapLog.Debug("closing connection...")
	mu.Lock()
	defer mu.Unlock()
	sessionId = strings.TrimSpace(sessionId)
	ns = strings.TrimSpace(ns)
	if val, ok := sessionConnections[sessionId]; ok && sessionId != "" && val != nil && val.changeStream != nil {
		val.cancel()
		if err := val.changeStream.Close(ctx); err != nil {
			h.zapLog.Debug(err.Error())
		}
		delete(sessionConnections, sessionId)
	}
	if val, ok := clientConnections[ns]; ok && ns != "" {
		if val > 1 {
			clientConnections[ns] -= 1
		} else {
			delete(clientConnections, ns)
		}
	}
}

func (h *defaultHandler) validateMaxStreams(ctx context.Context, sessionId, ns string) error {
	mu.Lock()
	defer mu.Unlock()
	sessionId = strings.TrimSpace(sessionId)
	ns = strings.TrimSpace(ns)
	// close previous connection if present
	if val, ok := sessionConnections[sessionId]; ok && sessionId != "" && val != nil && val.changeStream != nil {
		if err := val.changeStream.Close(ctx); err != nil {
			h.zapLog.Debug(err.Error())
		}
		if val, ok := clientConnections[ns]; ok {
			if val > 0 {
				clientConnections[ns] = val - 1
			} else {
				delete(clientConnections, ns)
			}
		}
		delete(sessionConnections, sessionId)
	}
	// check if max connections is reached
	totalConnections := 0
	for k, v := range clientConnections {
		h.zapLog.Debug(fmt.Sprintf("k: %s, v: %d", k, v))
		totalConnections += v
	}
	h.zapLog.Debug(fmt.Sprintf("total project connections are: %d", len(clientConnections)))
	if len(clientConnections) >= maxStreamConnections {
		return errors.New(fmt.Sprintf("Max stream connections per client reached %d. Streams are expensive so remember to clean the up properly.", maxStreamConnections))
	}
	h.zapLog.Debug(fmt.Sprintf("total connections are: %d", totalConnections))
	if totalConnections >= maxStreamConnections {
		return errors.New("max stream connections for server")
	}
	// measure client connections
	if ns != "" {
		// only allow single connection per session if namespace is set
		if val, ok := clientConnections[ns]; ok {
			clientConnections[ns] = val + 1
		} else {
			clientConnections[ns] = 1
		}
	}
	return nil
}
