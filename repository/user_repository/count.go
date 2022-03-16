package user_repository

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/softcorp-io/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongoRepository) Count(ctx context.Context, namespace string) (int64, error) {
	filter := bson.M{"namespace": namespace}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *mongoRepository) Delete(ctx context.Context, user *go_block.User) error {
	prepare(actionGet, user)
	if err := r.validate(actionGet, user); err != nil {
		return err
	}
	filter := bson.M{}
	if user.Id != "" {
		filter = bson.M{"_id": user.Id, "namespace": user.Namespace}
	} else if user.Email != "" {
		filter = bson.M{"email_hash": fmt.Sprintf("%x", md5.Sum([]byte(user.Email))), "namespace": user.Namespace}
	} else if user.OptionalId != "" {
		filter = bson.M{"optional_id": user.OptionalId, "namespace": user.Namespace}
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
