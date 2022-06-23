package handler

import (
	"context"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/models"
	"github.com/nuntiodev/hera/repository/token_repository"
)

/*
	GetTokens - this method returns information about all tokens for a specific user. UserId is included in the req.Token.
*/
func (h *defaultHandler) GetTokens(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		tokenRepository token_repository.TokenRepository
		tokens          []*models.Token
	)
	tokenRepository, err = h.repository.TokenRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	tokens, err = tokenRepository.GetTokens(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{
		Tokens: models.TokensToProto(tokens),
	}, nil
}
