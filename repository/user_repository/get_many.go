package user_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/models"
	"go.mongodb.org/mongo-driver/bson"
)

/*
	GetMany - this method fetches an array of users by id.
*/
func (r *mongodbRepository) GetMany(ctx context.Context, users []*go_hera.User) ([]*models.User, error) {
	if users == nil {
		return nil, errors.New("users array is nil")
	}
	var userIds []string
	for _, user := range users {
		if user == nil || user.Id == "" {
			return nil, errors.New("user is nil or user id is empty in array")
		}
		userIds = append(userIds, user.Id)
	}
	var resp []*models.User
	cursor, err := r.collection.Find(ctx, bson.M{
		"_id": bson.M{"$in": userIds},
	})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		// check if external encryption has been applied
		if err := r.crypto.Decrypt(&user); err != nil {
			return nil, err
		}
		// check if we should upgrade the encryption level
		if upgradable, _ := r.crypto.Upgradeble(&user); upgradable {
			if err := r.upgradeEncryptionLevel(ctx, &user); err != nil {
				return nil, err
			}
		}
		resp = append(resp, &user)
	}
	if len(users) != len(resp) {
		return nil, errors.New("length of users does not equal length of resp which means that not all users could be found")
	}
	return resp, nil
}
