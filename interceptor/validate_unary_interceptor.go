package interceptor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nuntiodev/hera-proto/go_hera"
	"google.golang.org/grpc"
)

var (
	TokenIsNil          = errors.New("token is nil")
	TokenPointerIsEmpty = errors.New("token pointer is nil")
	UpdateIsNil         = errors.New("update is nil")
	UserIsNil           = errors.New("user is nil")
	UsersIsNil          = errors.New("users is nil")
	QueryIsNil          = errors.New("query is nil")
	NamespaceIsEmpty    = errors.New("namespace is empty")
	UserBatchIsNil      = errors.New("user batch is nil")
	ConfigIsNil         = errors.New("config is nil")
)

func (i *DefaultInterceptor) WithValidateUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info == nil {
		return nil, errors.New("invalid request")
	}
	translatedReq, ok := req.(*go_hera.HeraRequest)
	if !ok {
		translatedReq = &go_hera.HeraRequest{}
	}
	method := strings.Split(info.FullMethod, ProjectName)
	if len(method) != 2 {
		return nil, errors.New(fmt.Sprintf("invalid method call: %s", info.FullMethod))
	}
	switch method[1] {
	case Heartbeat, PublicKeys, DeleteNamespace,
		RegisterRsaKey, RemovePublicKey, GetConfig,
		/*todo: implment*/ ResetPassword:
		break
	case CreateUser, GetUser, ValidateCredentials,
		Login, DeleteUser, SendVerificationEmail,
		VerifyEmail, SendVerificationText, VerifyPhone,
		SendResetPasswordEmail, SendResetPasswordText:
		if translatedReq.User == nil {
			return nil, UserIsNil
		}
	case CreateNamespace, UpdateConfig, DeleteConfig:
		if translatedReq.Config == nil {
			return nil, ConfigIsNil
		}
	case DeleteUsers:
		if translatedReq.Users == nil {
			return nil, UsersIsNil
		}
	case UpdateUserMetadata, UpdateUserProfile, UpdateUserContact,
		UpdateUserPassword:
		if translatedReq.User == nil {
			return nil, UserIsNil
		} else if translatedReq.UserUpdate == nil {
			return nil, UpdateIsNil
		}
	case SearchForUser, ListUsers:
		if translatedReq.Query == nil {
			return nil, QueryIsNil
		}
	case CreateTokenPair:
		if translatedReq.User == nil {
			return nil, UserIsNil
		} else if translatedReq.Token == nil {
			return nil, TokenIsNil
		}
	case ValidateToken:
		if translatedReq.TokenPointer == "" {
			return nil, TokenIsNil
		}
	case BlockToken:
		if translatedReq.TokenPointer == "" && translatedReq.Token == nil {
			return nil, TokenIsNil
		}
	case RefreshToken, GetTokens:
		if translatedReq.Token == nil {
			return nil, TokenIsNil
		}

	default:
		return nil, errors.New(fmt.Sprintf("invalid request: %s", info.FullMethod))
	}
	return handler(ctx, req) // make actual request
}
