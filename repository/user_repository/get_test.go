package user_repository

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/io-nuntio/block-proto/go_block"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/x/cryptox"
	"github.com/stretchr/testify/assert"
)

func TestGetByIdIEEncrypted(t *testing.T) {
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
	for index, userRepository := range clients {
		internalLevel := len(userRepository.internalEncryptionKeys)
		// create some metadata
		metadata, err := json.Marshal(&CustomMetadata{
			Name:      gofakeit.Name(),
			ClassYear: 3,
		})
		assert.NoError(t, err)
		password := gofakeit.Password(true, true, true, true, true, 30)
		user := &go_block.User{
			OptionalId: uuid.NewV4().String(),
			Email:      gofakeit.Email(),
			Password:   password,
			Image:      gofakeit.ImageURL(10, 10),
			Metadata:   string(metadata),
		}
		createdUser, err := userRepository.Create(context.Background(), user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		// set new encryption key
		encryptionKey, err := userRepository.crypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		userRepository.internalEncryptionKeys = append(userRepository.internalEncryptionKeys, encryptionKey)
		// act
		getUser, err := userRepository.Get(context.Background(), &go_block.User{
			Id: createdUser.Id,
		}, true)
		// validate
		assert.NoError(t, err)
		assert.NotNil(t, getUser)
		assert.NoError(t, compareUsers(createdUser, getUser, false))
		// validate we are at level 3
		getUser, err = userRepository.Get(context.Background(), &go_block.User{
			Id: createdUser.Id,
		}, true)
		assert.NoError(t, err)
		assert.Equal(t, int32(internalLevel+1), getUser.InternalEncryptionLevel, index)
	}
}

func TestGetByEmailIEEncrypted(t *testing.T) {
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
	for _, userRepository := range clients {
		// create some metadata
		metadata, err := json.Marshal(&CustomMetadata{
			Name:      gofakeit.Name(),
			ClassYear: 3,
		})
		assert.NoError(t, err)
		password := gofakeit.Password(true, true, true, true, true, 30)
		user := &go_block.User{
			OptionalId: uuid.NewV4().String(),
			Email:      gofakeit.Email(),
			Password:   password,
			Image:      gofakeit.ImageURL(10, 10),
			Metadata:   string(metadata),
		}
		createdUser, err := userRepository.Create(context.Background(), user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		// set new encryption key
		encryptionKey, err := userRepository.crypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		userRepository.internalEncryptionKeys = append(userRepository.internalEncryptionKeys, encryptionKey)
		// act
		getUser, err := userRepository.Get(context.Background(), &go_block.User{
			Email: createdUser.Email,
		}, true)
		assert.NoError(t, err)
		assert.NotNil(t, getUser)
		assert.NoError(t, compareUsers(createdUser, getUser, false))
	}
}

func TestGetByOptionalIdIEEncrypted(t *testing.T) {
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
	for _, userRepository := range clients {
		// create some metadata
		metadata, err := json.Marshal(&CustomMetadata{
			Name:      gofakeit.Name(),
			ClassYear: 3,
		})
		assert.NoError(t, err)
		password := gofakeit.Password(true, true, true, true, true, 30)
		user := &go_block.User{
			OptionalId: uuid.NewV4().String(),
			Email:      gofakeit.Email(),
			Password:   password,
			Image:      gofakeit.ImageURL(10, 10),
			Metadata:   string(metadata),
		}
		createdUser, err := userRepository.Create(context.Background(), user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		// set new encryption key
		encryptionKey, err := userRepository.crypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		userRepository.internalEncryptionKeys = append(userRepository.internalEncryptionKeys, encryptionKey)
		// act
		getUser, err := userRepository.Get(context.Background(), &go_block.User{
			OptionalId: createdUser.OptionalId,
		}, true)
		assert.NoError(t, err)
		assert.NotNil(t, getUser)
		assert.NoError(t, compareUsers(createdUser, getUser, false))
	}
}
