package interceptor

import (
	"errors"
	"google.golang.org/grpc"
)

func (i *DefaultInterceptor) WithValidateStreamInterceptor(req interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if info == nil {
		return errors.New("invalid request")
	}
	return handler(req, ss) // make actual request
}
