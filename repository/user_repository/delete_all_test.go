package user_repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteAllIEEncrypted(t *testing.T) {
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
		// create user two
		userTwo := getTestUser()
		dbUserTwo, err := userRepository.Create(context.Background(), &userTwo)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserTwo)
		// create user three
		userThree := getTestUser()
		dbUserThree, err := userRepository.Create(context.Background(), &userThree)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserThree)
		// act
		err = userRepository.DeleteAll(context.Background())
		// validate
		assert.NoError(t, err)
		// validate repository is not empty
		getUsers, err := userRepository.GetAll(context.Background(), nil)
		assert.Nil(t, getUsers)
		assert.Equal(t, 0, len(getUsers), index)
	}
}
