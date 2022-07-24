package user_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/models"
)

/*
	Create - this method creates a user and encrypts it if keys are present.
*/
func (r *mongodbRepository) Create(ctx context.Context, user *go_hera.User) (*go_hera.User, error) {
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
	if user.Password != nil && user.Password.Body != "" {
		if r.hasher == nil {
			return nil, errors.New("hasher is nil")
		}
		hashedPassword, err := r.hasher.Generate(user.Password.Body)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}
	create := models.ProtoUserToUser(user)
	// create hashes
	create.EmailHash, create.UsernameHash, create.PhoneHash = generateUserHashes(user)
	if err := r.crypto.Encrypt(create); err != nil {
		return nil, err
	}
	if _, err := r.collection.InsertOne(ctx, create); err != nil {
		return nil, err
	}
	// set new data for user created
	return user, nil
}
