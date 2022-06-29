package user_repository

import (
	"context"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (r *mongodbRepository) UpdateEmailVerificationCode(ctx context.Context, user *go_hera.User) error {
	if user == nil {
		return UserIsNilErr
	}
	prepare(actionUpdateEmailVerificationCode, user)
	filter, err := getUserFilter(user)
	if err != nil {
		return err
	}
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(user.EmailVerificationCode), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.EmailVerificationCode = string(hashedCode)
	mongoUpdate := bson.M{
		"$set": bson.M{
			"email_verification_code":       user.EmailVerificationCode,
			"verification_email_sent_at":    time.Now(),
			"verification_email_expires_at": time.Now().Add(r.maxCodeVerificationAge),
			"updated_at":                    time.Now(),
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
