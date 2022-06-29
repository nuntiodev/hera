package token

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/nuntiodev/hera-sdks/go_hera"
)

func (c *defaultToken) GenerateToken(privateKey *rsa.PrivateKey, tokenId, userId, refreshTokenId, tokenType string, expiresAt time.Duration) (string, *go_hera.CustomClaims, error) {
	expiresAtInt64 := int64(0)
	if expiresAt.Seconds() != 0 {
		expiresAtInt64 = time.Now().Add(expiresAt).Unix()
	} else if tokenId == "" {
		return "", nil, errors.New("missing required token id")
	}
	if tokenType != AccessToken && tokenType != RefreshToken {
		return "", nil, errors.New("invalid token type")
	}
	if tokenType == AccessToken && refreshTokenId == "" {
		return "", nil, errors.New("missing required refreshTokenId")
	}
	claims := go_hera.CustomClaims{
		UserId:         userId,
		Type:           tokenType,
		RefreshTokenId: refreshTokenId,
		StandardClaims: jwt.StandardClaims{
			Id: tokenId,
			// In JWT, the expiry time is expressed as unix
			ExpiresAt: expiresAtInt64,
			Issuer:    Issuer,
		},
	}
	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	if err != nil {
		return "", nil, err
	}
	return signedToken, &claims, nil
}
