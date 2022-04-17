package user_repository

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/x/cryptox"
	"github.com/stretchr/testify/assert"
)

func TestDeleteBatchIEEncrypted(t *testing.T) {
	// setup available clients
	var clients []*mongodbRepository
	userRepositoryFullEncryption, err := getTestUserRepository(context.Background(), true, true, "")
	assert.NoError(t, err)
	userRepositoryInternalEncryption, err := getTestUserRepository(context.Background(), true, false, "")
	assert.NoError(t, err)
	userRepositoryExternalEncryption, err := getTestUserRepository(context.Background(), false, true, "")
	assert.NoError(t, err)
	userRepositoryNoEncryption, err := getTestUserRepository(context.Background(), false, false, "")
	assert.NoError(t, err)
	clients = []*mongodbRepository{userRepositoryFullEncryption, userRepositoryInternalEncryption, userRepositoryExternalEncryption, userRepositoryNoEncryption}
	// delete all users from other tests (we use the same collection)
	err = userRepositoryExternalEncryption.DeleteAll(context.Background())
	assert.NoError(t, err)
	for _, userRepository := range clients {
		assert.NoError(t, err)
		password := gofakeit.Password(true, true, true, true, true, 30)
		createdUserOne, err := userRepository.Create(context.Background(), &go_block.User{
			Password: password,
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdUserOne)
		createdUserTwo, err := userRepository.Create(context.Background(), &go_block.User{
			Password: password,
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdUserTwo)
		createdUserThree, err := userRepository.Create(context.Background(), &go_block.User{
			Password: password,
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdUserThree)
		// set new encryption key
		encryptionKey, err := userRepository.crypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		userRepository.internalEncryptionKeys = append(userRepository.internalEncryptionKeys, encryptionKey)
		// act
		err = userRepository.DeleteBatch(context.Background(), []*go_block.User{createdUserOne, createdUserTwo})
		// validate
		assert.NoError(t, err)
		// validate
		getUser, err := userRepository.Get(context.Background(), createdUserOne, false)
		assert.Nil(t, getUser)
		assert.Error(t, err)
		getUser, err = userRepository.Get(context.Background(), createdUserTwo, false)
		assert.Nil(t, getUser)
		assert.Error(t, err)
		getUser, err = userRepository.Get(context.Background(), createdUserThree, false)
		assert.NotNil(t, getUser)
		assert.NoError(t, err)
	}
}
