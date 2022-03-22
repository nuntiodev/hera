package interceptor

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

func (i *DefaultInterceptor) WithLogUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	i.zapLog.Debug(fmt.Sprintf("%v", req))
	h, err := handler(ctx, req) // make actual request
	if err != nil {
		i.zapLog.Error(fmt.Sprintf("Method:%s	Duration:%s   Error:%v",
			info.FullMethod,
			time.Since(start),
			err))
	} else {
		i.zapLog.Info(fmt.Sprintf("Method:%s	Duration:%s",
			info.FullMethod,
			time.Since(start),
		))
	}
	return h, err
}
