package interceptor

import (
	"context"
	"errors"
	"fmt"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"google.golang.org/grpc"
	"strings"
)

const (
	ProjectName         = "/BlockUser.Service/"
	Heartbeat           = "Heartbeat"
	Create              = "Create"
	UpdatePassword      = "UpdatePassword"
	UpdateProfile       = "UpdateProfile"
	UpdateSecurity      = "UpdateSecurity"
	Get                 = "Get"
	GetAll              = "GetAll"
	Search              = "Search"
	ValidateCredentials = "ValidateCredentials"
	Delete              = "Delete"
	DeleteNamespace     = "DeleteNamespace"
)

var (
	ErrorReqIsNil    = errors.New("req is nil")
	UpdateIsNil      = errors.New("update is nil")
	UserIsNil        = errors.New("user is nil")
	NamespaceIsEmpty = errors.New("namespace is empty")
)

func (i *DefaultInterceptor) WithValidateUnaryInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	if info == nil {
		return nil, errors.New("invalid request")
	}
	translatedReq, ok := req.(*block_user.UserRequest)
	if !ok {
		translatedReq = &block_user.UserRequest{}
	}
	method := strings.Split(info.FullMethod, ProjectName)
	if len(method) != 2 {
		return nil, errors.New(fmt.Sprintf("invalid method call: %s", info.FullMethod))
	}
	switch method[1] {
	case Heartbeat:
		break
	case Create, Get:
		if translatedReq == nil {
			return nil, ErrorReqIsNil
		} else if translatedReq.User == nil {
			return &block_user.UserResponse{}, UpdateIsNil
		}
	case UpdatePassword, UpdateSecurity, UpdateProfile:
		if translatedReq == nil {
			return nil, ErrorReqIsNil
		} else if translatedReq.Update == nil {
			return &block_user.UserResponse{}, UpdateIsNil
		} else if translatedReq.User == nil {
			return &block_user.UserResponse{}, UpdateIsNil
		}
	case GetAll, Search:
		if translatedReq == nil {
			return nil, ErrorReqIsNil
		}
	case ValidateCredentials, Delete:
		if translatedReq == nil {
			return nil, ErrorReqIsNil
		} else if translatedReq.User == nil {
			return nil, UserIsNil
		}
	case DeleteNamespace:
		if translatedReq == nil {
			return nil, ErrorReqIsNil
		} else if translatedReq.Namespace == "" {
			return nil, NamespaceIsEmpty
		}
	default:
		return &block_user.UserResponse{}, errors.New(fmt.Sprintf("invalid request: %s", info.FullMethod))
	}
	h, err := handler(ctx, req) // make actual request
	return h, err
}
