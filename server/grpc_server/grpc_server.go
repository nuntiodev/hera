package grpc_server

import (
	"fmt"
	"net"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/handler"
	"github.com/nuntiodev/hera/interceptor"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
		port = "9000"
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
				s.interceptor.WithAuthenticateUnaryInterceptor,
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				s.interceptor.WithLogStreamInterceptor,
				s.interceptor.WithValidateStreamInterceptor,
			),
		),
	)
	reflection.Register(grpcServer)
	go_hera.RegisterServiceServer(grpcServer, s.handler)
	return grpcServer.Serve(lis)
}
