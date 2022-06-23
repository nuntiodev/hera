package handler

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-proto/go_hera"
)

/*
	DeleteNamespace - this method deletes an entire namespace. This includes deleting all users and the namespace config.
*/
func (h *defaultHandler) DeleteNamespace(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	if req.Namespace == "" {
		return nil, errors.New("cannot drop default database")
	}
	return nil, h.repository.DropDatabase(ctx, req.Namespace)
}
