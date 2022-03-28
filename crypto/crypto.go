package crypto

import (
	"github.com/softcorp-io/block-proto/go_block"
	"time"
)

const (
	TokenTypeAccess  = "access_token"
	TokenTypeRefresh = "refresh_token"
	Issuer           = "Block User Service"
)

type Crypto interface {
	Encrypt(stringToEncrypt string, keyString string) (string, error)
	Decrypt(encryptedString string, keyString string) (string, error)
	EncryptUser(key string, user *go_block.User) error
	DecryptUser(key string, user *go_block.User) error
	GenerateToken(userId, refreshTokenId, tokenType string, expiresAt time.Duration) (string, *go_block.CustomClaims, error)
	ValidateToken(jwtToken string) (*go_block.CustomClaims, error)
}

type defaultCrypto struct {
	jwtPrivateKey []byte
	jwtPublicKey  []byte
}

func New(jwtPrivateKey, jwtPublicKey []byte) (Crypto, error) {
	return &defaultCrypto{
		jwtPrivateKey: jwtPrivateKey,
		jwtPublicKey:  jwtPublicKey,
	}, nil
}
