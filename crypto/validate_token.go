package crypto

import (
	"errors"
	"github.com/golang-jwt/jwt"
)

func (c *defaultCrypto) ValidateToken(jwtToken string) (*CustomClaims, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(c.jwtPublicKey)
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(
		jwtToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return key, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if token.Valid == false {
		return nil, errors.New("token is not valid")
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}
	if claims.Issuer != Issuer {
		return nil, errors.New("invalid issuer")
	}
	return claims, nil
}
