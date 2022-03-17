package user_repository

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/softcorp-io/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongoRepository) DeleteBatch(ctx context.Context, userBatch []*go_block.User) error {
	var ids []string
	var emails []string
	var optionalIds []string
	for _, user := range userBatch {
		if user == nil {
			return errors.New("a user is nil")
		}
		prepare(actionGet, user)
		if user.Id != "" {
			ids = append(ids, user.Id)
		} else if user.Email != "" {
			emails = append(emails, fmt.Sprintf("%x", md5.Sum([]byte(user.Email))))
		} else if user.OptionalId != "" {
			ids = append(optionalIds, user.OptionalId)
		}
	}
	var idsFilter bson.D
	var emailsFilter bson.D
	var optionalIdsFilter bson.D
	if len(ids) > 0 {
		idsFilter = bson.D{{"$in", ids}}
	}
	if len(emails) > 0 {
		emailsFilter = bson.D{{"$in", emails}}
	}
	if len(optionalIds) > 0 {
		optionalIdsFilter = bson.D{{"$in", optionalIds}}
	}
	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"_id", idsFilter}},
			bson.D{{"email", emailsFilter}},
			bson.D{{"optional_id", optionalIdsFilter}},
		},
		},
	}
	if _, err := r.collection.DeleteMany(ctx, filter); err != nil {
		return err
	}
	return nil
}
