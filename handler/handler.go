package handler

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/crypto"
	"github.com/softcorp-io/block-user-service/repository"
	"go.uber.org/zap"
	"os"
	"strconv"
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
	GetStream(request *go_block.UserRequest, server go_block.UserService_GetStreamServer) error
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

func initialize() error {
	maxStreamConnectionsString, ok := os.LookupEnv("MAX_STREAM_CONNECTIONS")
	if !ok {
		return nil
	}
	count, err := strconv.Atoi(maxStreamConnectionsString)
	if err == nil {
		maxStreamConnections = count
	}
	return nil
}

func New(zapLog *zap.Logger, repository repository.Repository, crypto crypto.Crypto) (Handler, error) {
	zapLog.Info("creating handler")
	handler := &defaultHandler{
		repository: repository,
		crypto:     crypto,
		zapLog:     zapLog,
	}
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				handler.cleanupConnections()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return handler, nil
}
