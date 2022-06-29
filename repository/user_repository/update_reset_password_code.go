package user_repository

import (
	"context"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (r *mongodbRepository) UpdateResetPasswordCode(ctx context.Context, user *go_hera.User) error {
	if user == nil {
		return UserIsNilErr
	}
	prepare(actionUpdateResetPasswordEmailSent, user)
	filter, err := getUserFilter(user)
	if err != nil {
		return err
	}
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(user.ResetPasswordCode), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.ResetPasswordCode = string(hashedCode)
	mongoUpdate := bson.M{
		"$set": bson.M{
			"reset_password_code":            user.ResetPasswordCode,
			"reset_password_code_sent_at":    time.Now(),
			"reset_password_code_expires_at": time.Now().Add(r.maxCodeVerificationAge),
			"updated_at":                     time.Now(),
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
