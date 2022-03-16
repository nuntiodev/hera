package user_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/softcorp-io/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *mongoRepository) GetAll(ctx context.Context, userFilter *go_block.UserFilter, namespace string, encryptionKey string) ([]*go_block.User, error) {
	var resp []*go_block.User
	sortOptions := options.FindOptions{}
	limitOptions := options.Find()
	limitOptions.SetLimit(maximumGetLimit)
	filter := bson.M{"namespace": namespace}
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
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		protoUser := UserToProtoUser(&user)
		if user.Encrypted == true && encryptionKey != "" {
			if err := r.crypto.DecryptUser(encryptionKey, protoUser); err != nil {
				return nil, err
			}
		}
		resp = append(resp, protoUser)
	}

	return resp, nil
}
