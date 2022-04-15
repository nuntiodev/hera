package user_repository

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/io-nuntio/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) Get(ctx context.Context, user *go_block.User, upgrade bool) (*go_block.User, error) {
	prepare(actionGet, user)
	if err := r.validate(actionGet, user); err != nil {
		return nil, err
	}
	filter := bson.M{}
	if user.Id != "" {
		filter = bson.M{"_id": user.Id}
	} else if user.Email != "" {
		filter = bson.M{"email_hash": fmt.Sprintf("%x", md5.Sum([]byte(user.Email)))}
	} else if user.OptionalId != "" {
		filter = bson.M{"optional_id": user.OptionalId}
	}
	resp := User{}
	if err := r.collection.FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, err
	}
	if resp.InternalEncrypted || resp.ExternalEncrypted {
		if err := r.decryptUser(ctx, &resp, upgrade); err != nil {
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
