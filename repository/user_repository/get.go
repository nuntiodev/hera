package user_repository

import (
	"context"
	"fmt"
	"github.com/nuntiodev/block-proto/go_block"
)

func (r *mongodbRepository) Get(ctx context.Context, user *go_block.User, upgrade bool) (*go_block.User, error) {
	prepare(actionGet, user)
	if err := r.validate(actionGet, user); err != nil {
		return nil, err
	}
	filter, err := getUserFilter(user)
	if err != nil {
		return nil, err
	}
	resp := User{}
	if err := r.collection.FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, fmt.Errorf("could not find user with id: %v and err: %v", filter, err)
	}
	if resp.InternalEncryptionLevel > 0 || resp.ExternalEncryptionLevel > 0 {
		if err := r.decryptUser(&resp); err != nil {
			return nil, err
		}
	}
	// check if we should upgrade the encryption level
	if upgrade && r.isEncryptionLevelUpgradable(&resp) {
		if err := r.upgradeInternalEncryptionLevel(ctx, &resp); err != nil {
			return nil, err
		}
	}
	return UserToProtoUser(&resp), nil
}
