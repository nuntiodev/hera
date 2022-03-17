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

func TestUpdateOptionalId(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{})
	user.Id = ""
	namespace := uuid.NewV4().String()
	createUser, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      user,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	// act
	newOptionalId := uuid.NewV4().String()
	createUser.User.OptionalId = newOptionalId
	updateUser, err := testClient.UpdateOptionalId(ctx, &go_block.UserRequest{
		Update:    createUser.User,
		User:      createUser.User,
		Namespace: namespace,
	})
	// validate
	assert.NoError(t, err)
	assert.NotNil(t, updateUser)
	assert.NotNil(t, updateUser.User)
	assert.Equal(t, updateUser.User.OptionalId, newOptionalId)
}

func TestUpdateOptionalIdNoUser(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{})
	namespace := uuid.NewV4().String()
	_, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      user,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.UpdateOptionalId(ctx, &go_block.UserRequest{
		Namespace: namespace,
	})
	// validate
	assert.Error(t, err)
}

func TestUpdateOptionalIdNoReq(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	user.Id = ""
	namespace := uuid.NewV4().String()
	_, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      user,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.UpdateOptionalId(ctx, nil)
	// validate
	assert.Error(t, err)
}
