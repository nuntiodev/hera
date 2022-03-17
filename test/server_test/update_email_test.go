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

func TestUpdateEmail(t *testing.T) {
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
	newEmail := gofakeit.Email()
	createUser.User.Email = newEmail
	updateUser, err := testClient.UpdateEmail(ctx, &go_block.UserRequest{
		Update:    createUser.User,
		User:      createUser.User,
		Namespace: namespace,
	})
	// validate
	assert.NoError(t, err)
	assert.NotNil(t, updateUser)
	assert.NotNil(t, updateUser.User)
	assert.Equal(t, updateUser.User.Email, newEmail)
}

func TestUpdateEmailWithEncryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	user.Id = ""
	namespace := uuid.NewV4().String()
	createUser, err := testClient.Create(ctx, &go_block.UserRequest{
		User:          user,
		EncryptionKey: encryptionKey,
		Namespace:     namespace,
	})
	assert.NoError(t, err)
	// act
	newEmail := gofakeit.Email()
	createUser.User.Email = newEmail
	updateUser, err := testClient.UpdateEmail(ctx, &go_block.UserRequest{
		Update:        createUser.User,
		User:          createUser.User,
		EncryptionKey: encryptionKey,
		Namespace:     namespace,
	})
	// validate
	assert.NoError(t, err)
	assert.NotNil(t, updateUser)
	assert.NotNil(t, updateUser.User)
	assert.Equal(t, updateUser.User.Email, newEmail)
}

func TestUpdateEmailNoUser(t *testing.T) {
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
	_, err = testClient.UpdateEmail(ctx, &go_block.UserRequest{
		Namespace: namespace,
	})
	// validate
	assert.Error(t, err)
}

func TestUpdateEmailNoReq(t *testing.T) {
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
	_, err = testClient.UpdateEmail(ctx, nil)
	// validate
	assert.Error(t, err)
}
