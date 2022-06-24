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
	accessTokenExpiry              = time.Minute * 30
	refreshTokenExpiry             = time.Hour * 24 * 7
	publicKey                      *rsa.PublicKey
	publicKeyString                = ""
	privateKey                     *rsa.PrivateKey
	emailVerificationTemplatePath  = ""
	emailResetPasswordTemplatePath = ""
)

type Handler interface {
	Heartbeat(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	CreateUser(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	UpdateUserMetadata(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	UpdateUserProfile(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	UpdateUserContact(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	UpdateUserPassword(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	SearchForUser(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	GetUser(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	GetUsers(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error)
	ListUsers(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	ValidateCredentials(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	Login(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	DeleteUser(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	DeleteUsers(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	CreateTokenPair(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	ValidateToken(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	BlockToken(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	RefreshToken(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	GetTokens(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	PublicKeys(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	SendVerificationEmail(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	VerifyEmail(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	SendVerificationText(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	VerifyPhone(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	SendResetPasswordEmail(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	SendResetPasswordText(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	ResetPassword(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	DeleteNamespace(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	CreateNamespace(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	RegisterRsaKey(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	RemovePublicKey(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	GetConfig(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	UpdateConfig(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
	DeleteConfig(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
}

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
	return nil
}

func (h *defaultHandler) initializeDefaultConfig() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	_, err := h.GetConfig(ctx, &go_hera.HeraRequest{})
	heraConfig := models.Config{}
	if err != nil {
		h.zapLog.Info("hera config does not exits - creating...")
		// load json file
		jsonFile, err := os.Open("hera_config.json")
		if err == nil {
			byteValue, err := ioutil.ReadAll(jsonFile)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(byteValue, &heraConfig); err != nil {
				return err
			}
		}
		if _, err := h.CreateNamespace(ctx, &go_hera.HeraRequest{
			Config: models.ConfigToProtoConfig(&heraConfig),
		}); err != nil {
			return err
		}
		h.zapLog.Info("hera config was successfully created...")
	} else {
		h.zapLog.Info("hera config already exists...")
	}
	if _, err := h.UpdateConfig(ctx, &go_hera.HeraRequest{
		Config: models.ConfigToProtoConfig(&heraConfig),
	}); err != nil {
		return err
	}
	return nil
}

func New(zapLog *zap.Logger, repository repository.Repository, token token.Token, email email.Email, text text.Text, maxEmailVerificationAge time.Duration) (Handler, error) {
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
	if err := handler.initializeDefaultConfig(); err != nil {
		return nil, err
	}
	return handler, nil
}
