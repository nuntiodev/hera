package handler

import (
	"context"
	"errors"

	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/repository/token_repository"
)

/*
	BlockToken will block an access or refresh token in the database for a specific login-session.
	You can either provide a token which will then be blocked or a pointer to that token, which will then be blocked.
*/
func (h *defaultHandler) BlockToken(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		claims          *go_hera.CustomClaims
		token           *go_hera.Token
		tokenRepository token_repository.TokenRepository
	)
	if req.GetToken() != nil {
		token = req.GetToken()
	} else if req.GetTokenPointer() != "" {
		// validate requested token and get id of the token
		claims, err = h.token.ValidateToken(publicKey, req.GetTokenPointer())
		if err != nil {
			return nil, err
		}
		token = &go_hera.Token{
			UserId: claims.UserId,
			Id:     claims.Id,
		}
	} else {
		return nil, errors.New("not a valid token")
	}
	// validate if token is blocked in db
	tokenRepository, err = h.repository.TokenRepositoryBuilder().SetNamespace(req.GetNamespace()).Build(ctx)
	if err != nil {
		return nil, err
	}
	if err = tokenRepository.Block(ctx, token); err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{}, nil
}
