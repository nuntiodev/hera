package authenticator

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
)

var (
	SystemAuthenticator Authenticator
)

type Authenticator interface {
	AuthenticateRequest(ctx context.Context, req *go_block.UserRequest) error
}

type NoAuthenticator struct{}

func (*NoAuthenticator) AuthenticateRequest(ctx context.Context, req *go_block.UserRequest) error {
	return nil
}

func New() Authenticator {
	if SystemAuthenticator == nil {
		return &NoAuthenticator{}
	}
	return SystemAuthenticator
}
