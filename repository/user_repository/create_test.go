package user_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateIEEncrypted(t *testing.T) {
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	for _, userRepository := range clients {
		user := getTestUser()
		c := user
		// act
		createdUser, err := userRepository.Create(context.Background(), &user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		// assert new fields are present
		assert.NotEmpty(t, createdUser.Id)
		// assert that old fields are the same
		assert.Equal(t, createdUser.Email.Body, c.Email)
		assert.NotEqual(t, user.Password, c.Password)
		assert.Equal(t, createdUser.Image.Body, c.Image)
		assert.Equal(t, createdUser.Metadata.Body, c.Metadata)
		assert.Equal(t, createdUser.FirstName.Body, c.FirstName)
		assert.Equal(t, createdUser.LastName.Body, c.LastName)
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte(c.Password)))
		// delete all at the end
		assert.NoError(t, userRepository.DeleteBatch(context.Background(), []*go_block.User{
			models.UserToProtoUser(createdUser),
		}))
	}
}

func TestCreateInvalidPassword(t *testing.T) {
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	for _, userRepository := range clients {
		user := getTestUser()
		user.Password = "test1234"
		// act
		createdUser, err := userRepository.Create(context.Background(), &user)
		assert.Error(t, err)
		assert.Nil(t, createdUser)
	}
}

func TestCreateInvalidEmail(t *testing.T) {
	// setup user client
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	for _, userRepository := range clients {
		user := getTestUser()
		user.Email = "info@@nuntio.io"
		// act
		createdUser, err := userRepository.Create(context.Background(), &user)
		assert.Error(t, err)
		assert.Nil(t, createdUser)
	}
}

func TestCreateInvalidMetadata(t *testing.T) {
	// setup user client
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	for _, userRepository := range clients {
		user := getTestUser()
		user.Metadata = "some invalid metadata"
		// act
		createdUser, err := userRepository.Create(context.Background(), &user)
		assert.Error(t, err)
		assert.Nil(t, createdUser)
	}
}

func TestCreateNilUser(t *testing.T) {
	// setup user client
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	for _, userRepository := range clients {
		// act
		createdUser, err := userRepository.Create(context.Background(), nil)
		assert.Error(t, err)
		assert.Nil(t, createdUser)
	}
}
