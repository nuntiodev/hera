package user_repository

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (r *mongodbRepository) UpdatePhoneNumber(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdatePhoneNumber, update)
	if err := r.validate(actionUpdatePhoneNumber, update); err != nil {
		return nil, err
	}
	phoneNumberHash := fmt.Sprintf("%x", md5.Sum([]byte(update.Email)))
	get, err := r.Get(ctx, get, true) // check if user encryption is turned on
	if err != nil {
		return nil, err
	}
	// validate update is required
	if get.PhoneNumberHash == phoneNumberHash {
		return nil, errors.New("email is identical to current email")
	}
	updateUser := ProtoUserToUser(&go_block.User{
		PhoneNumber: update.PhoneNumber,
		UpdatedAt:   update.UpdatedAt,
	})
	// transfer data from get to update
	updateUser.ExternalEncryptionLevel = int(get.ExternalEncryptionLevel)
	updateUser.InternalEncryptionLevel = int(get.InternalEncryptionLevel)
	updateUser.PhoneNumberIsVerified = false
	updateUser.VerificationTextSentAt = time.Time{}
	// encrypt user if user has previously been encrypted
	if updateUser.InternalEncryptionLevel > 0 || updateUser.ExternalEncryptionLevel > 0 {
		if err := r.encryptUser(ctx, actionUpdateEmail, updateUser); err != nil {
			return nil, err
		}
	}
	// check if new email already is verified previously; if so -> set to true
	for _, verifiedPhoneNumber := range get.VerifiedPhoneNumbers {
		if verifiedPhoneNumber == phoneNumberHash {
			updateUser.PhoneNumberIsVerified = true
		}
	}
	updateUser.PhoneNumberHash = phoneNumberHash
	mongoUpdate := bson.M{
		"$set": bson.M{
			"phone_number":              updateUser.PhoneNumber,
			"phone_number_is_verified":  updateUser.PhoneNumberIsVerified,
			"verification_text_sent_at": updateUser.VerificationTextSentAt,
			"phone_number_hash":         updateUser.PhoneNumberHash,
			"updated_at":                updateUser.UpdatedAt,
		},
	}
	updateResult, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": get.Id},
		mongoUpdate,
	)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("could not find get")
	}
	// set updated fields
	get.Email = update.Email
	get.UpdatedAt = ts.New(updateUser.UpdatedAt)
	get.EmailIsVerified = false
	return get, nil
}
