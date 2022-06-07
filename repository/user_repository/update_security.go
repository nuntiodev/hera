package user_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/nuntiodev/nuntio-user-block/models"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) UpdateSecurity(ctx context.Context, get *go_block.User) (*models.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	eKeys, _ := r.crypto.GetExternalEncryptionKeys()
	iKeys, _ := r.crypto.GetInternalEncryptionKeys()
	if len(eKeys) == 0 && len(iKeys) == 0 {
		return nil, errors.New("need at least one external or internal encryption key to update the encryption level")
	}
	resp, err := r.Get(ctx, get) // check if user encryption is turned on
	if err != nil {
		return nil, err
	}
	external, internal := r.crypto.EncryptionLevel(resp)
	if external == 0 && internal == 0 {
		// user is not encrypted -> encrypt and store
		if err := r.crypto.Encrypt(resp); err != nil {
			return nil, err
		}
	} else {
		// user is already encrypted and has been decrypted
		// not set zero in all stringx fields
		if err := r.crypto.SetZero(resp); err != nil {
			return nil, err
		}
	}
	resp.UpdatedAt = time.Now()
	mongoUpdate := bson.M{
		"$set": bson.M{
			"username":     resp.Username,
			"email":        resp.Email,
			"image":        resp.Image,
			"metadata":     resp.Metadata,
			"first_name":   resp.FirstName,
			"last_name":    resp.LastName,
			"birthdate":    resp.Birthdate,
			"phone_number": resp.PhoneNumber,
			"updated_at":   resp.UpdatedAt,
			"encrypted_at": resp.EncryptedAt,
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
	/*
		copy.Username.InternalEncryptionLevel = resp.Username.InternalEncryptionLevel
		copy.Username.ExternalEncryptionLevel = resp.Username.ExternalEncryptionLevel
		copy.Email.InternalEncryptionLevel = resp.Email.InternalEncryptionLevel
		copy.Email.ExternalEncryptionLevel = resp.Email.ExternalEncryptionLevel
		copy.Image.InternalEncryptionLevel = resp.Image.InternalEncryptionLevel
		copy.Image.ExternalEncryptionLevel = resp.Image.ExternalEncryptionLevel
		copy.Metadata.InternalEncryptionLevel = resp.Metadata.InternalEncryptionLevel
		copy.Metadata.ExternalEncryptionLevel = resp.Metadata.ExternalEncryptionLevel
		copy.FirstName.InternalEncryptionLevel = resp.FirstName.InternalEncryptionLevel
		copy.FirstName.ExternalEncryptionLevel = resp.FirstName.ExternalEncryptionLevel
		copy.LastName.InternalEncryptionLevel = resp.LastName.InternalEncryptionLevel
		copy.LastName.ExternalEncryptionLevel = resp.LastName.ExternalEncryptionLevel
		copy.Birthdate.InternalEncryptionLevel = resp.Birthdate.InternalEncryptionLevel
		copy.Birthdate.ExternalEncryptionLevel = resp.Birthdate.ExternalEncryptionLevel
		copy.PhoneNumber.InternalEncryptionLevel = resp.PhoneNumber.InternalEncryptionLevel
		copy.PhoneNumber.ExternalEncryptionLevel = resp.PhoneNumber.ExternalEncryptionLevel
	*/
	fmt.Println(resp)
	return resp, nil
}
