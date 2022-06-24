package user_repository

import (
	"context"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (r *mongodbRepository) UpdatePassword(ctx context.Context, get *go_hera.User, update *go_hera.User) error {
	if update == nil {
		return UpdateIsNil
	} else if err := validatePassword(update.Password); err != nil {
		return err
	}
	prepare(actionGet, get)
	prepare(actionUpdatePassword, update)
	filter, err := getUserFilter(get)
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(update.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	update.Password = string(hashedPassword)
	updateUser := models.ProtoUserToUser(update)
	mongoUpdate := bson.M{
		"$set": bson.M{
			"password":                       updateUser.Password,
			"reset_password_attempts":        int32(0),
			"reset_password_code_sent_at":    time.Time{},
			"reset_password_code_expires_at": time.Time{},
			"updated_at":                     time.Now(),
		},
	}
	if _, err := r.collection.UpdateOne(
		ctx,
		filter,
		mongoUpdate,
	); err != nil {
		return err
	}
	return nil
}
