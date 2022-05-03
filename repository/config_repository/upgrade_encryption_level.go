package config_repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (cr *defaultConfigRepository) upgradeEncryptionLevel(ctx context.Context, config Config) error {
	if err := cr.EncryptConfig(actionCreate, &config); err != nil {
		return err
	}
	updateGeneralText := bson.M{}
	if config.GeneralText != nil {
		updateGeneralText = bson.M{
			"missing_password_title":   config.GeneralText.MissingPasswordTitle,
			"missing_password_details": config.GeneralText.MissingPasswordDetails,
			"missing_email_title":      config.GeneralText.MissingEmailTitle,
			"missing_email_details":    config.GeneralText.MissingEmailDetails,
			"created_by":               config.GeneralText.CreatedBy,
			"password_hint":            config.GeneralText.PasswordHint,
			"email_hint":               config.GeneralText.EmailHint,
			"error_title":              config.GeneralText.ErrorTitle,
			"error_description":        config.GeneralText.ErrorDescription,
			"no_wifi_title":            config.GeneralText.NoWifiTitle,
			"no_wifi_description":      config.GeneralText.NoWifiDescription,
		}
	}
	updateWelcomeText := bson.M{}
	if config.WelcomeText != nil {
		updateWelcomeText = bson.M{
			"welcome_title":        config.WelcomeText.WelcomeTitle,
			"welcome_details":      config.WelcomeText.WelcomeTitle,
			"continue_with_nuntio": config.WelcomeText.ContinueWithNuntio,
			"logo":                 config.WelcomeText.Logo,
		}
	}
	updateRegisterText := bson.M{}
	if config.RegisterText != nil {
		updateRegisterText = bson.M{
			"register_button":               config.RegisterText.RegisterButton,
			"register_title":                config.RegisterText.RegisterTitle,
			"register_details":              config.RegisterText.RegisterDetails,
			"password_do_not_match_title":   config.RegisterText.PasswordDoNotMatchTitle,
			"password_do_not_match_details": config.RegisterText.PasswordDoNotMatchDetails,
			"repeat_password_hint":          config.RegisterText.RepeatPasswordHint,
			"contains_special_char":         config.RegisterText.ContainsSpecialChar,
			"contains_number_char":          config.RegisterText.ContainsNumberChar,
			"password_must_match":           config.RegisterText.PasswordMustMatch,
			"contains_eight_chars":          config.RegisterText.ContainsEightChars,
		}
	}
	updateLoginText := bson.M{}
	if config.LoginText != nil {
		updateLoginText = bson.M{
			"login_button":    config.LoginText.LoginButton,
			"login_title":     config.LoginText.LoginTitle,
			"login_details":   config.LoginText.LoginDetails,
			"forgot_password": config.LoginText.ForgotPassword,
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"name":                      config.Name,
			"website":                   config.Website,
			"about":                     config.About,
			"email":                     config.Email,
			"logo":                      config.Logo,
			"general_text":              updateGeneralText,
			"welcome_text":              updateWelcomeText,
			"register_text":             updateLoginText,
			"login_text":                updateRegisterText,
			"internal_encryption_level": int32(len(cr.internalEncryptionKeys)),
			"updated_at":                time.Now(),
		},
	}
	if _, err := cr.collection.UpdateOne(ctx, bson.M{"_id": config.Id}, mongoUpdate); err != nil {
		return err
	}
	return nil
}
