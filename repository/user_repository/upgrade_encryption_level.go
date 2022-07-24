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
			"phone":        copy.Phone,
			"gender":       copy.Gender,
			"country":      copy.Country,
			"city":         copy.City,
			"address":      copy.Address,
			"postal_code":  copy.PostalCode,
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
	user.Username.EncryptionLevel = copy.Username.EncryptionLevel
	user.Email.EncryptionLevel = copy.Email.EncryptionLevel
	user.Image.EncryptionLevel = copy.Image.EncryptionLevel
	user.Metadata.EncryptionLevel = copy.Metadata.EncryptionLevel
	user.FirstName.EncryptionLevel = copy.FirstName.EncryptionLevel
	user.LastName.EncryptionLevel = copy.LastName.EncryptionLevel
	user.Birthdate.EncryptionLevel = copy.Birthdate.EncryptionLevel
	user.Phone.EncryptionLevel = copy.Phone.EncryptionLevel
	return nil
}
