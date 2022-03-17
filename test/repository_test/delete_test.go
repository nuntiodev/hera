package respository_test

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-user-service/repository/user_repository"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDelete(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	users, err := testRepository.Users(ctx, "")
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, "")
	assert.Nil(t, err)
	// act
	err = users.Delete(ctx, createdUser)
	assert.Nil(t, err)
	// validate
	getUser, err := users.Get(ctx, createdUser, "")
	assert.Error(t, err)
	assert.Nil(t, getUser)
	// validate in database
	_, err = users.Get(ctx, createdUser, "")
	assert.Error(t, err)
}

func TestDeleteDifferentNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	usersOne, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser, err := usersOne.Create(ctx, user, "")
	assert.Nil(t, err)
	// act
	usersTwo, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	err = usersTwo.Delete(ctx, createdUser)
	assert.Error(t, err)
	assert.Equal(t, user_repository.NoUsersDeletedErr, err)
	// validate
	getUser, err := usersOne.Get(ctx, createdUser, "")
	assert.Nil(t, err)
	assert.Nil(t, user_mock.CompareUsers(getUser, createdUser))
}
