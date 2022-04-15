package user_repository

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"

	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (r *mongodbRepository) UpdateEmail(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateEmail, update)
	if err := r.validate(actionUpdateEmail, update); err != nil {
		return nil, err
	}
	emailHash := ""
	if update.Email != "" {
		emailHash = fmt.Sprintf("%x", md5.Sum([]byte(update.Email)))
	}
	get, err := r.Get(ctx, get, true) // check if user encryption is turned on
	if err != nil {
		return nil, err
	}
	updateUser := ProtoUserToUser(&go_block.User{
		Email:     update.Email,
		UpdatedAt: update.UpdatedAt,
	})
	// transfer data from get to update
	updateUser.ExternalEncrypted = get.ExternalEncrypted
	updateUser.ExternalEncryptionLevel = int(get.ExternalEncryptionLevel)
	updateUser.InternalEncrypted = get.InternalEncrypted
	updateUser.InternalEncryptionLevel = int(get.InternalEncryptionLevel)
	// encrypt user if user has previously been encrypted
	if updateUser.ExternalEncrypted || updateUser.InternalEncrypted {
		if err := r.encryptUser(ctx, actionUpdateEmail, updateUser); err != nil {
			return nil, err
		}
	}
	updateUser.EmailHash = emailHash
	mongoUpdate := bson.M{
		"$set": bson.M{
			"email":      updateUser.Email,
			"email_hash": updateUser.EmailHash,
			"updated_at": updateUser.UpdatedAt,
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
	get.Email = update.Email
	get.UpdatedAt = ts.New(updateUser.UpdatedAt)
	return get, nil
}
