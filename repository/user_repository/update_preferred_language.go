package user_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (r *mongodbRepository) UpdatePreferredLanguage(ctx context.Context, get, update *go_block.User) (*models.User, error) {
	prepare(actionUpdatePreferredLanguage, update)
	if err := r.validate(actionUpdatePreferredLanguage, update); err != nil {
		return nil, err
	}
	filter, err := getUserFilter(get)
	if err != nil {
		return nil, err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"preferred_language": update.PreferredLanguage,
			"updated_at":         time.Now(),
		},
	}
	result := r.collection.FindOneAndUpdate(
		ctx,
		filter,
		mongoUpdate,
	)
	if err := result.Err(); err != nil {
		return nil, err
	}
	var resp models.User
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	if err := r.crypto.Decrypt(&resp); err != nil {
		return nil, err
	}
	// set updated fields
	resp.PreferredLanguage = update.PreferredLanguage
	resp.UpdatedAt = time.Now()
	return &resp, nil
}
