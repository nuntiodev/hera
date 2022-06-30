package user_repository

import (
	"context"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/models"
	"golang.org/x/crypto/bcrypt"
)

/*
	Create - this method creates a user and encrypts it if keys are present.
*/
func (r *mongodbRepository) Create(ctx context.Context, user *go_hera.User) error {
	// validate data
	if user == nil {
		return UserIsNilErr
	} else if err := validateEmail(user.GetEmail()); err != nil {
		return err
	} else if err := validatePassword(user.Password); err != nil && r.validatePassword {
		return err
	} else if err := validateMetadata(user.Metadata); err != nil {
		return err
	} else if err := validatePhone(user.GetPhone()); err != nil {
		return err
	}
	prepare(actionCreate, user)
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}
	create := models.ProtoUserToUser(user)
	// create hashes
	create.EmailHash, create.UsernameHash, create.PhoneHash = generateUserHashes(user)
	if err := r.crypto.Encrypt(create); err != nil {
		return err
	}
	if _, err := r.collection.InsertOne(ctx, create); err != nil {
		return err
	}
	// set new data for user created
	return nil
}
