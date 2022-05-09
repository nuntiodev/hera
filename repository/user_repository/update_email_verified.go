package user_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (r *mongodbRepository) UpdateEmailVerified(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error) {
	prepare(actionGet, get)
	prepare(actionUpdateEmailVerified, update)
	if err := r.validate(actionUpdateEmailVerified, update); err != nil {
		return nil, err
	}
	filter, err := getUserFilter(get)
	if err != nil {
		return nil, err
	}
	if update.EmailIsVerified {
		update.EmailVerifiedAt = ts.Now()
	}
	updateUser := ProtoUserToUser(&go_block.User{
		EmailIsVerified: update.EmailIsVerified,
		EmailVerifiedAt: update.EmailVerifiedAt,
		UpdatedAt:       update.UpdatedAt,
	})
	mongoUpdate := bson.M{
		"$set": bson.M{
			"email_is_verified": updateUser.EmailIsVerified,
			"email_verified_at": updateUser.EmailVerifiedAt,
			"updated_at":        updateUser.UpdatedAt,
		},
		"$inc": bson.D{{"verify_email_attempts", 1}},
	}
	if _, err := r.collection.UpdateOne(
		ctx,
		filter,
		mongoUpdate,
	); err != nil {
		return nil, err
	}
	// set updated fields
	get.EmailVerifiedAt = ts.New(updateUser.EmailVerifiedAt)
	get.EmailIsVerified = updateUser.EmailIsVerified
	get.UpdatedAt = ts.New(updateUser.UpdatedAt)
	return get, nil
}
