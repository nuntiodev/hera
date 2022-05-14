package text_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (t *defaultTextRepository) UpdateLogin(ctx context.Context, text *go_block.Text) (*go_block.Text, error) {
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
		LoginText: text.LoginText,
	})
	if get.InternalEncryptionLevel > 0 {
		update.InternalEncryptionLevel = get.InternalEncryptionLevel
		if err := t.EncryptText(actionUpdate, update); err != nil {
			return nil, err
		}
	}
	updateLoginText := bson.M{}
	if text.LoginText != nil {
		updateLoginText = bson.M{
			"login_button":    update.LoginText.LoginButton,
			"login_title":     update.LoginText.LoginTitle,
			"login_details":   update.LoginText.LoginDetails,
			"forgot_password": update.LoginText.ForgotPassword,
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"login_text": updateLoginText,
			"updated_at": time.Now(),
		},
	}
	if _, err := t.collection.UpdateOne(ctx, bson.M{"_id": text.Id}, mongoUpdate); err != nil {
		return nil, err
	}
	// set updated fields
	get.RegisterText = text.RegisterText
	return get, nil
}
