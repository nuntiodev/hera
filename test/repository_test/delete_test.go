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
	createdUser, err := testRepo.Create(ctx, user)
	assert.Nil(t, err)
	// act
	err = testRepo.Delete(ctx, createdUser)
	assert.Nil(t, err)
	// validate
	getUser, err := testRepo.GetById(ctx, createdUser)
	assert.Error(t, err)
	assert.Nil(t, getUser)
	// validate in database
	_, err = testRepo.GetById(ctx, createdUser)
	assert.Error(t, err)
}

func TestDeleteDifferentNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	namespace := user.Namespace
	createdUser, err := testRepo.Create(ctx, user)
	assert.Nil(t, err)
	// act
	createdUser.Namespace = uuid.NewV4().String()
	err = testRepo.Delete(ctx, createdUser)
	assert.Error(t, err)
	assert.Equal(t, user_repository.NoUsersDeletedErr, err)
	// validate
	createdUser.Namespace = namespace
	getUser, err := testRepo.GetById(ctx, createdUser)
	assert.Nil(t, err)
	assert.Nil(t, user_mock.CompareUsers(getUser, createdUser))
}
