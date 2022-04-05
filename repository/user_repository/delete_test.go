package user_repository

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/x/cryptox"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteIEEncrypted(t *testing.T) {
	// setup available clients
	var clients []*mongoRepository
	userRepositoryFullEncryption, err := getTestUserRepository(context.Background(), true, true, "")
	assert.NoError(t, err)
	userRepositoryInternalEncryption, err := getTestUserRepository(context.Background(), true, false, "")
	assert.NoError(t, err)
	userRepositoryExternalEncryption, err := getTestUserRepository(context.Background(), false, true, "")
	assert.NoError(t, err)
	userRepositoryNoEncryption, err := getTestUserRepository(context.Background(), false, false, "")
	assert.NoError(t, err)
	clients = []*mongoRepository{userRepositoryFullEncryption, userRepositoryInternalEncryption, userRepositoryExternalEncryption, userRepositoryNoEncryption}
	// delete all users from other tests (we use the same collection)
	err = userRepositoryExternalEncryption.DeleteAll(context.Background())
	assert.NoError(t, err)
	for index, userRepository := range clients {
		assert.NoError(t, err)
		password := gofakeit.Password(true, true, true, true, true, 30)
		createdUserOne, err := userRepository.Create(context.Background(), &go_block.User{
			Password: password,
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdUserOne)
		createdUser, err := userRepository.Create(context.Background(), &go_block.User{
			Password: password,
		})
		// set new encryption key
		encryptionKey, err := userRepository.crypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		userRepository.internalEncryptionKeys = append(userRepository.internalEncryptionKeys, encryptionKey)
		// act
		err = userRepository.Delete(context.Background(), createdUser)
		// validate
		assert.NoError(t, err)
		// validate repository is not empty
		getUser, err := userRepository.Get(context.Background(), createdUser, false)
		assert.Error(t, err, index)
		assert.Nil(t, getUser)
	}
}
