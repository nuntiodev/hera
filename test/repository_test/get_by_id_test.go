package respository_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetById(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, "")
	assert.Nil(t, err)
	// act
	getUser, err := users.Get(ctx, &go_block.User{
		Id: createdUser.Id,
	}, "")
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, createdUser)
	assert.Nil(t, user_mock.CompareUsers(getUser, createdUser))
}

func TestGetByIdDifferentNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, "")
	assert.Nil(t, err)
	// act
	usersTwo, err := testRepository.Users(ctx, "")
	assert.NoError(t, err)
	getUser, err := usersTwo.Get(ctx, createdUser, "")
	assert.Error(t, err)
	// validate
	assert.Nil(t, getUser)
}
