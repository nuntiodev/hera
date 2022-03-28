package crypto

import (
	"errors"
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
	"time"
)

func (c *defaultCrypto) GenerateToken(userId, refreshTokenId, tokenType string, expiresAt time.Duration) (string, *CustomClaims, error) {
	expiresAtInt64 := int64(0)
	if expiresAt.Seconds() != 0 {
		expiresAtInt64 = time.Now().UTC().Add(expiresAt).Unix()
	}
	if tokenType != TokenTypeAccess && tokenType != TokenTypeRefresh {
		return "", nil, errors.New("invalid token type")
	}
	if tokenType == TokenTypeAccess && refreshTokenId == "" {
		return "", nil, errors.New("missing required refreshTokenId")
	}
	claims := CustomClaims{
		UserId:         userId,
		Type:           tokenType,
		RefreshTokenId: refreshTokenId,
		StandardClaims: jwt.StandardClaims{
			Id: uuid.NewV4().String(),
			// In JWT, the expiry time is expressed as unix
			ExpiresAt: expiresAtInt64,
			Issuer:    Issuer,
		},
	}
	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM(c.jwtPrivateKey)
	if err != nil {
		return "", nil, err
	}
	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(signingKey)
	if err != nil {
		return "", nil, err
	}
	return signedToken, &claims, nil
}
