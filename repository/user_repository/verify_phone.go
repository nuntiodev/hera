package user_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (r *mongodbRepository) VerifyPhone(ctx context.Context, user *go_hera.User, isVerified bool) error {
	if user == nil {
		return UserIsNilErr
	}
	prepare(actionGet, user)
	_, _, phoneHash := generateUserHashes(&go_hera.User{Email: user.Email})
	if phoneHash == "" {
		return errors.New("no valid phone hash")
	}
	mongoUpdate := bson.M{}
	setData := bson.M{
		"updated_at": time.Now(),
	}
	if isVerified {
		setData["verify_phone_attempts"] = int32(0)
	} else {
		mongoUpdate["$inc"] = bson.D{{"verify_phone_attempts", int32(1)}}
		mongoUpdate["$addToSet"] = bson.D{{"verified_phone_numbers", phoneHash}}
	}
	if _, err := r.collection.UpdateOne(
		ctx,
		bson.M{"phone_hash": phoneHash},
		mongoUpdate,
	); err != nil {
		return err
	}
	return nil
}
