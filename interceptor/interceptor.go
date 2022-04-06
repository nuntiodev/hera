package interceptor

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type DefaultInterceptor struct {
	zapLog *zap.Logger
}

type Interceptor interface {
	WithLogUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
	WithLogStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
	WithValidateUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
	WithValidateStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
}

func New(zapLog *zap.Logger) (Interceptor, error) {
	return &DefaultInterceptor{
		zapLog: zapLog,
	}, nil
}
