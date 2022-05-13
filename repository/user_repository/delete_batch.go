package user_repository

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"

	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) DeleteBatch(ctx context.Context, userBatch []*go_block.User) error {
	var ids []string
	var emails []string
	var usernames []string
	for _, user := range userBatch {
		if user == nil {
			return errors.New("a user is nil")
		}
		prepare(actionGet, user)
		if user.Id != "" {
			ids = append(ids, user.Id)
		} else if user.Email != "" {
			emails = append(emails, fmt.Sprintf("%x", md5.Sum([]byte(user.Email))))
		} else if user.Username != "" {
			ids = append(usernames, user.Username)
		}
	}
	var idsFilter bson.D
	var emailsFilter bson.D
	var usernamesFilter bson.D
	if len(ids) > 0 {
		idsFilter = bson.D{{"$in", ids}}
	}
	if len(emails) > 0 {
		emailsFilter = bson.D{{"$in", emails}}
	}
	if len(usernames) > 0 {
		usernamesFilter = bson.D{{"$in", usernames}}
	}
	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"_id", idsFilter}},
			bson.D{{"email", emailsFilter}},
			bson.D{{"username", usernamesFilter}},
		},
		},
	}
	if _, err := r.collection.DeleteMany(ctx, filter); err != nil {
		return err
	}
	return nil
}
