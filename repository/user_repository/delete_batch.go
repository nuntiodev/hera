package user_repository

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"

	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

/*
	DeleteBatch - this method deletes a batch of users by id, email or username.
*/
func (r *mongodbRepository) DeleteBatch(ctx context.Context, userBatch []*go_block.User) error {
	var ids []string
	var emails []string
	var usernames []string
	var phoneNumbers []string
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
			usernames = append(usernames, fmt.Sprintf("%x", md5.Sum([]byte(user.Username))))
		} else if user.PhoneNumber != "" {
			phoneNumbers = append(phoneNumbers, fmt.Sprintf("%x", md5.Sum([]byte(user.PhoneNumber))))
		}
	}
	var idsFilter bson.D
	var emailsFilter bson.D
	var usernamesFilter bson.D
	var phoneNumberFilter bson.D
	if len(ids) > 0 {
		idsFilter = bson.D{{"$in", ids}}
	}
	if len(emails) > 0 {
		emailsFilter = bson.D{{"$in", emails}}
	}
	if len(usernames) > 0 {
		usernamesFilter = bson.D{{"$in", usernames}}
	}
	if len(phoneNumbers) > 0 {
		phoneNumberFilter = bson.D{{"$in", phoneNumbers}}
	}
	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"_id", idsFilter}},
			bson.D{{"email_hash", emailsFilter}},
			bson.D{{"username_hash", usernamesFilter}},
			bson.D{{"phone_number_hash", phoneNumberFilter}},
		},
		},
	}
	if _, err := r.collection.DeleteMany(ctx, filter); err != nil {
		return err
	}
	return nil
}
