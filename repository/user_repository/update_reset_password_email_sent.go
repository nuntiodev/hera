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
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(user.EmailResetPasswordCode), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.EmailResetPasswordCode = string(hashedCode)
	mongoUpdate := bson.M{
		"$set": bson.M{
			"email_reset_password_code":    user.EmailResetPasswordCode,
			"reset_password_email_sent_at": time.Now().UTC(),
			"updated_at":                   time.Now().UTC(),
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
	resp.ResetPasswordEmailSentAt = time.Now().UTC()
	resp.UpdatedAt = time.Now().UTC()
	return UserToProtoUser(&resp), nil
}
