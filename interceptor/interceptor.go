package interceptor

import (
	"context"
	"github.com/nuntiodev/hera/authenticator"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	ProjectName            = "/Hera.Service/"
	Heartbeat              = "Heartbeat"
	CreateUser             = "CreateUser"
	UpdateUserMetadata     = "UpdateUserMetadata"
	UpdateUserProfile      = "UpdateUserProfile"
	UpdateUserContact      = "UpdateUserContact"
	UpdateUserPassword     = "UpdateUserPassword"
	SearchForUser          = "SearchForUser"
	GetUser                = "GetUser"
	ListUsers              = "ListUsers"
	ValidateCredentials    = "ValidateCredentials"
	Login                  = "Login"
	DeleteUser             = "DeleteUser"
	DeleteUsers            = "DeleteUsers"
	CreateTokenPair        = "CreateTokenPair"
	ValidateToken          = "ValidateToken"
	BlockToken             = "BlockToken"
	RefreshToken           = "RefreshToken"
	GetTokens              = "GetTokens"
	PublicKeys             = "PublicKeys"
	SendVerificationEmail  = "SendVerificationEmail"
	VerifyEmail            = "VerifyEmail"
	SendVerificationText   = "SendVerificationText"
	VerifyPhone            = "VerifyPhone"
	SendResetPasswordEmail = "SendResetPasswordEmail"
	SendResetPasswordText  = "SendResetPasswordText"
	ResetPassword          = "ResetPassword"
	DeleteNamespace        = "DeleteNamespace"
	CreateNamespace        = "CreateNamespace"
	RegisterRsaKey         = "RegisterRsaKey"
	RemovePublicKey        = "RemovePublicKey"
	GetConfig              = "GetConfig"
	UpdateConfig           = "UpdateConfig"
	DeleteConfig           = "DeleteConfig"
)

type DefaultInterceptor struct {
	zapLog        *zap.Logger
	authenticator authenticator.Authenticator
}

type Interceptor interface {
	WithAuthenticateUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
	WithLogUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
	WithLogStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
	WithValidateUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
	WithValidateStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
}

func New(zapLog *zap.Logger, authenticator authenticator.Authenticator) (Interceptor, error) {
	return &DefaultInterceptor{
		zapLog:        zapLog,
		authenticator: authenticator,
	}, nil
}
