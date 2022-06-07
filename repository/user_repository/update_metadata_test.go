package user_repository

import (
	"context"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/x/cryptox"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateMetadataIEEncrypted(t *testing.T) {
	// setup available clients
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	err = clients[0].DeleteAll(context.Background())
	for _, userRepository := range clients {
		userOne := getTestUser()
		dbUserOne, err := userRepository.Create(context.Background(), &userOne)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserOne)
		// set new internal and external encryption key
		encryptionKey, err := cryptox.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		// internal
		internalKeys, _ := userRepository.crypto.GetInternalEncryptionKeys()
		internalKeys = append(internalKeys, encryptionKey)
		assert.NoError(t, userRepository.crypto.SetInternalEncryptionKeys(internalKeys))
		// external
		externalKeys, _ := userRepository.crypto.GetExternalEncryptionKeys()
		externalKeys = append(externalKeys, encryptionKey)
		assert.NoError(t, userRepository.crypto.SetExternalEncryptionKeys(externalKeys))
		// act
		newMetadata, err := json.Marshal(CustomMetadata{
			Name: gofakeit.Name(),
		})
		assert.NoError(t, err)
		dbUpdatedUser, err := userRepository.UpdateMetadata(context.Background(), models.UserToProtoUser(dbUserOne), &go_block.User{
			Metadata: string(newMetadata),
		})
		assert.NoError(t, err)
		assert.NotNil(t, dbUpdatedUser)
		assert.Equal(t, string(newMetadata), dbUpdatedUser.Metadata.Body)
		// validate change has been updated in db
		getUser, err := userRepository.Get(context.Background(), models.UserToProtoUser(dbUpdatedUser))
		assert.NoError(t, err)
		assert.Equal(t, string(newMetadata), getUser.Metadata.Body)
		// validate encryption level
		internalThree, externalThree := userRepository.crypto.EncryptionLevel(getUser)
		assert.Equal(t, int32(len(internalKeys)), internalThree)
		assert.Equal(t, int32(len(externalKeys)), externalThree)
	}
}

func TestUpdateMetadataInvalid(t *testing.T) {
	// setup available clients
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	err = clients[0].DeleteAll(context.Background())
	for _, userRepository := range clients {
		userOne := getTestUser()
		dbUserOne, err := userRepository.Create(context.Background(), &userOne)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserOne)
		// set new internal and external encryption key
		encryptionKey, err := cryptox.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		// internal
		internalKeys, _ := userRepository.crypto.GetInternalEncryptionKeys()
		internalKeys = append(internalKeys, encryptionKey)
		assert.NoError(t, userRepository.crypto.SetInternalEncryptionKeys(internalKeys))
		// external
		externalKeys, _ := userRepository.crypto.GetExternalEncryptionKeys()
		externalKeys = append(externalKeys, encryptionKey)
		assert.NoError(t, userRepository.crypto.SetExternalEncryptionKeys(externalKeys))
		// act
		newMetadata := "invalid metadata"
		dbUpdatedUser, err := userRepository.UpdateMetadata(context.Background(), models.UserToProtoUser(dbUserOne), &go_block.User{
			Metadata: newMetadata,
		})
		assert.Error(t, err)
		assert.Nil(t, dbUpdatedUser)
	}
}

func TestUpdateMetadataNilUpdate(t *testing.T) {
	// setup available clients
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	err = clients[0].DeleteAll(context.Background())
	for _, userRepository := range clients {
		userOne := getTestUser()
		dbUserOne, err := userRepository.Create(context.Background(), &userOne)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserOne)
		// set new internal and external encryption key
		encryptionKey, err := cryptox.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		// internal
		internalKeys, _ := userRepository.crypto.GetInternalEncryptionKeys()
		internalKeys = append(internalKeys, encryptionKey)
		assert.NoError(t, userRepository.crypto.SetInternalEncryptionKeys(internalKeys))
		// external
		externalKeys, _ := userRepository.crypto.GetExternalEncryptionKeys()
		externalKeys = append(externalKeys, encryptionKey)
		assert.NoError(t, userRepository.crypto.SetExternalEncryptionKeys(externalKeys))
		// act
		dbUpdatedUser, err := userRepository.UpdateMetadata(context.Background(), models.UserToProtoUser(dbUserOne), nil)
		assert.Error(t, err)
		assert.Nil(t, dbUpdatedUser)
	}
}
