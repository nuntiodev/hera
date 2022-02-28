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
	WithValidateUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

func New(zapLog *zap.Logger) (Interceptor, error) {
	return &DefaultInterceptor{
		zapLog: zapLog,
	}, nil
}
