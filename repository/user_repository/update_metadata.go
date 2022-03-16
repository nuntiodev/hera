package user_repository

import (
	"context"
	"errors"
	"github.com/softcorp-io/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongoRepository) UpdateMetadata(ctx context.Context, get *go_block.User, update *go_block.User, encryptionKey string) (*go_block.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateMetadata, update)
	if err := r.validate(actionUpdateMetadata, update); err != nil {
		return nil, err
	}
	getUser, err := r.Get(ctx, get, encryptionKey) // check if user encryption is turned on
	if err != nil {
		return nil, err
	}
	resp := *update
	if err := r.handleEncryption(getUser.Encrypted, update, encryptionKey); err != nil {
		return nil, err
	}
	updateUser := ProtoUserToUser(&go_block.User{
		Metadata:  update.Metadata,
		UpdatedAt: update.UpdatedAt,
	})
	mongoUpdate := bson.M{
		"$set": bson.M{
			"metadata":     updateUser.Metadata,
			"updated_at":   updateUser.UpdatedAt,
			"encrypted_at": updateUser.EncryptedAt,
		},
	}
	updateResult, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": getUser.Id, "namespace": getUser.Namespace},
		mongoUpdate,
	)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("could not find get")
	}
	return &resp, nil
}
