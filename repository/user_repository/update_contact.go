package user_repository

import (
	"context"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/models"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (r *mongodbRepository) UpdateContact(ctx context.Context, get *go_hera.User, update *go_hera.User) error {
	prepare(actionGet, get)
	prepare(actionUpdateContact, update)
	// check if we can find the user
	filter, err := getUserFilter(get)
	if err != nil {
		return err
	}
	// validate data
	if update.Email == nil && update.Username == nil && update.Phone == nil {
		return NothingToUpdateErr
	} else if err := validateEmail(update.GetEmail()); err != nil {
		return err
	} else if err := validatePhone(update.GetPhone()); err != nil {
		return err
	}
	// encrypt user and build mongo data update
	updateUser := models.ProtoUserToUser(&go_hera.User{
		Email:     update.Email,
		Phone:     update.Phone,
		Username:  update.Username,
		UpdatedAt: update.UpdatedAt,
	})
	if err := r.crypto.Encrypt(updateUser); err != nil {
		return err
	}
	mongoUpdate := bson.M{}
	// some fields are optional and should only be updated if the fields are not nil
	emailHash, usernameHash, phoneHash := generateUserHashes(update)
	if update.Email != nil {
		mongoUpdate["verification_email_sent_at"] = time.Time{}
		mongoUpdate["email"] = updateUser.Email
		mongoUpdate["email_hash"] = emailHash
	}
	if update.Phone != nil {
		mongoUpdate["phone"] = updateUser.Phone
		mongoUpdate["phone_hash"] = phoneHash
	}
	if update.Username != nil {
		mongoUpdate["username"] = updateUser.Username
		mongoUpdate["username_hash"] = usernameHash
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
