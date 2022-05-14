package text_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (t *defaultTextRepository) UpdateRegister(ctx context.Context, text *go_block.Text) (*go_block.Text, error) {
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
		RegisterText: text.RegisterText,
	})
	if get.InternalEncryptionLevel > 0 {
		update.InternalEncryptionLevel = get.InternalEncryptionLevel
		if err := t.EncryptText(actionUpdate, update); err != nil {
			return nil, err
		}
	}
	updateRegisterText := bson.M{}
	if text.RegisterText != nil {
		updateRegisterText = bson.M{
			"register_button":               update.RegisterText.RegisterButton,
			"register_title":                update.RegisterText.RegisterTitle,
			"register_details":              update.RegisterText.RegisterDetails,
			"password_do_not_match_title":   update.RegisterText.PasswordDoNotMatchTitle,
			"password_do_not_match_details": update.RegisterText.PasswordDoNotMatchDetails,
			"repeat_password_hint":          update.RegisterText.RepeatPasswordHint,
			"contains_special_char":         update.RegisterText.ContainsSpecialChar,
			"contains_number_char":          update.RegisterText.ContainsNumberChar,
			"password_must_match":           update.RegisterText.PasswordMustMatch,
			"contains_eight_chars":          update.RegisterText.ContainsEightChars,
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"register_text": updateRegisterText,
			"updated_at":    time.Now(),
		},
	}
	if _, err := t.collection.UpdateOne(ctx, bson.M{"_id": text.Id}, mongoUpdate); err != nil {
		return nil, err
	}
	// set updated fields
	get.RegisterText = text.RegisterText
	return get, nil
}
