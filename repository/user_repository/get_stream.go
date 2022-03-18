package user_repository

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *mongoRepository) GetStream(ctx context.Context, get *go_block.User) (*mongo.ChangeStream, error) {
	pipeline := mongo.Pipeline{}
	if get != nil && get.Id != "" {
		pipeline = mongo.Pipeline{
			bson.D{
				{"$match",
					bson.D{{"documentKey._id", get.Id}},
				},
			},
		}
	}
	userStream, err := r.collection.Watch(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	return userStream, nil
}
