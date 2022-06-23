package user_repository

import (
	"context"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/models"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (r *mongodbRepository) UpdateMetadata(ctx context.Context, get *go_hera.User, update *go_hera.User) error {
	prepare(actionGet, get)
	prepare(actionUpdateMetadata, update)
	filter, err := getUserFilter(get)
	if err != nil {
		return err
	}
	if err := validateMetadata(update.Metadata); err != nil {
		return err
	}
	updateUser := models.ProtoUserToUser(&go_hera.User{
		Metadata:  update.Metadata,
		UpdatedAt: update.UpdatedAt,
	})
	if err := r.crypto.Encrypt(updateUser); err != nil {
		return err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"metadata":   updateUser.Metadata,
			"updated_at": time.Now(),
		},
	}
	if _, err := r.collection.UpdateOne(
		ctx,
		filter,
		mongoUpdate,
	); err != nil {
		return err
	}
	return nil
}
