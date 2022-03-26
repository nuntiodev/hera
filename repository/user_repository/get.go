package user_repository

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/softcorp-io/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongoRepository) Get(ctx context.Context, user *go_block.User) (*go_block.User, error) {
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
	protoUser := UserToProtoUser(&resp)
	if resp.Encrypted == true && r.encryptionKey != "" {
		if err := r.crypto.DecryptUser(r.encryptionKey, protoUser); err != nil {
			return nil, err
		}
	}
	return protoUser, nil
}
