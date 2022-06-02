package interceptor

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/authenticator"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	ProjectName             = "/BlockUser.UserService/"
	Heartbeat               = "Heartbeat"
	Create                  = "Create"
	UpdatePassword          = "UpdatePassword"
	UpdateMetadata          = "UpdateMetadata"
	UpdateImage             = "UpdateImage"
	UpdateEmail             = "UpdateEmail"
	UpdatePhoneNumber       = "UpdatePhoneNumber"
	UpdateName              = "UpdateName"
	UpdateBirthdate         = "UpdateBirthdate"
	UpdateUsername          = "UpdateUsername"
	UpdatePreferredLanguage = "UpdatePreferredLanguage"
	UpdateSecurity          = "UpdateSecurity"
	Get                     = "Get"
	GetAll                  = "GetAll"
	ValidateCredentials     = "ValidateCredentials"
	Login                   = "Login"
	ValidateToken           = "ValidateToken"
	BlockToken              = "BlockToken"
	BlockTokenById          = "BlockTokenById"
	RefreshToken            = "RefreshToken"
	GetTokens               = "GetTokens"
	PublicKeys              = "PublicKeys"
	RecordActiveMeasurement = "RecordActiveMeasurement"
	UserActiveHistory       = "UserActiveHistory"
	NamespaceActiveHistory  = "NamespaceActiveHistory"
	SendVerificationEmail   = "SendVerificationEmail"
	VerifyEmail             = "VerifyEmail"
	SendResetPasswordEmail  = "SendResetPasswordEmail"
	ResetPassword           = "ResetPassword"
	Delete                  = "Delete"
	DeleteBatch             = "DeleteBatch"
	DeleteNamespace         = "DeleteNamespace"
	CreateNamespaceConfig   = "CreateNamespaceConfig"
	UpdateConfigSettings    = "UpdateConfigSettings"
	UpdateConfigDetails     = "UpdateConfigDetails"
	GetConfig               = "GetConfig"
	DeleteConfig            = "DeleteConfig"
	InitializeApplication   = "InitializeApplication"
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
