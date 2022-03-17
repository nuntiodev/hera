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

func TestDelete(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	createUser, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      user,
		Namespace: namespace,
	})
	// act
	_, err = testClient.Delete(ctx, &go_block.UserRequest{
		User:      createUser.User,
		Namespace: namespace,
	})
	// validate
	assert.NoError(t, err)
	// validate in database
	_, err = testClient.Get(ctx, &go_block.UserRequest{
		User:      createUser.User,
		Namespace: namespace,
	})
	assert.Error(t, err)
}

func TestDeleteNoUser(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	// act
	_, err := testClient.Delete(ctx, &go_block.UserRequest{
		Namespace: uuid.NewV4().String(),
	})
	assert.Error(t, err)
}

func TestDeleteNoUReq(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	// act
	_, err := testClient.Delete(ctx, nil)
	assert.Error(t, err)
}
