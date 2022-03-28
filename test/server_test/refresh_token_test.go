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

func TestRefreshToken(t *testing.T) {
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
	loginUser, err := testClient.Login(ctx, &go_block.UserRequest{
		User: &go_block.User{
			Password: password,
			Id:       createdUser.User.Id,
		},
		Namespace: namespace,
	})
	assert.NoError(t, err)
	// act
	refreshToken, err := testClient.RefreshToken(ctx, &go_block.UserRequest{
		Token: &go_block.Token{
			RefreshToken: loginUser.Token.RefreshToken,
		},
	})
	// validate
	assert.NoError(t, err)
	assert.NotNil(t, refreshToken)
	assert.NotNil(t, refreshToken.Token)
	// validate new access token can make request
	validateUser, err := testClient.ValidateToken(ctx, &go_block.UserRequest{
		Token: &go_block.Token{
			AccessToken: refreshToken.Token.AccessToken,
		},
	})
	// validate
	assert.NoError(t, err)
	assert.NotNil(t, validateUser.User)
	assert.Equal(t, loginUser.User.Id, createdUser.User.Id)
}

func TestRefreshTokenExpiry(t *testing.T) {
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
	loginUser, err := testClient.Login(ctx, &go_block.UserRequest{
		User: &go_block.User{
			Password: password,
			Id:       createdUser.User.Id,
		},
		Namespace: namespace,
	})
	assert.NoError(t, err)
	// act
	time.Sleep(refreshTokenExpiresAt + 1*time.Second)
	_, err = testClient.RefreshToken(ctx, &go_block.UserRequest{
		Token: &go_block.Token{
			RefreshToken: loginUser.Token.RefreshToken,
		},
	})
	// validate
	assert.Error(t, err)
}
