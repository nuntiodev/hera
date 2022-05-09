package user_repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (r *mongodbRepository) isEncryptionLevelUpgradable(user *User) bool {
	if user == nil {
		return false
	}
	validExternalEncConfig := user.ExternalEncryptionLevel == 0 || (user.ExternalEncryptionLevel > 0 && r.externalEncryptionKey != "")
	userNeedsInternalUpgrading := len(r.internalEncryptionKeys) > user.InternalEncryptionLevel
	userIsLevelZero := user.InternalEncryptionLevel == 0 && len(r.internalEncryptionKeys) > 0
	return (validExternalEncConfig && userNeedsInternalUpgrading) || userIsLevelZero
}

func (r *mongodbRepository) upgradeInternalEncryptionLevel(ctx context.Context, user *User) error {
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
	if user.ExternalEncryptionLevel > 0 && r.externalEncryptionKey == "" {
		return errors.New("the requested user is encrypted by an external encryption key and can no be upgraded without that external encryption key")
	}
	// handle upgrade of internal key
	if get.InternalEncryptionLevel < int32(len(r.internalEncryptionKeys)) {
		update := ProtoUserToUser(get)
		if err := r.encryptUser(ctx, actionUpgradeEncryption, update); err != nil {
			return err
		}
		update.InternalEncryptionLevel = len(r.internalEncryptionKeys)
		if r.externalEncryptionKey != "" {
			user.ExternalEncryptionLevel = 1
		}
		mongoUpdate := bson.M{
			"$set": bson.M{
				"email":                     update.Email,
				"image":                     update.Image,
				"first_name":                update.FirstName,
				"last_name":                 update.LastName,
				"birthdate":                 update.Birthdate,
				"internal_encryption_level": update.InternalEncryptionLevel,
				"external_encryption_level": update.ExternalEncryptionLevel,
				"metadata":                  update.Metadata,
				"updated_at":                time.Now().UTC(),
				"encrypted_at":              time.Now().UTC(),
			},
		}
		if err := r.collection.FindOneAndUpdate(
			ctx,
			bson.M{"_id": get.Id},
			mongoUpdate,
		).Err(); err != nil {
			return err
		}
		// set updated fields for user
		user.InternalEncryptionLevel = update.InternalEncryptionLevel
		user.ExternalEncryptionLevel = update.ExternalEncryptionLevel
	}
	return nil
}
