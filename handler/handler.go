package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/crypto"
	"github.com/softcorp-io/block-user-service/repository"
	"go.uber.org/zap"
)

type Handler interface {
	Heartbeat(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	Create(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdatePassword(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateMetadata(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateImage(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateEmail(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateOptionalId(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	UpdateSecurity(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	Get(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	GetAll(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	ValidateCredentials(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	Delete(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	DeleteBatch(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	DeleteNamespace(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
}

type defaultHandler struct {
	repository repository.Repository
	crypto     crypto.Crypto
	zapLog     *zap.Logger
}

func New(zapLog *zap.Logger, repository repository.Repository, crypto crypto.Crypto) (Handler, error) {
	zapLog.Info("creating handler")
	handler := &defaultHandler{
		repository: repository,
		crypto:     crypto,
		zapLog:     zapLog,
	}
	return handler, nil
}
