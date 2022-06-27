package interceptor

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

func (i *DefaultInterceptor) WithLogUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req) // make actual request
	if err != nil {
		i.logger.Error(fmt.Sprintf("Hera: Method:%s	Duration:%s   Error:%v",
			info.FullMethod,
			time.Since(start),
			err))
	} else {
		i.logger.Debug(fmt.Sprintf("Method:%s	Duration:%s",
			info.FullMethod,
			time.Since(start),
		))
	}
	return h, err
}
