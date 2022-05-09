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
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(user.EmailVerificationCode), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.EmailVerificationCode = string(hashedCode)
	mongoUpdate := bson.M{
		"$set": bson.M{
			"email_verification_code":       user.EmailVerificationCode,
			"verification_email_sent_at":    time.Now().UTC(),
			"verification_email_expires_at": time.Now().UTC().Add(time.Minute * 15),
			"verify_email_attempts":         int32(0),
			"updated_at":                    time.Now().UTC(),
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
	resp.VerificationEmailSentAt = time.Now().UTC()
	resp.VerificationEmailExpiresAt = time.Now().UTC().Add(time.Minute * 15)
	resp.UpdatedAt = time.Now().UTC()
	return UserToProtoUser(&resp), nil
}
