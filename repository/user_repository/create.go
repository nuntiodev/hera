package user_repository

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/nuntiodev/block-proto/go_block"
	"golang.org/x/crypto/bcrypt"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (r *mongodbRepository) Create(ctx context.Context, user *go_block.User) (*go_block.User, error) {
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
	create := ProtoUserToUser(user)
	create.EmailHash = emailHash
	if err := r.encryptUser(ctx, actionCreate, create); err != nil {
		return nil, err
	}
	if _, err := r.collection.InsertOne(ctx, create); err != nil {
		return nil, err
	}
	// set new data for user created
	user.EncryptedAt = ts.New(create.EncryptedAt)
	user.ExternalEncrypted = create.ExternalEncrypted
	user.InternalEncrypted = create.InternalEncrypted
	user.ExternalEncryptionLevel = int32(create.ExternalEncryptionLevel)
	user.InternalEncryptionLevel = int32(create.InternalEncryptionLevel)
	return user, nil
}
