package user_repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}
	return count, nil
}
