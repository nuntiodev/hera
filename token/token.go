package token

import (
	"crypto/rsa"
	"time"

	"github.com/nuntiodev/hera-sdks/go_hera"
)

const (
	AccessToken  = "access_token"
	RefreshToken = "refresh_token"
	Issuer       = "Block User Service"
)

type Token interface {
	GenerateToken(privateKey *rsa.PrivateKey, tokenId, userId, refreshTokenId, tokenType string, expiresAt time.Duration) (string, *go_hera.CustomClaims, error)
	ValidateToken(publicKey *rsa.PublicKey, jwtToken string) (*go_hera.CustomClaims, error)
}

type defaultToken struct{}

func New() (Token, error) {
	return &defaultToken{}, nil
}
