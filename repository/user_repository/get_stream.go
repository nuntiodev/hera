package user_repository

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *mongoRepository) GetStream(ctx context.Context, get *go_block.User) (*mongo.ChangeStream, error) {
	matchPipeline := bson.D{}
	if get != nil && get.Id != "" {
		matchPipeline = bson.D{
			{"$match",
				bson.D{{"documentKey._id", get.Id}},
			},
		}
	}
	userStream, err := r.collection.Watch(ctx, mongo.Pipeline{matchPipeline})
	if err != nil {
		return nil, err
	}
	return userStream, nil
}
