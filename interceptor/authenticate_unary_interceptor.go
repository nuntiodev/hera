package interceptor

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"google.golang.org/grpc"
)

func (i *DefaultInterceptor) WithAuthenticateUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info == nil {
		return nil, errors.New("invalid request")
	}
	translatedReq, ok := req.(*go_block.UserRequest)
	if !ok {
		translatedReq = &go_block.UserRequest{}
	}
	if err := i.authenticator.AuthenticateRequest(ctx, translatedReq); err != nil {
		return &go_block.UserResponse{}, err
	}
	return handler(ctx, translatedReq)
}
