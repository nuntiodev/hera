package handler

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/repository"
	"github.com/nuntiodev/nuntio-user-block/token"
	"github.com/nuntiodev/x/cryptox"
	"go.uber.org/zap"
)

var (
	accessTokenExpiry  = time.Minute * 30
	refreshTokenExpiry = time.Hour * 24 * 30
	publicKey          *rsa.PublicKey
	publicKeyString    = ""
	privateKey         *rsa.PrivateKey
	validatePassword   = false
)

type Handler interface {
	Heartbeat(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	Create(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdatePassword(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateMetadata(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateImage(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateEmail(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateName(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateBirthdate(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateOptionalId(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateSecurity(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	Get(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	GetAll(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	ValidateCredentials(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	Login(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	ValidateToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	BlockToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	BlockTokenById(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	RefreshToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	GetTokens(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	PublicKeys(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	RecordActiveMeasurement(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UserActiveHistory(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	NamespaceActiveHistory(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	Delete(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	DeleteBatch(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	DeleteNamespace(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	CreateNamespaceConfig(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateConfigSettings(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateConfigDetails(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateConfigGeneralText(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateConfigWelcomeText(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateConfigRegisterText(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateConfigLoginText(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)

	GetConfig(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	DeleteConfig(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
}

type defaultHandler struct {
	repository repository.Repository
	crypto     cryptox.Crypto
	token      token.Token
	zapLog     *zap.Logger
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
	validatePasswordString, ok := os.LookupEnv("VALIDATE_PASSWORD")
	if ok && validatePasswordString != "" {
		// on error it returns false
		validatePassword, _ = strconv.ParseBool(validatePasswordString)
	}
	return nil
}

func New(zapLog *zap.Logger, repository repository.Repository, crypto cryptox.Crypto, token token.Token) (Handler, error) {
	zapLog.Info("creating handler")
	if err := initialize(); err != nil {
		return nil, err
	}
	handler := &defaultHandler{
		repository: repository,
		crypto:     crypto,
		token:      token,
		zapLog:     zapLog,
	}
	return handler, nil
}
