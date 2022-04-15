package user_repository

import (
	"context"
	"errors"

	"github.com/io-nuntio/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (r *mongodbRepository) UpdateImage(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateImage, update)
	if err := r.validate(actionUpdateImage, update); err != nil {
		return nil, err
	}
	get, err := r.Get(ctx, get, true) // check if user encryption is turned on
	if err != nil {
		return nil, err
	}
	updateUser := ProtoUserToUser(&go_block.User{
		Image:     update.Image,
		UpdatedAt: update.UpdatedAt,
	})
	// transfer data from get to update
	updateUser.ExternalEncrypted = get.ExternalEncrypted
	updateUser.ExternalEncryptionLevel = int(get.ExternalEncryptionLevel)
	updateUser.InternalEncrypted = get.InternalEncrypted
	updateUser.InternalEncryptionLevel = int(get.InternalEncryptionLevel)
	// encrypt user if user has previously been encrypted
	if updateUser.ExternalEncrypted || updateUser.InternalEncrypted {
		if err := r.encryptUser(ctx, actionUpdateImage, updateUser); err != nil {
			return nil, err
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"image":        updateUser.Image,
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
	get.Image = update.Image
	get.UpdatedAt = ts.New(updateUser.UpdatedAt)
	return get, nil
}
