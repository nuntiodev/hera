package server_test

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

func TestLogin(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	password := "Test1234"
	user := user_mock.GetRandomUser(&go_block.User{
		Image:    gofakeit.ImageURL(10, 10),
		Password: password,
	})
	user.Id = ""
	namespace := uuid.NewV4().String()
	createdUser, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      user,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	// act
	loginUser, err := testClient.Login(ctx, &go_block.UserRequest{
		User: &go_block.User{
			Password: password,
			Id:       createdUser.User.Id,
		},
		Namespace: namespace,
	})
	// validate
	assert.NoError(t, err)
	assert.NotNil(t, loginUser)
	assert.NotNil(t, loginUser.Token)
	assert.NotEmpty(t, loginUser.Token.RefreshToken)
	assert.NotEmpty(t, loginUser.Token.AccessToken)
}

func TestLoginNoPassword(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	password := "Test1234"
	user := user_mock.GetRandomUser(&go_block.User{
		Image:    gofakeit.ImageURL(10, 10),
		Password: password,
	})
	user.Id = ""
	namespace := uuid.NewV4().String()
	createdUser, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      user,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.Login(ctx, &go_block.UserRequest{
		User: &go_block.User{
			Id: createdUser.User.Id,
		},
		Namespace: namespace,
	})
	// validate
	assert.Error(t, err)
}

func TestLoginInvalidPassword(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	password := "Test1234"
	user := user_mock.GetRandomUser(&go_block.User{
		Image:    gofakeit.ImageURL(10, 10),
		Password: password,
	})
	user.Id = ""
	namespace := uuid.NewV4().String()
	createdUser, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      user,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.Login(ctx, &go_block.UserRequest{
		User: &go_block.User{
			Id:       createdUser.User.Id,
			Password: password + "2",
		},
		Namespace: namespace,
	})
	// validate
	assert.Error(t, err)
}
