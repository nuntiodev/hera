package user_repository

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/models"
	"testing"

	"github.com/nuntiodev/block-proto/go_block"
	"github.com/stretchr/testify/assert"
)

func TestDeleteBatchIEEncrypted(t *testing.T) {
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	err = clients[0].DeleteAll(context.Background())
	assert.NoError(t, err)
	for _, userRepository := range clients {
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
		err = userRepository.DeleteBatch(context.Background(), []*go_block.User{models.UserToProtoUser(dbUserOne), models.UserToProtoUser(dbUserTwo)})
		// validate
		assert.NoError(t, err)
		// validate
		_, err = userRepository.Get(context.Background(), models.UserToProtoUser(dbUserOne))
		assert.Error(t, err)
		_, err = userRepository.Get(context.Background(), models.UserToProtoUser(dbUserTwo))
		assert.Error(t, err)
		getUser, err := userRepository.Get(context.Background(), models.UserToProtoUser(dbUserThree))
		assert.NoError(t, err)
		assert.NotNil(t, getUser)
	}
}
