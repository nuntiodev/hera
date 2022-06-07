package user_repository

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/models"
	"testing"

	"github.com/nuntiodev/block-proto/go_block"
	"github.com/stretchr/testify/assert"
)

/*
	TestCountIEEncrypted - this method creates three users and tests that the number of users in the database is 3.
	It does so under multiple different configurations of the crypto object.
*/
func TestCountIEEncrypted(t *testing.T) {
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	for index, userRepository := range clients {
		// create the first user
		userOne := getTestUser()
		dbUserOne, err := userRepository.Create(context.Background(), &userOne)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserOne)
		// create the second user
		userTwo := getTestUser()
		dbUserTwo, err := userRepository.Create(context.Background(), &userTwo)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserTwo)
		// create the third user
		userThree := getTestUser()
		dbUserThree, err := userRepository.Create(context.Background(), &userThree)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserThree)
		// act
		count, err := userRepository.Count(context.Background())
		// validate
		assert.NoError(t, err)
		assert.Equal(t, 3, int(count), index)
		// delete all at the end
		assert.NoError(t, userRepository.DeleteBatch(context.Background(), []*go_block.User{
			models.UserToProtoUser(dbUserOne),
			models.UserToProtoUser(dbUserTwo),
			models.UserToProtoUser(dbUserThree),
		}))
	}
}
