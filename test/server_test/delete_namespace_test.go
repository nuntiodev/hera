package server_test

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDeleteNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	_, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      user_mock.GetRandomUser(&go_block.User{}),
		Namespace: namespace,
	})
	assert.NoError(t, err)
	_, err = testClient.Create(ctx, &go_block.UserRequest{
		User: user_mock.GetRandomUser(&go_block.User{}),
	})
	assert.NoError(t, err)
	newNamespace := uuid.NewV4().String()
	createUser, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      user_mock.GetRandomUser(&go_block.User{}),
		Namespace: newNamespace,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.DeleteNamespace(ctx, &go_block.UserRequest{
		Namespace: namespace,
	})
	// validate
	assert.NoError(t, err)
	// validate in database
	getUsers, err := testClient.GetAll(ctx, &go_block.UserRequest{
		Namespace: namespace,
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(getUsers.Users))
	getUsers, err = testClient.GetAll(ctx, &go_block.UserRequest{
		Namespace: newNamespace,
	})
	assert.Equal(t, 1, len(getUsers.Users))
	assert.NoError(t, user_mock.CompareUsers(getUsers.Users[0], createUser.User))
}

func TestDeleteNamespaceNoUser(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	// act
	_, err := testClient.DeleteNamespace(ctx, &go_block.UserRequest{})
	assert.Error(t, err)
}

func TestDeleteNamespaceNoUReq(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	// act
	_, err := testClient.DeleteNamespace(ctx, nil)
	assert.Error(t, err)
}
