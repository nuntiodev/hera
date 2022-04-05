package user_repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (r *mongoRepository) upgradeInternalEncryptionLevel(ctx context.Context, user *User) error {
	if user == nil {
		return UserIsNilErr
	}
	if len(r.internalEncryptionKeys) == 0 {
		return errors.New("no encryption keys to upgrade internal encryption")
	}
	get, err := r.Get(ctx, UserToProtoUser(user), false)
	if err != nil {
		return err
	}
	// make sure external encryption key is present or user is not encrypted by external key
	if user.ExternalEncrypted && r.externalEncryptionKey == "" {
		return errors.New("the requested user is encrypted by an external encryption key and can no be upgraded without that external encryption key")
	}
	// handle upgrade of internal key
	if (get.InternalEncrypted || get.InternalEncryptionLevel == 0) && get.InternalEncryptionLevel < int32(len(r.internalEncryptionKeys)) {
		update := ProtoUserToUser(get)
		if err := r.encryptUser(ctx, actionUpgradeEncryption, update); err != nil {
			return err
		}
		mongoUpdate := bson.M{
			"$set": bson.M{
				"email":                     update.Email,
				"image":                     update.Image,
				"external_encrypted":        update.ExternalEncrypted,
				"internal_encrypted":        update.InternalEncrypted,
				"internal_encryption_level": len(r.internalEncryptionKeys),
				"metadata":                  update.Metadata,
				"updated_at":                time.Now(),
				"encrypted_at":              time.Now(),
			},
		}
		if err := r.collection.FindOneAndUpdate(
			ctx,
			bson.M{"_id": get.Id},
			mongoUpdate,
		).Err(); err != nil {
			return err
		}
		user.InternalEncrypted = true
		user.InternalEncryptionLevel = len(r.internalEncryptionKeys)
	}
	return nil
}
