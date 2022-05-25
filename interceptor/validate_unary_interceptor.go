package interceptor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nuntiodev/block-proto/go_block"
	"google.golang.org/grpc"
)

var (
	TokenIsNil          = errors.New("token is nil")
	TokenPointerIsEmpty = errors.New("token pointer is nil")
	UpdateIsNil         = errors.New("update is nil")
	UserIsNil           = errors.New("user is nil")
	NamespaceIsEmpty    = errors.New("namespace is empty")
	UserBatchIsNil      = errors.New("user batch is nil")
	MeasurementIsNil    = errors.New("measurement is nil")
	ConfigIsNil         = errors.New("config is nil")
	AuthConfigIsNil     = errors.New("auth config is nil")
	TextIsNil           = errors.New("text is nil")
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
	case Heartbeat, GetAll, PublicKeys,
		NamespaceActiveHistory, GetConfig, DeleteConfig,
		InitializeApplication:
		break
	case BlockToken:
		if translatedReq.TokenPointer == "" {
			return nil, TokenPointerIsEmpty
		}
	case RecordActiveMeasurement, UserActiveHistory:
		if translatedReq.ActiveMeasurement == nil {
			return &go_block.UserResponse{}, MeasurementIsNil
		}
	case Get, Create, VerifyEmail,
		SendVerificationEmail, SendResetPasswordEmail, ResetPassword:
		if translatedReq.User == nil {
			return &go_block.UserResponse{}, UserIsNil
		}
	case UpdatePassword, UpdateMetadata,
		UpdateImage, UpdateEmail, UpdateUsername,
		UpdateName, UpdateBirthdate, UpdatePhoneNumber,
		UpdatePreferredLanguage:
		if translatedReq.Update == nil {
			return &go_block.UserResponse{}, UpdateIsNil
		} else if translatedReq.User == nil {
			return &go_block.UserResponse{}, UpdateIsNil
		}
	case ValidateCredentials, Delete, UpdateSecurity, Login:
		if translatedReq.User == nil {
			return nil, UserIsNil
		}
	case RefreshToken, GetTokens, BlockTokenById:
		if translatedReq.Token == nil {
			return nil, TokenIsNil
		}
	case ValidateToken:
		if translatedReq.TokenPointer == "" {
			return nil, errors.New("token pointer is nil")
		}
	case DeleteNamespace:
		if translatedReq.Namespace == "" {
			return nil, NamespaceIsEmpty
		}
	case DeleteBatch:
		if translatedReq.UserBatch == nil {
			return nil, UserBatchIsNil
		}
	case CreateNamespaceConfig, UpdateConfigSettings,
		UpdateConfigDetails:
		if translatedReq.Config == nil {
			return nil, ConfigIsNil
		}
	default:
		return &go_block.UserResponse{}, errors.New(fmt.Sprintf("invalid request: %s", info.FullMethod))
	}
	return handler(ctx, req) // make actual request
}
