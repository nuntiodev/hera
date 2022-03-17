package user_repository

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/softcorp-io/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongoRepository) UpdateEmail(ctx context.Context, get *go_block.User, update *go_block.User, encryptionKey string) (*go_block.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateEmail, update)
	if err := r.validate(actionUpdateEmail, update); err != nil {
		return nil, err
	}
	emailHash := ""
	if update.Email != "" {
		emailHash = fmt.Sprintf("%x", md5.Sum([]byte(update.Email)))
	}
	getUser, err := r.Get(ctx, get, encryptionKey) // check if user encryption is turned on
	if err != nil {
		return nil, err
	}
	resp := *update
	if err := r.handleEncryption(getUser.Encrypted, update, encryptionKey); err != nil {
		return nil, err
	}
	updateUser := ProtoUserToUser(&go_block.User{
		Email:     update.Email,
		UpdatedAt: update.UpdatedAt,
	})
	updateUser.EmailHash = emailHash
	mongoUpdate := bson.M{
		"$set": bson.M{
			"email":        updateUser.Email,
			"email_hash":   updateUser.EmailHash,
			"updated_at":   updateUser.UpdatedAt,
			"encrypted_at": updateUser.EncryptedAt,
		},
	}
	updateResult, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": getUser.Id},
		mongoUpdate,
	)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("could not find get")
	}
	return &resp, nil
}
