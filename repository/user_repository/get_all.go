package user_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuntiodev/nuntio-user-block/models"

	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
	GetAll - this method fetches all users matching the filter.
*/
func (r *mongodbRepository) GetAll(ctx context.Context, userFilter *go_block.UserFilter) ([]*models.User, error) {
	var resp []*models.User
	sortOptions := options.FindOptions{}
	limitOptions := options.Find()
	limitOptions.SetLimit(maximumGetLimit)
	filter := bson.M{}
	if userFilter != nil {
		order := -1
		if userFilter.Order == go_block.UserFilter_INC {
			order = 1
		}
		switch userFilter.Sort {
		case go_block.UserFilter_CREATED_AT:
			sortOptions.SetSort(bson.D{{"created_at", order}, {"_id", order}})
		case go_block.UserFilter_UPDATE_AT:
			sortOptions.SetSort(bson.D{{"updated_at", order}, {"_id", order}})
		default:
			return nil, errors.New("invalid sorting")
		}
		if userFilter.From >= 0 && userFilter.To > 0 {
			if userFilter.To-userFilter.From > maximumGetLimit {
				return nil, errors.New(fmt.Sprintf("exceeding maximum range of %d", maximumGetLimit))
			}
			limitOptions.SetLimit(int64(userFilter.To - userFilter.From))
			limitOptions.SetSkip(int64(userFilter.From))
		}
	}
	cursor, err := r.collection.Find(ctx, filter, &sortOptions, limitOptions)
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

	return resp, nil
}
