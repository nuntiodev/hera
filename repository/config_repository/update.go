package config_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/models"
	"github.com/nuntiodev/x/randx"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (c *defaultConfigRepository) Update(ctx context.Context, config *go_hera.Config) error {
	if config == nil {
		return errors.New("missing required config")
	} else if config.Name == "" {
		config.Name = "Hera App"
	}
	if config.Scrypt == nil {
		signerKey, err := randx.GenerateRandomString(88)
		if err != nil {
			return err
		}
		config.Scrypt = &go_hera.Scrypt{SignerKey: signerKey}
	}
	if config.Scrypt.SignerKey == "" {
		return errors.New("missing required signer key for scrypt config")
	}
	updateConfig := models.ProtoConfigToConfig(&go_hera.Config{
		Name:                     config.Name,
		Logo:                     config.Logo,
		DisableLogin:             config.DisableLogin,
		DisableSignup:            config.DisableSignup,
		ValidatePassword:         config.ValidatePassword,
		VerifyEmail:              config.VerifyEmail,
		VerifyPhone:              config.VerifyPhone,
		Scrypt:                   config.Scrypt,
		Bcrypt:                   config.Bcrypt,
		HasingAlgorithm:          config.HasingAlgorithm,
		SupportedLoginMechanisms: config.SupportedLoginMechanisms,
	})
	if err := c.crypto.Encrypt(updateConfig); err != nil {
		return err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"name":                       updateConfig.Name,
			"logo":                       updateConfig.Logo,
			"nuntio_verify_id":           updateConfig.NuntioVerifyId,
			"disable_signup":             updateConfig.DisableSignup,
			"disable_login":              updateConfig.DisableLogin,
			"validate_password":          updateConfig.ValidatePassword,
			"verify_email":               updateConfig.VerifyEmail,
			"verify_phone":               updateConfig.VerifyPhone,
			"supported_login_mechanisms": updateConfig.SupportedLoginMechanisms,
			"scrypt":                     updateConfig.Scrypt,
			"bcrypt":                     updateConfig.Bcrypt,
			"hashing_algorithm":          updateConfig.HashingAlgorithm,
			"updated_at":                 time.Now(),
		},
	}
	if _, err := c.collection.UpdateOne(ctx, bson.M{"_id": namespaceConfigName}, mongoUpdate); err != nil {
		return err
	}
	return nil
}
