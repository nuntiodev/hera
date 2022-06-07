package user_repository

import (
	"context"
	"fmt"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
)

/*
	Get - this method fetches a user either by id, username, or email.
*/
func (r *mongodbRepository) Get(ctx context.Context, user *go_block.User) (*models.User, error) {
	prepare(actionGet, user)
	if err := r.validate(actionGet, user); err != nil {
		return nil, err
	}
	filter, err := getUserFilter(user)
	if err != nil {
		return nil, err
	}
	resp := models.User{}
	if err := r.collection.FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, fmt.Errorf("could not find user with id: %v and err: %v", filter, err)
	}
	if err := r.crypto.Decrypt(&resp); err != nil {
		return nil, err
	}
	// check if we should upgrade the encryption level
	if upgradable, _ := r.crypto.Upgradeble(&resp); upgradable {
		if err := r.upgradeEncryptionLevel(ctx, &resp); err != nil {
			return nil, err
		}
	}
	return &resp, nil
}
