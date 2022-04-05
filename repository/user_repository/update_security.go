package user_repository

import (
	"context"
	"errors"
	"github.com/softcorp-io/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (r *mongoRepository) UpdateSecurity(ctx context.Context, get *go_block.User) (*go_block.User, error) {
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
	update := ProtoUserToUser(get)
	// check if we need to encrypt the user
	if get.ExternalEncrypted {
		// user is already encrypted - disable external encryption
		update.ExternalEncrypted = false
	} else {
		update.ExternalEncrypted = true
		update.InternalEncrypted = len(r.internalEncryptionKeys) > 0
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
			"external_encrypted":        update.ExternalEncrypted,
			"internal_encrypted":        update.InternalEncrypted,
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
	get.ExternalEncrypted = update.ExternalEncrypted
	get.InternalEncrypted = update.InternalEncrypted
	get.InternalEncryptionLevel = int32(update.InternalEncryptionLevel)
	return get, nil
}
