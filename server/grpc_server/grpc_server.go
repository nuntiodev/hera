package grpc_server

import (
	"errors"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"github.com/softcorp-io/block-user-service/handler"
	"github.com/softcorp-io/block-user-service/interceptor"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

var (
	port = ""
)

type Server struct {
	zapLog      *zap.Logger
	handler     handler.Handler
	interceptor interceptor.Interceptor
}

func initialize() error {
	var ok bool
	port, ok = os.LookupEnv("GRPC_PORT")
	if !ok || port == "" {
		return errors.New("missing required GRPC_PORT")
	}
	return nil
}

func New(zapLog *zap.Logger, handler handler.Handler, interceptor interceptor.Interceptor) (*Server, error) {
	if err := initialize(); err != nil {
		return nil, err
	}
	return &Server{
		zapLog:      zapLog,
		handler:     handler,
		interceptor: interceptor,
	}, nil
}

func (s *Server) Run() error {
	s.zapLog.Info("starting gRPC server")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}
	defer lis.Close()
	s.zapLog.Info(fmt.Sprintf("gRPC server running on port: %s", port))
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				s.interceptor.WithLogUnaryInterceptor,
				s.interceptor.WithValidateUnaryInterceptor,
			),
		),
	)
	reflection.Register(grpcServer)
	block_user.RegisterServiceServer(grpcServer, s.handler)
	return grpcServer.Serve(lis)
}
