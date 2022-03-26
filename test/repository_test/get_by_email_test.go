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

func TestGetByEmail(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
		Email: gofakeit.Email(),
	})
	users, err := testRepository.Users(ctx, uuid.NewV4().String(), "")
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user)
	assert.Nil(t, err)
	// act
	getUser, err := users.Get(ctx, &go_block.User{
		Email: createdUser.Email,
	})
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, createdUser)
	assert.Nil(t, user_mock.CompareUsers(getUser, createdUser))
}

func TestGetByEmailWithEncryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
		Email: gofakeit.Email(),
	})
	users, err := testRepository.Users(ctx, uuid.NewV4().String(), encryptionKey)
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user)
	assert.Nil(t, err)
	// act
	getUser, err := users.Get(ctx, &go_block.User{
		Email: createdUser.Email,
	})
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, createdUser)
	assert.Nil(t, user_mock.CompareUsers(getUser, createdUser))
}

func TestGetByEmailWithInvalidEncryptionKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
		Email: gofakeit.Email(),
	})
	usersOne, err := testRepository.Users(ctx, uuid.NewV4().String(), encryptionKey)
	assert.NoError(t, err)
	_, err = usersOne.Create(ctx, user)
	assert.Nil(t, err)
	// act
	usersTwo, err := testRepository.Users(ctx, uuid.NewV4().String(), invalidEncryptionKey)
	assert.NoError(t, err)
	_, err = usersTwo.Get(ctx, &go_block.User{
		Email: user.Email,
	})
	// validate
	assert.Error(t, err)
}

func TestGetByEmailWithEncryptionNoDecryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
		Email: gofakeit.Email(),
	})
	users, err := testRepository.Users(ctx, "", encryptionKey)
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user)
	assert.Nil(t, err)
	// act
	getUser, err := users.Get(ctx, &go_block.User{
		Email: createdUser.Email,
	})
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUser)
}

func TestGetByEmailDifferentNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	usersOne, err := testRepository.Users(ctx, uuid.NewV4().String(), "")
	assert.NoError(t, err)
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	createdUser, err := usersOne.Create(ctx, user)
	assert.Nil(t, err)
	// act
	usersTwo, err := testRepository.Users(ctx, uuid.NewV4().String(), "")
	assert.NoError(t, err)
	getUser, err := usersTwo.Get(ctx, createdUser)
	assert.Error(t, err)
	// validate
	assert.Nil(t, getUser)
}
