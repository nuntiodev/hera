package user_repository

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/softcorp-io/block-proto/go_block"
	"golang.org/x/crypto/bcrypt"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (r *mongoRepository) Create(ctx context.Context, user *go_block.User) (*go_block.User, error) {
	prepare(actionCreate, user)
	if err := r.validate(actionCreate, user); err != nil {
		return nil, err
	}
	emailHash := ""
	if user.Email != "" {
		emailHash = fmt.Sprintf("%x", md5.Sum([]byte(user.Email)))
	}
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}
	resp := *user
	if r.encryptionKey != "" {
		if err := r.crypto.EncryptUser(r.encryptionKey, user); err != nil {
			return nil, err
		}
		user.Encrypted = true
		user.EncryptedAt = ts.Now()
		resp.Encrypted = true
	} else {
		user.Encrypted = false
	}
	createUser := ProtoUserToUser(user)
	createUser.EmailHash = emailHash
	if _, err := r.collection.InsertOne(ctx, createUser); err != nil {
		return nil, err
	}
	return &resp, nil
}
