package interceptor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nuntiodev/block-proto/go_block"
	"google.golang.org/grpc"
)

const (
	ProjectName         = "/BlockUser.UserService/"
	Heartbeat           = "Heartbeat"
	Create              = "Create"
	UpdatePassword      = "UpdatePassword"
	UpdateMetadata      = "UpdateMetadata"
	UpdateEmail         = "UpdateEmail"
	UpdateOptionalId    = "UpdateOptionalId"
	UpdateImage         = "UpdateImage"
	UpdateSecurity      = "UpdateSecurity"
	Get                 = "Get"
	GetAll              = "GetAll"
	ValidateCredentials = "ValidateCredentials"
	Login               = "Login"
	ValidateToken       = "ValidateToken"
	BlockToken          = "BlockToken"
	BlockTokenById      = "BlockTokenById"
	RefreshToken        = "RefreshToken"
	GetTokens           = "GetTokens"
	PublicKeys          = "PublicKeys"
	Delete              = "Delete"
	DeleteNamespace     = "DeleteNamespace"
	DeleteBatch         = "DeleteBatch"
)

var (
	TokenIsNil          = errors.New("token is nil")
	TokenPointerIsEmpty = errors.New("token pointer is nil")
	UpdateIsNil         = errors.New("update is nil")
	UserIsNil           = errors.New("user is nil")
	NamespaceIsEmpty    = errors.New("namespace is empty")
	UserBatchIsNil      = errors.New("user batch is nil")
)

func (i *DefaultInterceptor) WithValidateUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info == nil {
		return nil, errors.New("invalid request")
	}
	translatedReq, ok := req.(*go_block.UserRequest)
	if !ok {
		translatedReq = &go_block.UserRequest{}
	}
	method := strings.Split(info.FullMethod, ProjectName)
	if len(method) != 2 {
		return nil, errors.New(fmt.Sprintf("invalid method call: %s", info.FullMethod))
	}
	switch method[1] {
	case Heartbeat, GetAll, PublicKeys:
		break
	case BlockToken:
		if translatedReq.TokenPointer == "" {
			return nil, TokenPointerIsEmpty
		}
	case Create, Get:
		if translatedReq.User == nil {
			return &go_block.UserResponse{}, UserIsNil
		}
	case UpdatePassword, UpdateMetadata,
		UpdateImage, UpdateEmail, UpdateOptionalId:
		if translatedReq.Update == nil {
			return &go_block.UserResponse{}, UpdateIsNil
		} else if translatedReq.User == nil {
			return &go_block.UserResponse{}, UpdateIsNil
		}
	case ValidateCredentials, Delete, UpdateSecurity, Login:
		if translatedReq.User == nil {
			return nil, UserIsNil
		}
	case ValidateToken, RefreshToken, GetTokens, BlockTokenById:
		if translatedReq.Token == nil {
			return nil, TokenIsNil
		}
	case DeleteNamespace:
		if translatedReq.Namespace == "" {
			return nil, NamespaceIsEmpty
		}
	case DeleteBatch:
		if translatedReq.UserBatch == nil {
			return nil, UserBatchIsNil
		}
	default:
		return &go_block.UserResponse{}, errors.New(fmt.Sprintf("invalid request: %s", info.FullMethod))
	}
	h, err := handler(ctx, req) // make actual request
	return h, err
}
