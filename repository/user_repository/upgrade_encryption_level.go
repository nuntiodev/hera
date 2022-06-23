package user_repository

import (
	"context"
	"fmt"
	"github.com/nuntiodev/hera/models"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (r *mongodbRepository) upgradeEncryptionLevel(ctx context.Context, user *models.User) error {
	if user == nil {
		return UserIsNilErr
	}
	if upgradable, err := r.crypto.Upgradeble(user); err != nil && !upgradable {
		return fmt.Errorf("could not upgrade with err: %v", err)
	}
	copy := *user
	if err := r.crypto.Encrypt(&copy); err != nil {
		return err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"username":     copy.Username,
			"email":        copy.Email,
			"image":        copy.Image,
			"metadata":     copy.Metadata,
			"first_name":   copy.FirstName,
			"last_name":    copy.LastName,
			"birthdate":    copy.Birthdate,
			"phone_number": copy.Phone,
			"updated_at":   time.Now(),
			"encrypted_at": time.Now(),
		},
	}
	filter, err := getUserFilter(models.UserToProtoUser(user))
	if err != nil {
		return err
	}
	if _, err := r.collection.UpdateOne(
		ctx,
		filter,
		mongoUpdate,
	); err != nil {
		return err
	}
	// update levels
	user.Username.InternalEncryptionLevel = copy.Username.InternalEncryptionLevel
	user.Username.ExternalEncryptionLevel = copy.Username.ExternalEncryptionLevel
	user.Email.InternalEncryptionLevel = copy.Email.InternalEncryptionLevel
	user.Email.ExternalEncryptionLevel = copy.Email.ExternalEncryptionLevel
	user.Image.InternalEncryptionLevel = copy.Image.InternalEncryptionLevel
	user.Image.ExternalEncryptionLevel = copy.Image.ExternalEncryptionLevel
	user.Metadata.InternalEncryptionLevel = copy.Metadata.InternalEncryptionLevel
	user.Metadata.ExternalEncryptionLevel = copy.Metadata.ExternalEncryptionLevel
	user.FirstName.InternalEncryptionLevel = copy.FirstName.InternalEncryptionLevel
	user.FirstName.ExternalEncryptionLevel = copy.FirstName.ExternalEncryptionLevel
	user.LastName.InternalEncryptionLevel = copy.LastName.InternalEncryptionLevel
	user.LastName.ExternalEncryptionLevel = copy.LastName.ExternalEncryptionLevel
	user.Birthdate.InternalEncryptionLevel = copy.Birthdate.InternalEncryptionLevel
	user.Birthdate.ExternalEncryptionLevel = copy.Birthdate.ExternalEncryptionLevel
	user.Phone.InternalEncryptionLevel = copy.Phone.InternalEncryptionLevel
	user.Phone.ExternalEncryptionLevel = copy.Phone.ExternalEncryptionLevel
	return nil
}
