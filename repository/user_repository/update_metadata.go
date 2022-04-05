package user_repository

import (
	"context"
	"errors"
	"github.com/softcorp-io/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongoRepository) UpdateMetadata(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateMetadata, update)
	if err := r.validate(actionUpdateMetadata, update); err != nil {
		return nil, err
	}
	get, err := r.Get(ctx, get, true) // check if user encryption is turned on
	if err != nil {
		return nil, err
	}
	updateUser := ProtoUserToUser(&go_block.User{
		Metadata:  update.Metadata,
		UpdatedAt: update.UpdatedAt,
	})
	// transfer data from get to update
	updateUser.ExternalEncrypted = get.ExternalEncrypted
	updateUser.ExternalEncryptionLevel = int(get.ExternalEncryptionLevel)
	updateUser.InternalEncrypted = get.InternalEncrypted
	updateUser.InternalEncryptionLevel = int(get.InternalEncryptionLevel)
	// encrypt user if user has previously been encrypted
	if updateUser.ExternalEncrypted || updateUser.InternalEncrypted {
		if err := r.encryptUser(ctx, actionUpdateMetadata, updateUser); err != nil {
			return nil, err
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"metadata":     updateUser.Metadata,
			"updated_at":   updateUser.UpdatedAt,
			"encrypted_at": updateUser.EncryptedAt,
		},
	}
	updateResult, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": get.Id},
		mongoUpdate,
	)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("could not find get")
	}
	// set updated fields
	get.Metadata = update.Metadata
	get.UpdatedAt = update.UpdatedAt
	return get, nil
}
