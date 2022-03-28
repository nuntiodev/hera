package crypto

import (
	"github.com/golang-jwt/jwt"
	"github.com/softcorp-io/block-proto/go_block"
	"time"
)

const (
	TokenTypeAccess  = "access_token"
	TokenTypeRefresh = "refresh_token"
	Issuer           = "Block User Service"
)

type CustomClaims struct {
	UserId string `json:"user_id"`
	Type   string `json:"type"`
	jwt.StandardClaims
}

type Crypto interface {
	Encrypt(stringToEncrypt string, keyString string) (string, error)
	Decrypt(encryptedString string, keyString string) (string, error)
	EncryptUser(key string, user *go_block.User) error
	DecryptUser(key string, user *go_block.User) error
	GenerateToken(userId, tokenType string, expiresAt time.Duration) (string, error)
	ValidateToken(jwtToken string) (*CustomClaims, error)
}

type defaultCrypto struct {
	jwtPrivateKey []byte `json:"jwt_private_key"`
	jwtPublicKey  []byte `json:"jwt_public_key"`
}

func New(jwtPrivateKey, jwtPublicKey []byte) (Crypto, error) {
	return &defaultCrypto{
		jwtPrivateKey: jwtPrivateKey,
		jwtPublicKey:  jwtPublicKey,
	}, nil
}
