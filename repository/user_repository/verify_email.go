package user_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-proto/go_hera"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (r *mongodbRepository) VerifyEmail(ctx context.Context, user *go_hera.User, isVerified bool) error {
	if user == nil {
		return UserIsNilErr
	}
	prepare(actionGet, user)
	_, emailHash, _ := generateUserHashes(&go_hera.User{Email: user.Email})
	if emailHash == "" {
		return errors.New("no valid email hash")
	}
	mongoUpdate := bson.M{}
	setData := bson.M{
		"updated_at": time.Now(),
	}
	if isVerified {
		setData["verify_email_attempts"] = int32(0)
	} else {
		mongoUpdate["$inc"] = bson.D{{"verify_email_attempts", int32(1)}}
		mongoUpdate["$addToSet"] = bson.D{{"verified_emails", emailHash}}
	}
	if _, err := r.collection.UpdateOne(
		ctx,
		bson.M{"email_hash": emailHash},
		mongoUpdate,
	); err != nil {
		return err
	}
	return nil
}
