package handler

import (
	"context"
	"errors"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/repository"
	"github.com/softcorp-io/block-user-service/repository/user_repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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
	DeleteNamespace(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error)
}

type defaultHandler struct {
	repository *repository.Repository
	zapLog     *zap.Logger
}

func New(zapLog *zap.Logger, repository *repository.Repository) (Handler, error) {
	zapLog.Info("creating handler")
	handler := &defaultHandler{
		repository: repository,
		zapLog:     zapLog,
	}
	return handler, nil
}

func (h *defaultHandler) Heartbeat(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	if err := h.repository.Liveness(ctx); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{}, nil
}

func (h *defaultHandler) Create(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	createdUser, err := h.repository.UserRepository.Create(ctx, req.User, &user_repository.EncryptionOptions{
		Key: req.EncryptionKey,
	})
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{
		User: createdUser,
	}, nil
}

func (h *defaultHandler) UpdatePassword(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	updatedUser, err := h.repository.UserRepository.UpdatePassword(ctx, req.User, req.Update)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{
		User: updatedUser,
	}, nil
}

func (h *defaultHandler) UpdateMetadata(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	updatedUser, err := h.repository.UserRepository.UpdateMetadata(ctx, req.User, req.Update, &user_repository.EncryptionOptions{
		Key: req.EncryptionKey,
	})
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		User: updatedUser,
	}, nil
}

func (h *defaultHandler) UpdateImage(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	updatedUser, err := h.repository.UserRepository.UpdateImage(ctx, req.User, req.Update, &user_repository.EncryptionOptions{
		Key: req.EncryptionKey,
	})
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		User: updatedUser,
	}, nil
}

func (h *defaultHandler) UpdateEmail(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	updatedUser, err := h.repository.UserRepository.UpdateEmail(ctx, req.User, req.Update, &user_repository.EncryptionOptions{
		Key: req.EncryptionKey,
	})
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		User: updatedUser,
	}, nil
}

func (h *defaultHandler) UpdateOptionalId(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	updatedUser, err := h.repository.UserRepository.UpdateOptionalId(ctx, req.User, req.Update)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		User: updatedUser,
	}, nil
}

func (h *defaultHandler) UpdateSecurity(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	updatedUser, err := h.repository.UserRepository.UpdateSecurity(ctx, req.User, req.Update, &user_repository.EncryptionOptions{
		Key: req.EncryptionKey,
	})
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		User: updatedUser,
	}, nil
}

func (h *defaultHandler) Get(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	getUser, err := h.repository.UserRepository.Get(ctx, req.User, &user_repository.EncryptionOptions{
		Key: req.EncryptionKey,
	})
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		User: getUser,
	}, nil
}

func (h *defaultHandler) GetAll(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	getUsers, err := h.repository.UserRepository.GetAll(ctx, req.Filter, req.Namespace, &user_repository.EncryptionOptions{
		Key: req.EncryptionKey,
	})
	if err != nil {
		return nil, err
	}
	usersInNamespace, err := h.repository.UserRepository.Count(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}
	return &go_block.UserResponse{
		Users:      getUsers,
		UsersAmout: usersInNamespace,
	}, nil
}

func (h *defaultHandler) ValidateCredentials(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	resp, err := h.Get(ctx, req)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	if resp.User.Password == "" {
		return &go_block.UserResponse{}, errors.New("please update the user with a non-empty password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(resp.User.Password), []byte(req.User.Password)); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{
		User: resp.User,
	}, nil
}

func (h *defaultHandler) Delete(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	return &go_block.UserResponse{}, h.repository.UserRepository.Delete(ctx, req.User)
}

func (h *defaultHandler) DeleteNamespace(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	return &go_block.UserResponse{}, h.repository.UserRepository.DeleteNamespace(ctx, req.Namespace)
}
