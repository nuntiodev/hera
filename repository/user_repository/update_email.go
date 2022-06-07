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

func (r *mongodbRepository) UpdateEmail(ctx context.Context, get *go_block.User, update *go_block.User) (*models.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateEmail, update)
	if err := r.validate(actionUpdateEmail, update); err != nil {
		return nil, err
	}
	emailHash := fmt.Sprintf("%x", md5.Sum([]byte(update.Email)))
	// validate update is required
	updateUser := models.ProtoUserToUser(&go_block.User{
		Email:     update.Email,
		UpdatedAt: update.UpdatedAt,
	})
	updateUser.EmailIsVerified = false
	updateUser.VerificationEmailSentAt = time.Time{}
	if err := r.crypto.Encrypt(updateUser); err != nil {
		return nil, err
	}
	// check if new email already is verified previously; if so -> set to true
	for _, verifiedEmail := range get.VerifiedEmails {
		if verifiedEmail == emailHash {
			updateUser.EmailIsVerified = true
		}
	}
	updateUser.EmailHash = emailHash
	mongoUpdate := bson.M{
		"$set": bson.M{
			"email":                      updateUser.Email,
			"email_is_verified":          updateUser.EmailIsVerified,
			"verification_email_sent_at": updateUser.VerificationEmailSentAt,
			"email_hash":                 updateUser.EmailHash,
			"updated_at":                 updateUser.UpdatedAt,
		},
	}
	result := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": get.Id},
		mongoUpdate,
	)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var resp models.User
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	if err := r.crypto.Decrypt(&resp); err != nil {
		return nil, err
	}
	// set updated fields
	resp.Email = cryptox.Stringx{
		Body:                    update.Email,
		InternalEncryptionLevel: resp.Email.InternalEncryptionLevel,
		ExternalEncryptionLevel: resp.Email.ExternalEncryptionLevel,
	}
	resp.UpdatedAt = updateUser.UpdatedAt
	resp.EmailIsVerified = false
	return &resp, nil
}
