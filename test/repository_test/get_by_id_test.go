package respository_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetById(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
	})
	createdUser, err := testRepo.Create(ctx, user, nil)
	assert.Nil(t, err)
	// act
	getUser, err := testRepo.Get(ctx, &block_user.User{
		Id:        createdUser.Id,
		Namespace: createdUser.Namespace,
	}, nil)
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, createdUser)
	assert.Nil(t, user_mock.CompareUsers(getUser, createdUser))
}

func TestGetByIdDifferentNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
	})
	createdUser, err := testRepo.Create(ctx, user, nil)
	assert.Nil(t, err)
	// act
	createdUser.Namespace = ""
	getUser, err := testRepo.Get(ctx, createdUser, nil)
	assert.Error(t, err)
	// validate
	assert.Nil(t, getUser)
}
