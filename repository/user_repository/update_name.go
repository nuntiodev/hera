package user_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (r *mongodbRepository) UpdateName(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateName, update)
	if err := r.validate(actionUpdateName, update); err != nil {
		return nil, err
	}
	get, err := r.Get(ctx, get, true) // check if user encryption is turned on
	if err != nil {
		return nil, err
	}
	updateUser := ProtoUserToUser(&go_block.User{
		FirstName: update.FirstName,
		LastName:  update.LastName,
		UpdatedAt: update.UpdatedAt,
	})
	// transfer data from get to update
	updateUser.ExternalEncryptionLevel = int(get.ExternalEncryptionLevel)
	updateUser.InternalEncryptionLevel = int(get.InternalEncryptionLevel)
	// encrypt user if user has previously been encrypted
	if updateUser.ExternalEncryptionLevel > 0 || updateUser.InternalEncryptionLevel > 0 {
		if err := r.encryptUser(ctx, actionUpdateName, updateUser); err != nil {
			return nil, err
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"first_name": updateUser.LastName,
			"last_name":  updateUser.FirstName,
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
	get.FirstName = update.FirstName
	get.LastName = update.LastName
	get.UpdatedAt = ts.New(updateUser.UpdatedAt)
	return get, nil
}
