package user_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (r *mongodbRepository) UpdateResetPasswordEmailSent(ctx context.Context, user *go_block.User) (*go_block.User, error) {
	prepare(actionUpdateResetPasswordEmailSent, user)
	if err := r.validate(actionUpdateResetPasswordEmailSent, user); err != nil {
		return nil, err
	}
	filter, err := getUserFilter(user)
	if err != nil {
		return nil, err
	}
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(user.ResetPasswordCode), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.ResetPasswordCode = string(hashedCode)
	mongoUpdate := bson.M{
		"$set": bson.M{
			"reset_password_code":             user.ResetPasswordCode,
			"reset_password_email_sent_at":    time.Now(),
			"reset_password_email_expires_at": time.Now().Add(r.maxEmailVerificationAge),
			"reset_password_attempts":         int32(0),
			"updated_at":                      time.Now(),
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
	resp.ResetPasswordEmailSentAt = time.Now()
	resp.UpdatedAt = time.Now()
	return UserToProtoUser(&resp), nil
}
