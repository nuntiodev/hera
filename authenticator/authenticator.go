package authenticator

import (
	"context"
	"github.com/nuntiodev/hera-proto/go_hera"
)

var (
	SystemAuthenticator Authenticator
)

type Info struct {
	IsGrpc bool
	IsHttp bool
}

type Authenticator interface {
	AuthenticateRequest(ctx context.Context, req *go_hera.HeraRequest, info *Info) error
}

type NoAuthenticator struct{}

func (*NoAuthenticator) AuthenticateRequest(ctx context.Context, req *go_hera.HeraRequest, info *Info) error {
	return nil
}

func New() Authenticator {
	if SystemAuthenticator == nil {
		return &NoAuthenticator{}
	}
	return SystemAuthenticator
}
