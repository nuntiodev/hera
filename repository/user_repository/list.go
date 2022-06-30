package user_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuntiodev/hera/models"

	"github.com/nuntiodev/hera-sdks/go_hera"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
	List - this method fetches all users matching the filter.
*/
func (r *mongodbRepository) List(ctx context.Context, query *go_hera.Query) ([]*go_hera.User, error) {
	if query == nil {
		return nil, errors.New("query is nil")
	}
	var resp []*go_hera.User
	sortOptions := options.FindOptions{}
	limitOptions := options.Find()
	limitOptions.SetLimit(maximumGetLimit)
	filter := bson.M{}
	if query != nil {
		order := -1
		if query.Order == go_hera.Query_INC {
			order = 1
		}
		switch query.Sort {
		case go_hera.Query_CREATED_AT:
			sortOptions.SetSort(bson.D{{"created_at", order}, {"_id", order}})
		case go_hera.Query_UPDATE_AT:
			sortOptions.SetSort(bson.D{{"updated_at", order}, {"_id", order}})
		default:
			return nil, errors.New("invalid sorting")
		}
		if query.From >= 0 && query.To > 0 {
			if query.To-query.From > maximumGetLimit {
				return nil, errors.New(fmt.Sprintf("exceeding maximum range of %d", maximumGetLimit))
			}
			limitOptions.SetLimit(int64(query.To - query.From))
			limitOptions.SetSkip(int64(query.From))
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
		resp = append(resp, models.UserToProtoUser(&user))
	}

	return resp, nil
}
