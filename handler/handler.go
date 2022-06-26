/*
	handler - the handler is the brain of this application. It has access to almost all other packages
	and uses these packages to answer API requests from the client. It is build on top of the gRPC framework.
*/
package handler

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/nuntiodev/hera/email"
	"github.com/nuntiodev/hera/models"
	"github.com/nuntiodev/hera/text"
	"io/ioutil"
	"os"
	"time"

	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/repository"
	"github.com/nuntiodev/hera/token"
	"github.com/nuntiodev/x/cryptox"
	"go.uber.org/zap"
)

var (
	accessTokenExpiry  = time.Minute * 30
	refreshTokenExpiry = time.Hour * 24 * 7
	publicKey          *rsa.PublicKey
	publicKeyString    = ""
	privateKey         *rsa.PrivateKey
	appName            = ""
)

type defaultHandler struct {
	repository         repository.Repository
	crypto             cryptox.Crypto
	token              token.Token
	email              email.Email
	text               text.Text
	emailEnabled       bool
	textEnabled        bool
	zapLog             *zap.Logger
	maxVerificationAge time.Duration
}

func decodeKeyPair(rsaPrivateKey, rsaPublicKey string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// Handle errors here
	block, _ := pem.Decode([]byte(rsaPrivateKey))
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}
	pubBlock, rest := pem.Decode([]byte(rsaPublicKey))
	if pubBlock == nil {
		return nil, nil, fmt.Errorf("pub block is nil with rest: %s", string(rest))
	}
	pubKey, err := x509.ParsePKCS1PublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}
	if privKey.PublicKey.Equal(pubKey) == false {
		return nil, nil, errors.New("keys do not match")
	}
	return privKey, &privKey.PublicKey, nil
}

func initialize() error {
	accessTokenExpiryString, ok := os.LookupEnv("ACCESS_TOKEN_EXPIRY")
	if ok {
		dur, err := time.ParseDuration(accessTokenExpiryString)
		if err == nil {
			accessTokenExpiry = dur
		}
	}
	refreshTokenExpiryString, ok := os.LookupEnv("REFRESH_TOKEN_EXPIRY")
	if ok {
		dur, err := time.ParseDuration(refreshTokenExpiryString)
		if err == nil {
			refreshTokenExpiry = dur
		}
	}
	publicKeyString, ok = os.LookupEnv("PUBLIC_KEY")
	if !ok || publicKeyString == "" {
		return errors.New("missing required PUBLIC_KEY")
	}
	privateKeyString, ok := os.LookupEnv("PRIVATE_KEY")
	if !ok || privateKeyString == "" {
		return errors.New("missing required PRIVATE_KEY")
	}
	var err error
	privateKey, publicKey, err = decodeKeyPair(privateKeyString, publicKeyString)
	if err != nil {
		return err
	}
	appName, ok = os.LookupEnv("APP_NAME")
	if !ok || appName == "" {
		appName = "Nuntio Hera App"
	}
	return nil
}

func New(zapLog *zap.Logger, repository repository.Repository, token token.Token, email email.Email, text text.Text, maxEmailVerificationAge time.Duration) (go_hera.ServiceServer, error) {
	zapLog.Info("creating handler")
	if err := initialize(); err != nil {
		return nil, err
	}
	emailEnabled := false
	if email != nil {
		emailEnabled = true
	}
	textEnabled := false
	if text != nil {
		textEnabled = true
	}
	handler := &defaultHandler{
		repository:         repository,
		token:              token,
		zapLog:             zapLog,
		email:              email,
		emailEnabled:       emailEnabled,
		text:               text,
		textEnabled:        textEnabled,
		maxVerificationAge: maxEmailVerificationAge,
	}
	if err := handler.initializeDefaultConfigAndUsers(textEnabled, emailEnabled); err != nil {
		return nil, err
	}
	return handler, nil
}

func (h *defaultHandler) initializeDefaultConfigAndUsers(textEnabled, emailEnabled bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	// INITIALIZE CONFIG
	var configUpdate *go_hera.Config
	var configCreate *go_hera.Config
	var users []*go_hera.User
	resp, err := h.GetConfig(ctx, &go_hera.HeraRequest{})
	if err != nil {
		resp, err := h.CreateNamespace(ctx, &go_hera.HeraRequest{
			Config: &go_hera.Config{Name: "Nuntio Hera App"},
		})
		configCreate = resp.Config
		if err != nil {
			return err
		}
		h.zapLog.Info("hera config was successfully created...")
	} else {
		configCreate = resp.Config
		h.zapLog.Info("hera config already exists...")
	}
	// load json file
	jsonFile, err := os.Open("hera_config.json")
	if err == nil {
		h.zapLog.Info("hera_config.json file found. Updating default config.")
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			return err
		}
		var heraConfig models.HeraConfig
		if err := json.Unmarshal(byteValue, &heraConfig); err != nil {
			return err
		}
		configUpdate, err = models.HeraConfigToProtoConfig(&heraConfig)
		if err != nil {
			return err
		}
		users = models.HeraConfigToProtoUsers(&heraConfig)
	} else {
		h.zapLog.Info("no hera_config.json file found. Create one to override default values.")
	}
	if configUpdate != nil {
		if _, err := h.UpdateConfig(ctx, &go_hera.HeraRequest{
			Config: configUpdate,
		}); err != nil {
			return err
		}
		configCreate = configUpdate
	}
	h.zapLog.Info(fmt.Sprintf("Hera is starting with config: %s", configCreate.String()))
	if configCreate.VerifyPhone && !textEnabled {
		return errors.New("default config requires phone verification, but no TextSender interfaces was provided")
	}
	if configCreate.VerifyEmail && !emailEnabled {
		return errors.New("default config requires email verification, but no EmailSender interfaces was provided")
	}
	// INITIALIZE USERS
	for _, user := range users {
		id := ""
		if user.GetId() != "" {
			id = user.GetId()
		} else if user.GetEmail() != "" {
			id = user.GetEmail()
		} else if user.GetPhone() != "" {
			id = user.GetPhone()
		} else if user.GetUsername() != "" {
			id = user.GetUsername()
		}
		// check if user already has been created
		if _, err := h.GetUser(ctx, &go_hera.HeraRequest{User: user}); err != nil {
			h.zapLog.Error("could not find user with err: " + err.Error())
			// user does not exists -> create user
			h.zapLog.Info("creating new user with: " + id)
			if _, err := h.CreateUser(ctx, &go_hera.HeraRequest{User: user}); err != nil {
				h.zapLog.Error("could not create user")
			}
		} else {
			h.zapLog.Info("user with identifier already exists: " + id)
		}
	}
	return nil
}
