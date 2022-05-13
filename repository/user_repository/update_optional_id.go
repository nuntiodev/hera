package user_repository

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"

	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) UpdateUsername(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateUsername, update)
	if err := r.validate(actionUpdateUsername, update); err != nil {
		return nil, err
	}
	updateUser := ProtoUserToUser(&go_block.User{
		Username:  update.Username,
		UpdatedAt: update.UpdatedAt,
	})
	mongoUpdate := bson.M{
		"$set": bson.M{
			"username":   updateUser.Username,
			"updated_at": updateUser.UpdatedAt,
		},
	}
	filter := bson.M{}
	if get.Id != "" {
		filter = bson.M{"_id": get.Id}
	} else if get.Email != "" {
		filter = bson.M{"email_hash": fmt.Sprintf("%x", md5.Sum([]byte(get.Email)))}
	} else if get.Username != "" {
		filter = bson.M{"username": get.Username}
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
