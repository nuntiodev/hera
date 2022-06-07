package user_repository

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/models"
	"testing"

	"github.com/google/uuid"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/x/cryptox"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUsernameIEEncrypted(t *testing.T) {
	// setup available clients
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	err = clients[0].DeleteAll(context.Background())
	assert.NoError(t, err)
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
		newUsername := uuid.NewString()
		updatedUser, err := userRepository.UpdateUsername(context.Background(), models.UserToProtoUser(dbUserOne), &go_block.User{
			Username: newUsername,
		})
		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		assert.Equal(t, newUsername, updatedUser.Username.Body)
		// validate change has been updated in db
		getUser, err := userRepository.Get(context.Background(), models.UserToProtoUser(updatedUser))
		assert.NoError(t, err)
		assert.Equal(t, newUsername, getUser.Username.Body)
		// assert.NoError(t, compareUsers(getUser, updatedUser, true)) todo: return a valid new state of user
		// validate encryption levels
		internalThree, externalThree := userRepository.crypto.EncryptionLevel(getUser)
		assert.Equal(t, int32(len(internalKeys)), internalThree)
		assert.Equal(t, int32(len(externalKeys)), externalThree)
	}
}

func TestUpdateUsernameNilUpdate(t *testing.T) {
	// setup available clients
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	err = clients[0].DeleteAll(context.Background())
	assert.NoError(t, err)
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
		updatedUser, err := userRepository.UpdateUsername(context.Background(), models.UserToProtoUser(dbUserOne), nil)
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	}
}
