package text_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (t *defaultTextRepository) UpdateWelcome(ctx context.Context, text *go_block.Text) (*go_block.Text, error) {
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
		WelcomeText: text.WelcomeText,
	})
	if get.InternalEncryptionLevel > 0 {
		update.InternalEncryptionLevel = get.InternalEncryptionLevel
		if err := t.EncryptText(actionUpdate, update); err != nil {
			return nil, err
		}
	}
	updateWelcomeText := bson.M{}
	if text.WelcomeText != nil {
		updateWelcomeText = bson.M{
			"welcome_title":        update.WelcomeText.WelcomeTitle,
			"welcome_details":      update.WelcomeText.WelcomeDetails,
			"continue_with_nuntio": update.WelcomeText.ContinueWithNuntio,
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"welcome_text": updateWelcomeText,
			"updated_at":   time.Now(),
		},
	}
	if _, err := t.collection.UpdateOne(ctx, bson.M{"_id": text.Id.String()}, mongoUpdate); err != nil {
		return nil, err
	}
	// set updated fields
	get.WelcomeText = text.WelcomeText
	return get, nil
}
