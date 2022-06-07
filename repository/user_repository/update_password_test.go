package user_repository

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/x/cryptox"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUpdatePasswordIEEncryptedById(t *testing.T) {
	// setup available clients
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	err = clients[0].DeleteAll(context.Background())
	assert.NoError(t, err)
	for _, userRepository := range clients {
		// create user
		userOne := getTestUser()
		dbUserOne, err := userRepository.Create(context.Background(), &userOne)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserOne)
		// act
		newPassword := "My@Secure3NewPassword1234!"
		dbUserOne.Password = newPassword
		dbUserOne.Username = cryptox.Stringx{}
		dbUserOne.Email = cryptox.Stringx{}
		updatedUser, err := userRepository.UpdatePassword(context.Background(), models.UserToProtoUser(dbUserOne), models.UserToProtoUser(dbUserOne))
		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		// validate updated fields
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte(newPassword)))
		// validate change has been updated in db
		getUser, err := userRepository.Get(context.Background(), models.UserToProtoUser(updatedUser))
		assert.NoError(t, err)
		assert.Equal(t, updatedUser.Password, getUser.Password)
		// assert.NoError(t, compareUsers(getUser, updatedUser, true)) todo: return a valid new state of user
	}
}

func TestUpdatePasswordIEEncryptedByEmail(t *testing.T) {
	// setup available clients
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	err = clients[0].DeleteAll(context.Background())
	assert.NoError(t, err)
	for _, userRepository := range clients {
		// create user
		userOne := getTestUser()
		dbUserOne, err := userRepository.Create(context.Background(), &userOne)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserOne)
		// act
		newPassword := "My@Secure3NewPassword1234!"
		dbUserOne.Password = newPassword
		dbUserOne.Username = cryptox.Stringx{}
		dbUserOne.Id = ""
		updatedUser, err := userRepository.UpdatePassword(context.Background(), models.UserToProtoUser(dbUserOne), models.UserToProtoUser(dbUserOne))
		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		// validate updated fields
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte(newPassword)))
	}
}

func TestUpdatePasswordIEEncryptedByUsername(t *testing.T) {
	// setup available clients
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	err = clients[0].DeleteAll(context.Background())
	assert.NoError(t, err)
	for _, userRepository := range clients {
		// create user
		userOne := getTestUser()
		dbUserOne, err := userRepository.Create(context.Background(), &userOne)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserOne)
		// act
		newPassword := "My@Secure3NewPassword1234!"
		dbUserOne.Password = newPassword
		dbUserOne.Email = cryptox.Stringx{}
		dbUserOne.Id = ""
		updatedUser, err := userRepository.UpdatePassword(context.Background(), models.UserToProtoUser(dbUserOne), models.UserToProtoUser(dbUserOne))
		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		// validate updated fields
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte(newPassword)))
	}
}

func TestUpdatePasswordWeakPassword(t *testing.T) {
	// setup available clients
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	err = clients[0].DeleteAll(context.Background())
	assert.NoError(t, err)
	for _, userRepository := range clients {
		// create user
		userOne := getTestUser()
		dbUserOne, err := userRepository.Create(context.Background(), &userOne)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserOne)
		// act
		newPassword := "newpassword"
		dbUserOne.Password = newPassword
		updatedUser, err := userRepository.UpdatePassword(context.Background(), models.UserToProtoUser(dbUserOne), models.UserToProtoUser(dbUserOne))
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	}
}

func TestUpdatePasswordNoUpdate(t *testing.T) {
	// setup available clients
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	err = clients[0].DeleteAll(context.Background())
	assert.NoError(t, err)
	for _, userRepository := range clients {
		// create user
		userOne := getTestUser()
		dbUserOne, err := userRepository.Create(context.Background(), &userOne)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserOne)
		// act
		dbUserOne.Id = ""
		dbUserOne.Email = cryptox.Stringx{}
		dbUserOne.Username = cryptox.Stringx{}
		updatedUser, err := userRepository.UpdatePassword(context.Background(), models.UserToProtoUser(dbUserOne), models.UserToProtoUser(dbUserOne))
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	}
}

func TestUpdatePasswordNilUpdate(t *testing.T) {
	// setup available clients
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	err = clients[0].DeleteAll(context.Background())
	assert.NoError(t, err)
	for _, userRepository := range clients {
		// create user
		userOne := getTestUser()
		dbUserOne, err := userRepository.Create(context.Background(), &userOne)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserOne)
		// act
		updatedUser, err := userRepository.UpdatePassword(context.Background(), models.UserToProtoUser(dbUserOne), nil)
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	}
}

func TestUpdatePasswordNilGet(t *testing.T) {
	// setup available clients
	clients, err := getUserRepositories()
	assert.NoError(t, err)
	// delete all users from other tests (we use the same collection)
	err = clients[0].DeleteAll(context.Background())
	assert.NoError(t, err)
	for _, userRepository := range clients {
		// create user
		userOne := getTestUser()
		dbUserOne, err := userRepository.Create(context.Background(), &userOne)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserOne)
		// act
		updatedUser, err := userRepository.UpdatePassword(context.Background(), nil, models.UserToProtoUser(dbUserOne))
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	}
}
