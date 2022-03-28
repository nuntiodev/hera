package crypto

import (
	"errors"
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
	"time"
)

func (c *defaultCrypto) GenerateToken(userId, tokenType string, expiresAt time.Duration) (string, error) {
	expiresAtInt64 := int64(0)
	if expiresAt.Seconds() != 0 {
		expiresAtInt64 = time.Now().UTC().Add(expiresAt).Unix()
	}
	if tokenType != TokenTypeAccess && tokenType != TokenTypeRefresh {
		return "", errors.New("invalid token type")
	}
	claims := CustomClaims{
		UserId: userId,
		Type:   tokenType,
		StandardClaims: jwt.StandardClaims{
			Id: uuid.NewV4().String(),
			// In JWT, the expiry time is expressed as unix
			ExpiresAt: expiresAtInt64,
			Issuer:    Issuer,
		},
	}
	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM(c.jwtPrivateKey)
	if err != nil {
		return "", err
	}
	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
