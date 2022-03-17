package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/repository/user_repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (h *defaultHandler) GetStream(req *go_block.UserRequest, server go_block.UserService_GetStreamServer) error {
	users, err := h.repository.Users(context.Background(), req.Namespace)
	if err != nil {
		return err
	}
	stream, err := users.GetStream(context.Background(), req.User)
	if err != nil {
		return err
	}
	for stream.Next(context.Background()) {
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
				h.zapLog.Debug(err.Error())
			}
		}
		go server.Send(&go_block.UserStream{
			StreamType: streamType,
			User:       userResp,
		})
	}
	return nil
}
