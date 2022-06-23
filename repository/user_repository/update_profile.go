package user_repository

import (
	"context"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/models"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (r *mongodbRepository) UpdateProfile(ctx context.Context, get *go_hera.User, update *go_hera.User) error {
	prepare(actionGet, get)
	prepare(actionUpdateProfile, update)
	// check if we can find the user
	filter, err := getUserFilter(get)
	if err != nil {
		return err
	}
	// validate data
	if update == nil {
		return UpdateIsNil
	} else if update.Image == nil && update.FirstName == nil && update.LastName == nil && update.Birthdate == nil {
		return NothingToUpdateErr
	}
	// encrypt user and build mongo data update
	updateUser := models.ProtoUserToUser(&go_hera.User{
		FirstName:         update.FirstName,
		LastName:          update.LastName,
		Image:             update.Image,
		Birthdate:         update.Birthdate,
		PreferredLanguage: update.PreferredLanguage,
	})
	if err := r.crypto.Encrypt(updateUser); err != nil {
		return err
	}
	mongoUpdate := bson.M{}
	// some fields are optional and should only be updated if the fields are not nil
	if update.Image != nil {
		mongoUpdate["birthdate"] = updateUser.Birthdate
	}
	if update.FirstName != nil {
		mongoUpdate["first_name"] = updateUser.FirstName
	}
	if update.LastName != nil {
		mongoUpdate["last_name"] = updateUser.LastName
	}
	if update.Birthdate != nil {
		mongoUpdate["birthdate"] = updateUser.Birthdate
	}
	if update.PreferredLanguage != nil {
		mongoUpdate["preferred_language"] = updateUser.PreferredLanguage
	}
	mongoUpdate["updated_at"] = time.Now()
	if _, err := r.collection.UpdateOne(
		ctx,
		filter,
		bson.M{
			"$set": mongoUpdate,
		},
	); err != nil {
		return err
	}
	return nil
}
