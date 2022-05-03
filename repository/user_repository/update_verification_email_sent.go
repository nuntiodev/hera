package user_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (r *mongodbRepository) UpdateVerificationEmailSent(ctx context.Context, get *go_block.User) (*go_block.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	filter, err := getUserFilter(get)
	if err != nil {
		return nil, err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"verification_email_sent_at": time.Now(),
			"updated_at":                 time.Now(),
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
	resp.VerificationEmailSentAt = time.Now()
	resp.UpdatedAt = time.Now()
	return UserToProtoUser(&resp), nil
}
