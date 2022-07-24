package user_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/models"
	"go.mongodb.org/mongo-driver/bson"
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
	if r.hasher == nil {
		return errors.New("hasher is nil")
	}
	hashedPassword, err := r.hasher.Generate(update.Password.Body)
	if err != nil {
		return err
	}
	update.Password = hashedPassword
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
