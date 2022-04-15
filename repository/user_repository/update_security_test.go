package user_repository

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/nuntiodev/block-proto/go_block"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/x/cryptox"
	"github.com/stretchr/testify/assert"
)

func TestUpdateSecurityIEEncrypted(t *testing.T) {
	// setup available clients
	var clients []*mongodbRepository
	userRepositoryFullEncryption, err := getTestUserRepository(context.Background(), true, true, "")
	assert.NoError(t, err)
	userRepositoryExternalEncryption, err := getTestUserRepository(context.Background(), false, true, "")
	assert.NoError(t, err)
	clients = []*mongodbRepository{userRepositoryFullEncryption, userRepositoryExternalEncryption}
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
		updatedUser, err := userRepository.UpdateSecurity(context.Background(), createdUser)
		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		assert.False(t, updatedUser.ExternalEncrypted)
		// assert that update has been propagated correctly to database
		getUser, err := userRepository.Get(context.Background(), updatedUser, true)
		assert.NoError(t, err)
		assert.NotNil(t, getUser)
	}
}

func TestUpdateSecurityUnencryptedUser(t *testing.T) {
	userRepository, err := getTestUserRepository(context.Background(), true, false, "")
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	userRepository.externalEncryptionKey = encryptionKey
	// act
	updatedUser, err := userRepository.UpdateSecurity(context.Background(), createdUser)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.True(t, updatedUser.ExternalEncrypted)
}

func TestUpdateSecurityNilUpdate(t *testing.T) {
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
		updatedUser, err := userRepository.UpdateSecurity(context.Background(), nil)
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	}
}
