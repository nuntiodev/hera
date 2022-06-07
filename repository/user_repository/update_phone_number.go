package user_repository

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/x/cryptox"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) UpdatePhoneNumber(ctx context.Context, get *go_block.User, update *go_block.User) (*models.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdatePhoneNumber, update)
	if err := r.validate(actionUpdatePhoneNumber, update); err != nil {
		return nil, err
	}
	phoneNumberHash := fmt.Sprintf("%x", md5.Sum([]byte(update.PhoneNumber)))
	updateUser := models.ProtoUserToUser(&go_block.User{
		PhoneNumber:     update.PhoneNumber,
		PhoneNumberHash: phoneNumberHash,
		UpdatedAt:       update.UpdatedAt,
	})
	updateUser.PhoneNumberIsVerified = false
	updateUser.VerificationTextSentAt = time.Time{}
	// check if new phone number already is verified previously; if so -> set to true
	for _, verifiedPhoneNumber := range get.VerifiedPhoneNumbers {
		if verifiedPhoneNumber == phoneNumberHash {
			updateUser.PhoneNumberIsVerified = true
		}
	}
	// encrypt user if user has previously been encrypted
	if err := r.crypto.Encrypt(updateUser); err != nil {
		return nil, err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"phone_number":              updateUser.PhoneNumber,
			"phone_number_is_verified":  updateUser.PhoneNumberIsVerified,
			"verification_text_sent_at": updateUser.VerificationTextSentAt,
			"phone_number_hash":         updateUser.PhoneNumberHash,
			"updated_at":                updateUser.UpdatedAt,
		},
	}
	result := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": get.Id},
		mongoUpdate,
	)
	if err := result.Err(); err != nil {
		return nil, err
	}
	var resp models.User
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	if err := r.crypto.Decrypt(&resp); err != nil {
		return nil, err
	}
	// set updated fields
	resp.PhoneNumber = cryptox.Stringx{
		Body:                    update.PhoneNumber,
		InternalEncryptionLevel: resp.PhoneNumber.InternalEncryptionLevel,
		ExternalEncryptionLevel: resp.PhoneNumber.ExternalEncryptionLevel,
	}
	resp.UpdatedAt = updateUser.UpdatedAt
	resp.PhoneNumberIsVerified = false
	resp.PhoneNumberHash = phoneNumberHash
	return &resp, nil
}
