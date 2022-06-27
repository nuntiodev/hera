package interceptor

import (
	"fmt"
	"google.golang.org/grpc"
)

func (i *DefaultInterceptor) WithLogStreamInterceptor(req interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	i.logger.Info(fmt.Sprintf("New streaming request:%s",
		info.FullMethod,
	))
	// make actual request
	if err := handler(req, ss); err != nil {
		i.logger.Error(err.Error())
		return err
	}
	return nil
}
