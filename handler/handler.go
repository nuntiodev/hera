package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/crypto"
	"github.com/softcorp-io/block-user-service/repository"
	"go.uber.org/zap"
	"time"
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
	Login(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	ValidateToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	BlockToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	RefreshToken(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	PublicKeys(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	Delete(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	DeleteBatch(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
	DeleteNamespace(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
}

type defaultHandler struct {
	repository        repository.Repository
	accessTokenExpiry time.Duration
	crypto            crypto.Crypto
	zapLog            *zap.Logger
	jwtPublicKey      []byte
}

func New(zapLog *zap.Logger, repository repository.Repository, crypto crypto.Crypto, accessTokenExpiry time.Duration, jwtPublicKey []byte) (Handler, error) {
	zapLog.Info("creating handler")
	handler := &defaultHandler{
		repository:        repository,
		accessTokenExpiry: accessTokenExpiry,
		crypto:            crypto,
		zapLog:            zapLog,
		jwtPublicKey:      jwtPublicKey,
	}
	return handler, nil
}
