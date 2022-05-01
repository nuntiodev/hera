package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nuntiodev/block-proto/go_block"
)

func (h *defaultHandler) CreateNamespaceConfig(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	// create test user in namespace
	users, err := h.repository.Users().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not build user with err: %v", err)
	}
	metadata, err := json.Marshal(map[string]string{
		"role": "test",
	})
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not marshal user with err: %v", err)
	}
	users.Create(ctx, &go_block.User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@user.io",
		Metadata:  string(metadata),
	})
	config, err := h.repository.Config(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not build config with err: %v", err)
	}
	resp, err := config.Create(ctx, req.Config)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not create config with err: %v", err)
	}
	return &go_block.UserResponse{
		Config: resp,
	}, nil
}
