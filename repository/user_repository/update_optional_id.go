package user_repository

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"

	"github.com/io-nuntio/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) UpdateOptionalId(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateOptionalId, update)
	if err := r.validate(actionUpdateOptionalId, update); err != nil {
		return nil, err
	}
	updateUser := ProtoUserToUser(&go_block.User{
		OptionalId: update.OptionalId,
		UpdatedAt:  update.UpdatedAt,
	})
	mongoUpdate := bson.M{
		"$set": bson.M{
			"optional_id": updateUser.OptionalId,
			"updated_at":  updateUser.UpdatedAt,
		},
	}
	filter := bson.M{}
	if get.Id != "" {
		filter = bson.M{"_id": get.Id}
	} else if get.Email != "" {
		filter = bson.M{"email_hash": fmt.Sprintf("%x", md5.Sum([]byte(get.Email)))}
	} else if get.OptionalId != "" {
		filter = bson.M{"optional_id": get.OptionalId}
	}
	updateResult, err := r.collection.UpdateOne(
		ctx,
		filter,
		mongoUpdate,
	)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("could not find get")
	}
	return update, nil
}
