package http_server

import (
	"fmt"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/authenticator"
	"github.com/nuntiodev/hera/interceptor"
	"go.uber.org/zap"
	"net/http"
	"os"
)

var (
	port = ""
)

func initialize() error {
	var ok bool
	port, ok = os.LookupEnv("HTTP_PORT")
	if !ok || port == "" {
		port = "9001"
	}
	return nil
}

type Server struct {
	handler       go_hera.ServiceServer
	interceptor   interceptor.Interceptor
	authenticator authenticator.Authenticator
	logger        *zap.Logger
}

func New(handler go_hera.ServiceServer, interceptor interceptor.Interceptor, authenticator authenticator.Authenticator, logger *zap.Logger) (*Server, error) {
	return &Server{
		handler:       handler,
		interceptor:   interceptor,
		authenticator: authenticator,
		logger:        logger,
	}, nil
}

func (s *Server) Run() error {
	s.logger.Info(fmt.Sprintf("starting Hera http server on port %s", port))
	return http.ListenAndServe(fmt.Sprintf(":%s", port), s.routes())
}
