package handler

import (
	"context"
	"errors"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/repository/user_repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	mongoInsert = "insert"
	mongoUpdate = "update"
	mongoDelete = "delete"
	refresh     = time.Second * 3
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

func validateStreamType(streamType []go_block.StreamType) (create bool, update bool, delete bool) {
	create = false
	update = false
	delete = false
	for _, st := range streamType {
		switch st {
		case go_block.StreamType_CREATE:
			create = true
		case go_block.StreamType_UPDATE:
			update = true
		case go_block.StreamType_DELETE:
			delete = true
		}
	}
	return create, update, delete
}

func (h *defaultHandler) handleUsersStream(ctx context.Context, filter *go_block.UserFilter, encryptionKey string, types []go_block.StreamType, namespace string, userBatch []*go_block.User, server go_block.UserService_GetStreamServer, autoFollowStream bool) error {
	stream, err := h.repository.UserRepository.GetUsersStream(ctx, namespace, userBatch)
	if err != nil {
		return err
	}
	defer stream.Close(ctx)
	// config
	_, updateStream, deleteStream := validateStreamType(types)
	autoEstablish := autoFollowStream && len(userBatch) > 0 && (updateStream || deleteStream)
	currentWatchedUsers := userBatch
	var newWatchedUsers []*go_block.User
	for stream.Next(context.Background()) {
		// handle event
		var changeEvent ChangeEvent
		var streamType go_block.StreamType
		if err := stream.Decode(&changeEvent); err != nil {
			return err
		}
		if changeEvent.OperationType == mongoDelete {
			streamType = go_block.StreamType_DELETE
			// handle delete event internally
			if autoEstablish {

				// delete from currently watched users
				for index, user := range currentWatchedUsers {
					if user.Id == changeEvent.DocumentKey.ID {
						currentWatchedUsers = append(currentWatchedUsers[:index], currentWatchedUsers[index+1:]...)
					}
				}
				// delete from new watched users
				for index, user := range newWatchedUsers {
					if user.Id == changeEvent.DocumentKey.ID {
						currentWatchedUsers = append(currentWatchedUsers[:index], currentWatchedUsers[index+1:]...)
					}
				}
			}
		} else if changeEvent.OperationType == mongoInsert {
			streamType = go_block.StreamType_CREATE
			// handle create event internally
			if autoEstablish && filter.Sort == go_block.UserFilter_CREATED_AT {
				newWatchedUsers = append(newWatchedUsers, user_repository.UserToProtoUser(&changeEvent.FullDocument))
			}
		} else if changeEvent.OperationType == mongoUpdate {
			// handle update event internally
			streamType = go_block.StreamType_UPDATE
			// handle create event internally
			if autoEstablish && filter.Sort == go_block.UserFilter_UPDATE_AT {
				newWatchedUsers = append(newWatchedUsers, user_repository.UserToProtoUser(&changeEvent.UpdateDescription.UpdatedFields))
			}
		}
		// return event to client
		userResp := &go_block.User{}
		switch streamType {
		case go_block.StreamType_CREATE:
			userResp = user_repository.UserToProtoUser(&changeEvent.FullDocument)
		case go_block.StreamType_UPDATE:
			userResp = user_repository.UserToProtoUser(&changeEvent.UpdateDescription.UpdatedFields)
		case go_block.StreamType_DELETE:
			userResp = user_repository.UserToProtoUser(&changeEvent.FullDocument)
		}
		if encryptionKey != "" && userResp.EncryptedAt.IsValid() {
			h.crypto.DecryptUser(encryptionKey, userResp)
		}
		go server.Send(&go_block.UserStream{
			StreamType: streamType,
			User:       userResp,
		})
		// check if changes has happened and auto establish a new connection
		if autoEstablish && len(currentWatchedUsers) == 0 {
			// rebase
			users, err := h.repository.UserRepository.GetAll(ctx, filter, namespace, encryptionKey)
			if err != nil {
				return err
			}
			h.zapLog.Info("rebasing stream with new users")
			if err := stream.Close(ctx); err != nil {
				return err
			}
			return h.handleUsersStream(ctx, filter, encryptionKey, types, namespace, users, server, autoFollowStream)
		} else if autoEstablish && len(newWatchedUsers) > 0 {
			// restart
			var users []*go_block.User
			var iterations int
			if len(newWatchedUsers) > len(userBatch) {
				iterations = len(userBatch)
			} else {
				iterations = len(newWatchedUsers)
			}
			index := 0
			for i := 0; i < len(userBatch); i++ {
				if i < iterations {
					// take the newest first (which is the last in the array)
					users = append(users, newWatchedUsers[len(newWatchedUsers)-1-i])
					continue
				}
				// take the newest elements (which is the first)
				users = append(users, currentWatchedUsers[index])
				index++
			}
			h.zapLog.Info("restarting stream with new users")
			if err := stream.Close(ctx); err != nil {
				return err
			}
			return h.handleUsersStream(ctx, filter, encryptionKey, types, namespace, users, server, autoFollowStream)
		}
	}
	return nil
}

func (h *defaultHandler) GetStream(req *go_block.UserRequest, server go_block.UserService_GetStreamServer) error {
	_, updateStream, deleteStream := validateStreamType(req.StreamType)
	if (updateStream || deleteStream) && len(req.UserBatch) == 0 {
		return errors.New("Invalid stream request. If you want yo open a stream and listen for updates/deletions you need to pass a user batch which contains the users with their ids to you want to listen to.")
	}
	if req.Filter == nil {
		req.Filter = &go_block.UserFilter{
			Sort: go_block.UserFilter_CREATED_AT,
		}
	}
	return h.handleUsersStream(context.Background(), req.Filter, req.EncryptionKey, req.StreamType, req.Namespace, req.UserBatch, server, req.AutoFollowStream)
}
