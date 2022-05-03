package user_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (r *mongodbRepository) UpdateEmailVerified(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
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
	}
	result := r.collection.FindOneAndUpdate(
		ctx,
		filter,
		mongoUpdate,
	)
	if err := result.Err(); err != nil {
		return nil, err
	}
	var resp User
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	// set updated fields
	resp.EmailVerifiedAt = updateUser.EmailVerifiedAt
	resp.EmailIsVerified = updateUser.EmailIsVerified
	resp.UpdatedAt = updateUser.UpdatedAt
	return UserToProtoUser(&resp), nil
}
