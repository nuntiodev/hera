package user_repository

import (
	"context"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/x/cryptox"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateEmailIEEncrypted(t *testing.T) {
	// setup available clients
	var clients []*mongoRepository
	/*
		userRepositoryFullEncryption, err := getTestUserRepository(context.Background(), true, true)
		assert.NoError(t, err)
		userRepositoryInternalEncryption, err := getTestUserRepository(context.Background(), true, false)
		assert.NoError(t, err)
		userRepositoryExternalEncryption, err := getTestUserRepository(context.Background(), false, true)
		assert.NoError(t, err)
	*/
	userRepositoryNoEncryption, err := getTestUserRepository(context.Background(), false, false, "")
	assert.NoError(t, err)
	clients = []*mongoRepository{userRepositoryNoEncryption}
	for index, userRepository := range clients {
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
		newEmail := gofakeit.Email()
		updatedUser, err := userRepository.UpdateEmail(context.Background(), createdUser, &go_block.User{
			Email: newEmail,
		})
		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		assert.Equal(t, newEmail, updatedUser.Email)
		// validate change has been updated in db
		getUser, err := userRepository.Get(context.Background(), updatedUser, true)
		assert.NoError(t, err, index)
		assert.Equal(t, newEmail, getUser.Email)
	}
}

func TestUpdateEmailInvalidEmail(t *testing.T) {
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
		newEmail := "info@@softcorp.io"
		updatedUser, err := userRepository.UpdateEmail(context.Background(), createdUser, &go_block.User{
			Email: newEmail,
		})
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	}
}

func TestUpdateEmailNilUpdate(t *testing.T) {
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
		updatedUser, err := userRepository.UpdateEmail(context.Background(), createdUser, nil)
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	}
}

func TestUpdateEmailNilGet(t *testing.T) {
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
		updatedUser, err := userRepository.UpdateEmail(context.Background(), nil, createdUser)
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	}
}
