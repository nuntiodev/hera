package user_repository

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/models"
	"testing"

	"github.com/nuntiodev/x/cryptox"
	"github.com/stretchr/testify/assert"
)

func TestGetByIdIEEncrypted(t *testing.T) {
	// setup available clients
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	err = clients[0].DeleteAll(context.Background())
	for _, userRepository := range clients {
		// create a user
		userOne := getTestUser()
		copy := userOne
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
		users := []*models.User{{Email: dbUserOne.Email}, {Id: dbUserOne.Id}, {Username: dbUserOne.Username}}
		for _, user := range users {
			getUser, err := userRepository.Get(context.Background(), models.UserToProtoUser(user))
			// validate
			assert.NoError(t, err)
			assert.NotNil(t, getUser)
			assert.NoError(t, compareUsers(models.ProtoUserToUser(&copy), getUser))
			// validate encryption level has been updated
			internalOne, externalOne := userRepository.crypto.EncryptionLevel(getUser)
			assert.Equal(t, int32(len(internalKeys)), internalOne)
			assert.Equal(t, int32(len(externalKeys)), externalOne)
		}
	}
}
