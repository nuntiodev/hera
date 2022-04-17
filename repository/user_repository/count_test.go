package user_repository

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/x/cryptox"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCountIEEncrypted(t *testing.T) {
	// setup available clients
	var clients []*mongodbRepository
	ns := uuid.NewV4().String()
	userRepositoryFullEncryption, err := getTestUserRepository(context.Background(), true, true, ns)
	assert.NoError(t, err)
	userRepositoryInternalEncryption, err := getTestUserRepository(context.Background(), true, false, ns)
	assert.NoError(t, err)
	userRepositoryExternalEncryption, err := getTestUserRepository(context.Background(), false, true, ns)
	assert.NoError(t, err)
	userRepositoryNoEncryption, err := getTestUserRepository(context.Background(), false, false, ns)
	assert.NoError(t, err)
	clients = []*mongodbRepository{userRepositoryFullEncryption, userRepositoryInternalEncryption, userRepositoryExternalEncryption, userRepositoryNoEncryption}
	for index, userRepository := range clients {
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
		count, err := userRepository.Count(context.Background())
		// validate
		assert.NoError(t, err)
		assert.Equal(t, 3, int(count), index)
		// delete all at the end
		assert.NoError(t, userRepository.DeleteBatch(context.Background(), []*go_block.User{
			createdUserOne,
			createdUserTwo,
			createdUserThree,
		}))
	}
}
