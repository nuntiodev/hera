package handler

import (
	"context"
	"errors"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"github.com/softcorp-io/block-user-service/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Handler interface {
	Heartbeat(ctx context.Context, req *block_user.Request) (*block_user.Response, error)
	Create(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error)
	UpdatePassword(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error)
	UpdateEmail(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error)
	UpdateProfile(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error)
	UpdateNamespace(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error)
	Get(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error)
	GetAll(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error)
	Search(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error)
	ValidateCredentials(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error)
	Delete(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error)
	DeleteNamespace(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error)
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

func (h *defaultHandler) Heartbeat(ctx context.Context, req *block_user.Request) (*block_user.Response, error) {
	if err := h.repository.Liveness(ctx); err != nil {
		return &block_user.Response{}, err
	}
	return &block_user.Response{}, nil
}

func (h *defaultHandler) Create(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error) {
	createdUser, err := h.repository.UserRepository.Create(ctx, req.User)
	if err != nil {
		return &block_user.UserResponse{}, err
	}
	return &block_user.UserResponse{
		User: createdUser,
	}, nil
}

func (h *defaultHandler) UpdatePassword(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error) {
	updatedUser, err := h.repository.UserRepository.UpdatePassword(ctx, req.Update)
	if err != nil {
		return &block_user.UserResponse{}, err
	}
	return &block_user.UserResponse{
		User: updatedUser,
	}, nil
}

func (h *defaultHandler) UpdateEmail(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error) {
	updatedUser, err := h.repository.UserRepository.UpdateEmail(ctx, req.Update)
	if err != nil {
		return nil, err
	}
	return &block_user.UserResponse{
		User: updatedUser,
	}, nil
}

func (h *defaultHandler) UpdateProfile(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error) {
	updatedUser, err := h.repository.UserRepository.UpdateProfile(ctx, req.Update)
	if err != nil {
		return nil, err
	}
	return &block_user.UserResponse{
		User: updatedUser,
	}, nil
}

func (h *defaultHandler) UpdateNamespace(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error) {
	updatedUser, err := h.repository.UserRepository.UpdateNamespace(ctx, req.Update)
	if err != nil {
		return nil, err
	}
	return &block_user.UserResponse{
		User: updatedUser,
	}, nil
}

func (h *defaultHandler) Get(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error) {
	if req.User.Id != "" {
		getUser, err := h.repository.UserRepository.GetById(ctx, req.User)
		if err != nil {
			return nil, err
		}
		return &block_user.UserResponse{
			User: getUser,
		}, nil
	} else if req.User.Email != "" {
		getUser, err := h.repository.UserRepository.GetByEmail(ctx, req.User)
		if err != nil {
			return nil, err
		}
		return &block_user.UserResponse{
			User: getUser,
		}, nil
	}
	return &block_user.UserResponse{}, errors.New("missing required search parameter (mail or id)")
}

func (h *defaultHandler) GetAll(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error) {
	getUsers, err := h.repository.UserRepository.GetAll(ctx, req.Filter, req.Namespace)
	if err != nil {
		return nil, err
	}
	return &block_user.UserResponse{
		Users: getUsers,
	}, nil
}

func (h *defaultHandler) Search(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error) {
	getUsers, err := h.repository.UserRepository.Search(ctx, req.Search, req.Namespace)
	if err != nil {
		return nil, err
	}
	return &block_user.UserResponse{
		Users: getUsers,
	}, nil
}

func (h *defaultHandler) ValidateCredentials(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error) {
	resp, err := h.Get(ctx, req)
	if err != nil {
		return &block_user.UserResponse{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(resp.User.Password), []byte(req.User.Password)); err != nil {
		return &block_user.UserResponse{}, err
	}
	return &block_user.UserResponse{}, nil
}

func (h *defaultHandler) Delete(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error) {
	return &block_user.UserResponse{}, h.repository.UserRepository.Delete(ctx, req.User)
}

func (h *defaultHandler) DeleteNamespace(ctx context.Context, req *block_user.UserRequest) (*block_user.UserResponse, error) {
	return &block_user.UserResponse{}, h.repository.UserRepository.DeleteNamespace(ctx, req.Namespace)
}
