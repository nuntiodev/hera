package user_repository

import (
	"context"
	"github.com/nuntiodev/hera-proto/go_hera"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (r *mongodbRepository) UpdatePhoneVerificationCode(ctx context.Context, user *go_hera.User) error {
	if user == nil {
		return UserIsNilErr
	}
	prepare(actionUpdateTextVerificationCode, user)
	filter, err := getUserFilter(user)
	if err != nil {
		return err
	}
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(user.PhoneVerificationCode), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PhoneVerificationCode = string(hashedCode)
	mongoUpdate := bson.M{
		"$set": bson.M{
			"phone_verification_code":      user.PhoneVerificationCode,
			"verification_text_sent_at":    time.Now(),
			"verification_text_expires_at": time.Now().Add(r.maxCodeVerificationAge),
			"updated_at":                   time.Now(),
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
