package user_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/softcorp-io/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	maxUserStreamSize = 75
)

func (r *mongoRepository) GetUsersStream(ctx context.Context, namespace string, userBatch []*go_block.User) (*mongo.ChangeStream, error) {
	if namespace == "" && len(userBatch) == 0 {
		return nil, errors.New("invalid request")
	}
	if len(userBatch) > maxUserStreamSize {
		return nil, errors.New(fmt.Sprintf("User batch size cannot exceed %d. Notice that we automatically handle following the newest users and update the stream.", maxUserStreamSize))
	}
	var matchPipeline bson.D
	// (create) stream match on specific namespace
	createStreamFilterNamespace := bson.D{
		{"operationType", "insert"},
		{"fullDocument.namespace", namespace},
	}
	var deleteStreamFilterMatchUsers bson.D
	var updateStreamFilterMatchUsers bson.D
	if len(userBatch) > 0 {
		var userIds []string
		for _, user := range userBatch {
			userIds = append(userIds, user.Id)
		}
		// (update) stream match for specific user ids
		deleteStreamFilterMatchUsers = bson.D{
			{"operationType", "delete"},
			{"documentKey._id", bson.D{{"$in", userIds}}},
		}
		// (delete) stream match for specific user ids
		updateStreamFilterMatchUsers = bson.D{
			{"operationType", "update"},
			{"documentKey._id", bson.D{{"$in", userIds}}},
		}
		if namespace == "" {
			matchPipeline = bson.D{
				{"$match",
					bson.M{"$or": bson.A{
						deleteStreamFilterMatchUsers,
						updateStreamFilterMatchUsers,
					},
					},
				},
			}
		} else {
			matchPipeline = bson.D{
				{"$match",
					bson.M{"$or": bson.A{
						createStreamFilterNamespace,
						deleteStreamFilterMatchUsers,
						updateStreamFilterMatchUsers,
					},
					},
				},
			}
		}
	} else {
		matchPipeline = bson.D{
			{"$match", createStreamFilterNamespace},
		}
	}
	userStream, err := r.collection.Watch(ctx, mongo.Pipeline{matchPipeline})
	if err != nil {
		return nil, err
	}
	return userStream, nil
}
