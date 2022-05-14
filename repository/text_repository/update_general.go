package text_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (t *defaultTextRepository) UpdateGeneral(ctx context.Context, text *go_block.Text) (*go_block.Text, error) {
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
		GeneralText: text.GeneralText,
	})
	if get.InternalEncryptionLevel > 0 {
		update.InternalEncryptionLevel = get.InternalEncryptionLevel
		if err := t.EncryptText(actionUpdate, update); err != nil {
			return nil, err
		}
	}
	updateGeneralText := bson.M{}
	if text.GeneralText != nil {
		updateGeneralText = bson.M{
			"missing_password_title":   update.GeneralText.MissingPasswordTitle,
			"missing_password_details": update.GeneralText.MissingPasswordDetails,
			"missing_email_title":      update.GeneralText.MissingEmailTitle,
			"created_by":               update.GeneralText.CreatedBy,
			"password_hint":            update.GeneralText.PasswordHint,
			"email_hint":               update.GeneralText.EmailHint,
			"error_title":              update.GeneralText.ErrorTitle,
			"error_description":        update.GeneralText.ErrorDescription,
			"no_wifi_title":            update.GeneralText.NoWifiTitle,
			"no_wifi_description":      update.GeneralText.NoWifiDescription,
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"general_text": updateGeneralText,
			"updated_at":   time.Now(),
		},
	}
	if _, err := t.collection.UpdateOne(ctx, bson.M{"_id": text.Id}, mongoUpdate); err != nil {
		return nil, err
	}
	// set updated fields
	get.RegisterText = text.RegisterText
	return get, nil
}
