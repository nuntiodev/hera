package user_repository

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/nuntiodev/nuntio-user-block/models"

	"github.com/nuntiodev/block-proto/go_block"
	"golang.org/x/crypto/bcrypt"
)

/*
	Create - this method creates a user and encrypts it if keys are present.
*/
func (r *mongodbRepository) Create(ctx context.Context, user *go_block.User) (*models.User, error) {
	prepare(actionCreate, user)
	if err := r.validate(actionCreate, user); err != nil {
		return nil, err
	}
	emailHash := ""
	if user.Email != "" {
		emailHash = fmt.Sprintf("%x", md5.Sum([]byte(user.Email)))
	}
	phoneNumberHash := ""
	if user.PhoneNumber != "" {
		phoneNumberHash = fmt.Sprintf("%x", md5.Sum([]byte(user.PhoneNumber)))
	}
	usernameHash := ""
	if user.Username != "" {
		usernameHash = fmt.Sprintf("%x", md5.Sum([]byte(user.Username)))
	}
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}
	create := models.ProtoUserToUser(user)
	create.EmailHash = emailHash
	create.PhoneNumberHash = phoneNumberHash
	create.UsernameHash = usernameHash
	copy := *create
	if err := r.crypto.Encrypt(create); err != nil {
		return nil, err
	}
	if _, err := r.collection.InsertOne(ctx, create); err != nil {
		return nil, err
	}
	// set new data for user created
	return &copy, nil
}
