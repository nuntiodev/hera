package user_repository

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/models"
	"testing"

	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/x/cryptox"
	"github.com/stretchr/testify/assert"
)

func TestGetAllIEEncrypted(t *testing.T) {
	// setup available clients
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	err = clients[0].DeleteAll(context.Background())
	assert.NoError(t, err)
	for index, userRepository := range clients {
		userRepository = clients[0]
		// create one user
		userOne := getTestUser()
		//copyOne := userOne
		dbUserOne, err := userRepository.Create(context.Background(), &userOne)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserOne)
		// create second user
		userTwo := getTestUser()
		//copyTwo := userTwo
		dbUserTwo, err := userRepository.Create(context.Background(), &userTwo)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserTwo)
		// create third user
		userThree := getTestUser()
		//copyThree := userThree
		dbUserThree, err := userRepository.Create(context.Background(), &userThree)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserThree)
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
		// act - get all users a couple of times
		for i := 0; i < 3; i++ {
			getUsers, err := userRepository.GetAll(context.Background(), nil)
			// validate
			assert.NoError(t, err)
			assert.NotNil(t, getUsers)
			assert.Equal(t, 3, len(getUsers), index)
			// validate order and user values
			assert.NoError(t, compareUsers(dbUserOne, getUsers[0]))
			assert.NoError(t, compareUsers(dbUserTwo, getUsers[1]))
			assert.NoError(t, compareUsers(dbUserThree, getUsers[2]))
			// validate encryption level
			internalOne, externalOne := userRepository.crypto.EncryptionLevel(getUsers[0])
			assert.Equal(t, int32(len(internalKeys)), internalOne)
			assert.Equal(t, int32(len(externalKeys)), externalOne)
			internalTwo, externalTwo := userRepository.crypto.EncryptionLevel(getUsers[1])
			assert.Equal(t, int32(len(internalKeys)), internalTwo)
			assert.Equal(t, int32(len(externalKeys)), externalTwo)
			internalThree, externalThree := userRepository.crypto.EncryptionLevel(getUsers[2])
			assert.Equal(t, int32(len(internalKeys)), internalThree)
			assert.Equal(t, int32(len(externalKeys)), externalThree)
		}
		// delete all at the end
		assert.NoError(t, userRepository.DeleteBatch(context.Background(), []*go_block.User{
			models.UserToProtoUser(dbUserOne),
			models.UserToProtoUser(dbUserTwo),
			models.UserToProtoUser(dbUserThree),
		}))
	}
}

func TestGetAllIEEncryptedWithFilters(t *testing.T) {
	// setup available clients
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	err = clients[0].DeleteAll(context.Background())
	assert.NoError(t, err)
	for index, userRepository := range clients {
		// create one user
		userOne := getTestUser()
		//copyOne := userOne
		dbUserOne, err := userRepository.Create(context.Background(), &userOne)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserOne)
		// create second user
		userTwo := getTestUser()
		//copyTwo := userTwo
		dbUserTwo, err := userRepository.Create(context.Background(), &userTwo)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserTwo)
		// create third user
		userThree := getTestUser()
		//copyThree := userThree
		dbUserThree, err := userRepository.Create(context.Background(), &userThree)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserThree)
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
		// act - get all users a couple of times
		for i := 0; i < 3; i++ {
			getUsers, err := userRepository.GetAll(context.Background(), &go_block.UserFilter{
				From:  0,
				To:    2,
				Order: go_block.UserFilter_DEC,
			}) // validate
			assert.NoError(t, err)
			assert.NotNil(t, getUsers)
			assert.Equal(t, 2, len(getUsers), index)
			// validate order and user values
			assert.NoError(t, compareUsers(dbUserTwo, getUsers[1]))
			assert.NoError(t, compareUsers(dbUserThree, getUsers[0]))
			// validate encryption level
			internalOne, externalOne := userRepository.crypto.EncryptionLevel(getUsers[0])
			assert.Equal(t, int32(len(internalKeys)), internalOne)
			assert.Equal(t, int32(len(externalKeys)), externalOne)
			internalTwo, externalTwo := userRepository.crypto.EncryptionLevel(getUsers[1])
			assert.Equal(t, int32(len(internalKeys)), internalTwo)
			assert.Equal(t, int32(len(externalKeys)), externalTwo)
		}
		// delete all at the end
		assert.NoError(t, userRepository.DeleteBatch(context.Background(), []*go_block.User{
			models.UserToProtoUser(dbUserOne),
			models.UserToProtoUser(dbUserTwo),
			models.UserToProtoUser(dbUserThree),
		}))
	}
}
