package interceptor

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/authenticator"
	"google.golang.org/grpc"
)

func (i *DefaultInterceptor) WithAuthenticateUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info == nil {
		return nil, errors.New("invalid request")
	}
	translatedReq, ok := req.(*go_hera.HeraRequest)
	if !ok {
		translatedReq = &go_hera.HeraRequest{}
	}
	if err := i.authenticator.AuthenticateRequest(ctx, translatedReq, &authenticator.Info{IsGrpc: true}); err != nil {
		return nil, err
	}
	return handler(ctx, translatedReq)
}
