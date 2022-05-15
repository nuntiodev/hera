package text_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (t *defaultTextRepository) UpdateProfile(ctx context.Context, text *go_block.Text) (*go_block.Text, error) {
	if text == nil {
		return nil, errors.New("missing required text")
	} else if text.Id == go_block.LanguageCode_INVALID_LANGUAGE_CODE {
		return nil, errors.New("missing required text language code id")
	}
	get, err := t.Get(ctx, text.Id)
	if err != nil {
		return nil, err
	}
	update := ProtoTextToText(&go_block.Text{
		ProfileText: text.ProfileText,
	})
	if get.InternalEncryptionLevel > 0 {
		update.InternalEncryptionLevel = get.InternalEncryptionLevel
		if err := t.EncryptText(actionUpdate, update); err != nil {
			return nil, err
		}
	}
	updateProfileText := bson.M{}
	if text.ProfileText != nil {
		updateProfileText = bson.M{
			"profile_title":               update.ProfileText.ProfileTitle,
			"logout":                      update.ProfileText.Logout,
			"change_email_title":          update.ProfileText.ChangeEmailTitle,
			"change_email_description":    update.ProfileText.ChangeEmailDescription,
			"change_password_title":       update.ProfileText.ChangePasswordTitle,
			"change_password_description": update.ProfileText.ChangePasswordDescription,
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"profile_text": updateProfileText,
			"updated_at":   time.Now(),
		},
	}
	if _, err := t.collection.UpdateOne(ctx, bson.M{"_id": text.Id.String()}, mongoUpdate); err != nil {
		return nil, err
	}
	// set updated fields
	get.ProfileText = text.ProfileText
	return get, nil
}
