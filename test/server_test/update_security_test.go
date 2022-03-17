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

func TestSecurityProfile(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	user.Id = ""
	namespace := uuid.NewV4().String()
	createUser, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      user,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	// act
	updateUser, err := testClient.UpdateSecurity(ctx, &go_block.UserRequest{
		Update:    createUser.User,
		User:      createUser.User,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updateUser)
}

func TestUpdateSecurityNoUser(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	namespace := uuid.NewV4().String()
	_, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      user,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.UpdateSecurity(ctx, &go_block.UserRequest{})
	// validate
	assert.Error(t, err)
}

func TestUpdateSecurityNoReq(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	user.Id = ""
	_, err := testClient.Create(ctx, &go_block.UserRequest{
		User: user,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.UpdateSecurity(ctx, nil)
	// validate
	assert.Error(t, err)
}
