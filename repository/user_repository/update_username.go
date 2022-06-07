package user_repository

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) UpdateUsername(ctx context.Context, get *go_block.User, update *go_block.User) (*models.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateUsername, update)
	if err := r.validate(actionUpdateUsername, update); err != nil {
		return nil, err
	}
	usernameHash := fmt.Sprintf("%x", md5.Sum([]byte(update.Username)))
	updateUser := models.ProtoUserToUser(&go_block.User{
		Username:     update.Username,
		UsernameHash: usernameHash,
		UpdatedAt:    update.UpdatedAt,
	})
	if err := r.crypto.Encrypt(updateUser); err != nil {
		return nil, err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"username":      updateUser.Username,
			"username_hash": updateUser.UsernameHash,
			"updated_at":    updateUser.UpdatedAt,
		},
	}
	filter, err := getUserFilter(get)
	if err != nil {
		return nil, err
	}
	result := r.collection.FindOneAndUpdate(
		ctx,
		filter,
		mongoUpdate,
	)
	if err := result.Err(); err != nil {
		return nil, err
	}
	var resp models.User
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	if err := r.crypto.Decrypt(&resp); err != nil {
		return nil, err
	}
	// set updated fields
	resp.Username = cryptox.Stringx{
		Body:                    update.Username,
		InternalEncryptionLevel: resp.Username.InternalEncryptionLevel,
		ExternalEncryptionLevel: resp.Username.ExternalEncryptionLevel,
	}
	resp.UsernameHash = usernameHash
	return &resp, nil
}
