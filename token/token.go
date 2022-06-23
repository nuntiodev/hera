package token

import (
	"crypto/rsa"
	"time"

	"github.com/nuntiodev/hera-proto/go_hera"
)

const (
	TokenTypeAccess  = "access_token"
	TokenTypeRefresh = "refresh_token"
	Issuer           = "Block User Service"
)

type Token interface {
	GenerateToken(privateKey *rsa.PrivateKey, tokenId, userId, refreshTokenId, tokenType string, expiresAt time.Duration) (string, *go_hera.CustomClaims, error)
	ValidateToken(publicKey *rsa.PublicKey, jwtToken string) (*go_hera.CustomClaims, error)
}

type defaultToken struct{}

func New() (Token, error) {
	return &defaultToken{}, nil
}
