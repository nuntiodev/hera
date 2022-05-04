package user_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (r *mongodbRepository) UpdateVerificationEmailSent(ctx context.Context, user *go_block.User) (*go_block.User, error) {
	prepare(actionUpdateVerificationEmailSent, user)
	if err := r.validate(actionUpdateVerificationEmailSent, user); err != nil {
		return nil, err
	}
	filter, err := getUserFilter(user)
	if err != nil {
		return nil, err
	}
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(user.VerificationCode), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.VerificationCode = string(hashedCode)
	mongoUpdate := bson.M{
		"$set": bson.M{
			"verification_code":          user.VerificationCode,
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
