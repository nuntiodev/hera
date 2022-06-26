package user_repository

import (
	"context"
	"github.com/nuntiodev/hera-proto/go_hera"
	"github.com/nuntiodev/hera/models"
	"golang.org/x/crypto/bcrypt"
)

/*
	Create - this method creates a user and encrypts it if keys are present.
*/
func (r *mongodbRepository) Create(ctx context.Context, user *go_hera.User) (*models.User, error) {
	// validate data
	if user == nil {
		return nil, UserIsNilErr
	} else if err := validateEmail(user.GetEmail()); err != nil {
		return nil, err
	} else if err := validatePassword(user.Password); err != nil && r.validatePassword {
		return nil, err
	} else if err := validateMetadata(user.Metadata); err != nil {
		return nil, err
	} else if err := validatePhone(user.GetPhone()); err != nil {
		return nil, err
	}
	prepare(actionCreate, user)
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}
	create := models.ProtoUserToUser(user)
	// create hashes
	create.EmailHash, create.UsernameHash, create.PhoneHash = generateUserHashes(user)
	resp := *create
	if err := r.crypto.Encrypt(create); err != nil {
		return nil, err
	}
	if _, err := r.collection.InsertOne(ctx, create); err != nil {
		return nil, err
	}
	// set new data for user created
	return &resp, nil
}
