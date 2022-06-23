package token

import (
	"crypto/rsa"
	"errors"

	"github.com/golang-jwt/jwt"
	"github.com/nuntiodev/hera-proto/go_hera"
)

func (c *defaultToken) ValidateToken(key *rsa.PublicKey, jwtToken string) (*go_hera.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		jwtToken,
		&go_hera.CustomClaims{},
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
	claims, ok := token.Claims.(*go_hera.CustomClaims)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}
	if claims.Issuer != Issuer {
		return nil, errors.New("invalid issuer")
	}
	return claims, nil
}
