package user_repository

import (
	"context"
	"errors"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) UpdateSecurity(ctx context.Context, get *go_block.User) (*go_block.User, error) {
	prepare(actionGet, get)
	if r.externalEncryptionKey == "" {
		return nil, errors.New("cannot update security profile without an external key")
	}
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	get, err := r.Get(ctx, get, true) // check if user encryption is turned on
	if err != nil {
		return nil, err
	}
	// validate all keys are present at specific level
	if get.InternalEncryptionLevel > int32(len(r.internalEncryptionKeys)) {
		return nil, errors.New("not enough valid internal encryption keys to upgrade security")
	}
	update := ProtoUserToUser(get)
	// check if we need to encrypt the user
	if get.ExternalEncryptionLevel > 0 {
		// user is already encrypted - disable external encryption
		update.ExternalEncryptionLevel = 0
		// we still want to encrypt user under internal encryption keys
		encryptionKey, err := r.crypto.CombineSymmetricKeys(r.internalEncryptionKeys, len(r.internalEncryptionKeys))
		if err != nil {
			return nil, err
		}
		if err := r.encrypt(update, encryptionKey); err != nil {
			return nil, err
		}
		// also update to the newest internal encryption level
		update.InternalEncryptionLevel = len(r.internalEncryptionKeys)
	} else {
		update.ExternalEncryptionLevel = 1
		update.InternalEncryptionLevel = len(r.internalEncryptionKeys)
		if err := r.encryptUser(ctx, actionUpdateSecurity, update); err != nil {
			return nil, err
		}
	}
	update.UpdatedAt = time.Now()
	mongoUpdate := bson.M{
		"$set": bson.M{
			"email":                     update.Email,
			"image":                     update.Image,
			"external_encryption_level": update.ExternalEncryptionLevel,
			"internal_encryption_level": update.InternalEncryptionLevel,
			"metadata":                  update.Metadata,
			"updated_at":                update.UpdatedAt,
			"encrypted_at":              update.EncryptedAt,
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
	get.ExternalEncryptionLevel = int32(update.ExternalEncryptionLevel)
	get.InternalEncryptionLevel = int32(update.InternalEncryptionLevel)
	return get, nil
}
