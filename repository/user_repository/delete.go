package user_repository

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) Delete(ctx context.Context, user *go_block.User) error {
	prepare(actionGet, user)
	if err := r.validate(actionGet, user); err != nil {
		return err
	}
	filter := bson.M{}
	if user.Id != "" {
		filter = bson.M{"_id": user.Id}
	} else if user.Email != "" {
		filter = bson.M{"email_hash": fmt.Sprintf("%x", md5.Sum([]byte(user.Email)))}
	} else if user.Username != "" {
		filter = bson.M{"username": user.Username}
	}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return NoUsersDeletedErr
	}
	return nil
}
