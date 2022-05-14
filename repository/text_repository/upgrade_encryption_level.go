package text_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (t *defaultTextRepository) upgradeEncryptionLevel(ctx context.Context, text Text) error {
	if len(t.internalEncryptionKeys) <= 0 {
		return errors.New("length of internal encryption keys is 0")
	} else if text.Id == go_block.LanguageCode_INVALID_LANGUAGE_CODE {
		return errors.New("invalid language code")
	}
	if err := t.EncryptText(actionCreate, &text); err != nil {
		return err
	}
	updateGeneralText := bson.M{}
	if text.GeneralText != nil {
		updateGeneralText = bson.M{
			"missing_password_title":   text.GeneralText.MissingPasswordTitle,
			"missing_password_details": text.GeneralText.MissingPasswordDetails,
			"missing_email_title":      text.GeneralText.MissingEmailTitle,
			"missing_email_details":    text.GeneralText.MissingEmailDetails,
			"created_by":               text.GeneralText.CreatedBy,
			"password_hint":            text.GeneralText.PasswordHint,
			"email_hint":               text.GeneralText.EmailHint,
			"error_title":              text.GeneralText.ErrorTitle,
			"error_description":        text.GeneralText.ErrorDescription,
			"no_wifi_title":            text.GeneralText.NoWifiTitle,
			"no_wifi_description":      text.GeneralText.NoWifiDescription,
		}
	}
	updateWelcomeText := bson.M{}
	if text.WelcomeText != nil {
		updateWelcomeText = bson.M{
			"welcome_title":        text.WelcomeText.WelcomeTitle,
			"welcome_details":      text.WelcomeText.WelcomeTitle,
			"continue_with_nuntio": text.WelcomeText.ContinueWithNuntio,
		}
	}
	updateRegisterText := bson.M{}
	if text.RegisterText != nil {
		updateRegisterText = bson.M{
			"register_button":               text.RegisterText.RegisterButton,
			"register_title":                text.RegisterText.RegisterTitle,
			"register_details":              text.RegisterText.RegisterDetails,
			"password_do_not_match_title":   text.RegisterText.PasswordDoNotMatchTitle,
			"password_do_not_match_details": text.RegisterText.PasswordDoNotMatchDetails,
			"repeat_password_hint":          text.RegisterText.RepeatPasswordHint,
			"contains_special_char":         text.RegisterText.ContainsSpecialChar,
			"contains_number_char":          text.RegisterText.ContainsNumberChar,
			"password_must_match":           text.RegisterText.PasswordMustMatch,
			"contains_eight_chars":          text.RegisterText.ContainsEightChars,
		}
	}
	updateLoginText := bson.M{}
	if text.LoginText != nil {
		updateLoginText = bson.M{
			"login_button":    text.LoginText.LoginButton,
			"login_title":     text.LoginText.LoginTitle,
			"login_details":   text.LoginText.LoginDetails,
			"forgot_password": text.LoginText.ForgotPassword,
		}
	}
	updateProfileText := bson.M{}
	if text.ProfileText != nil {
		updateProfileText = bson.M{
			"profile_title":               text.ProfileText.ProfileTitle,
			"logout":                      text.ProfileText.Logout,
			"change_email_title":          text.ProfileText.ChangeEmailTitle,
			"change_email_description":    text.ProfileText.ChangeEmailDescription,
			"change_password_title":       text.ProfileText.ChangePasswordTitle,
			"change_password_description": text.ProfileText.ChangePasswordDescription,
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"general_text":              updateGeneralText,
			"welcome_text":              updateWelcomeText,
			"register_text":             updateLoginText,
			"login_text":                updateRegisterText,
			"profile_text":              updateProfileText,
			"internal_encryption_level": int32(len(t.internalEncryptionKeys)),
			"updated_at":                time.Now(),
		},
	}
	if _, err := t.collection.UpdateOne(ctx, bson.M{"_id": text.Id}, mongoUpdate); err != nil {
		return err
	}
	return nil
}
