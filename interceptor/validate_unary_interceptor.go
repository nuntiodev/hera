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
	if err := ValidateRequest(method[1], translatedReq); err != nil {
		return &go_hera.HeraResponse{}, err
	}
	return handler(ctx, req) // make actual request
}

func ValidateRequest(name string, request *go_hera.HeraRequest) error {
	switch name {
	case Heartbeat, PublicKeys, DeleteNamespace,
		RegisterRsaKey, RemovePublicKey, GetConfig,
		/*todo: implment*/ ResetPassword:
		break
	case CreateUser, GetUser, ValidateCredentials,
		Login, DeleteUser, SendVerificationEmail,
		VerifyEmail, SendVerificationText, VerifyPhone,
		SendResetPasswordEmail, SendResetPasswordText:
		if request.User == nil {
			return UserIsNil
		}
	case CreateNamespace, UpdateConfig, DeleteConfig:
		if request.Config == nil {
			return ConfigIsNil
		}
	case DeleteUsers:
		if request.Users == nil {
			return UsersIsNil
		}
	case UpdateUserMetadata, UpdateUserProfile, UpdateUserContact,
		UpdateUserPassword:
		if request.User == nil {
			return UserIsNil
		} else if request.UserUpdate == nil {
			return UpdateIsNil
		}
	case SearchForUser, ListUsers:
		if request.Query == nil {
			return QueryIsNil
		}
	case CreateTokenPair:
		if request.User == nil {
			return UserIsNil
		} else if request.Token == nil {
			return TokenIsNil
		}
	case ValidateToken:
		if request.TokenPointer == "" {
			return TokenIsNil
		}
	case BlockToken:
		if request.TokenPointer == "" && request.Token == nil {
			return TokenIsNil
		}
	case RefreshToken, GetTokens:
		if request.Token == nil {
			return TokenIsNil
		}

	default:
		return errors.New(fmt.Sprintf("invalid request: %s", name))
	}
	return nil
}
