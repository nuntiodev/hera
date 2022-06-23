package authenticator

import (
	"context"
	"github.com/nuntiodev/hera-proto/go_hera"
)

var (
	SystemAuthenticator Authenticator
)

type Authenticator interface {
	AuthenticateRequest(ctx context.Context, req *go_hera.HeraRequest) error
}

type NoAuthenticator struct{}

func (*NoAuthenticator) AuthenticateRequest(ctx context.Context, req *go_hera.HeraRequest) error {
	return nil
}

func New() Authenticator {
	if SystemAuthenticator == nil {
		return &NoAuthenticator{}
	}
	return SystemAuthenticator
}
