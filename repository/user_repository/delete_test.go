package user_repository

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteIEEncrypted(t *testing.T) {
	// setup available clients
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	assert.NoError(t, clients[0].DeleteAll(context.Background()))
	for index, userRepository := range clients {
		// create user one
		userOne := getTestUser()
		dbUserOne, err := userRepository.Create(context.Background(), &userOne)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserOne)
		// act
		err = userRepository.Delete(context.Background(), models.UserToProtoUser(dbUserOne))
		// validate
		assert.NoError(t, err)
		// validate repository is not empty
		getUser, err := userRepository.Get(context.Background(), models.UserToProtoUser(dbUserOne))
		assert.Error(t, err, index)
		assert.Nil(t, getUser)
	}
}
