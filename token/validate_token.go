package token

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/softcorp-io/block-proto/go_block"
)

func (c *defaultToken) ValidateToken(key *rsa.PublicKey, jwtToken string) (*go_block.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		jwtToken,
		&go_block.CustomClaims{},
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
	claims, ok := token.Claims.(*go_block.CustomClaims)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}
	if claims.Issuer != Issuer {
		return nil, errors.New("invalid issuer")
	}
	return claims, nil
}
