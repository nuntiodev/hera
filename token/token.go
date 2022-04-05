package token

import (
	"crypto/rsa"
	"github.com/softcorp-io/block-proto/go_block"
	"time"
)

const (
	TokenTypeAccess  = "access_token"
	TokenTypeRefresh = "refresh_token"
	Issuer           = "Block User Service"
)

type Token interface {
	GenerateToken(privateKey *rsa.PrivateKey, userId, refreshTokenId, tokenType string, expiresAt time.Duration) (string, *go_block.CustomClaims, error)
	ValidateToken(publicKey *rsa.PublicKey, jwtToken string) (*go_block.CustomClaims, error)
}

type defaultToken struct{}

func New() (Token, error) {
	return &defaultToken{}, nil
}
