package user_repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongoRepository) DeleteNamespace(ctx context.Context, namespace string) error {
	filter := bson.M{"namespace": namespace}
	if _, err := r.collection.DeleteMany(ctx, filter); err != nil {
		return err
	}
	return nil
}
