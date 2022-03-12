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

func TestUpdateImage(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
	})
	user.Id = ""
	createUser, err := testClient.Create(ctx, &go_block.UserRequest{
		User: user,
	})
	assert.NoError(t, err)
	// act
	newImage := gofakeit.ImageURL(20, 10)
	createUser.User.Image = newImage
	updateUser, err := testClient.UpdateImage(ctx, &go_block.UserRequest{
		Update: createUser.User,
		User:   createUser.User,
	})
	// validate
	assert.NoError(t, err)
	assert.NotNil(t, updateUser)
	assert.NotNil(t, updateUser.User)
	assert.Equal(t, updateUser.User.Image, newImage)
}

func TestUpdateImageWithEncryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
	})
	user.Id = ""
	createUser, err := testClient.Create(ctx, &go_block.UserRequest{
		User:          user,
		EncryptionKey: encryptionKey,
	})
	assert.NoError(t, err)
	// act
	newImage := gofakeit.ImageURL(20, 10)
	createUser.User.Image = newImage
	updateUser, err := testClient.UpdateImage(ctx, &go_block.UserRequest{
		Update:        createUser.User,
		User:          createUser.User,
		EncryptionKey: encryptionKey,
	})
	// validate
	assert.NoError(t, err)
	assert.NotNil(t, updateUser)
	assert.NotNil(t, updateUser.User)
	assert.Equal(t, updateUser.User.Image, newImage)
}

func TestUpdateImageNoUser(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
	})
	_, err := testClient.Create(ctx, &go_block.UserRequest{
		User: user,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.UpdateImage(ctx, &go_block.UserRequest{})
	// validate
	assert.Error(t, err)
}

func TestUpdateImageNoReq(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
	})
	user.Id = ""
	_, err := testClient.Create(ctx, &go_block.UserRequest{
		User: user,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.UpdateImage(ctx, nil)
	// validate
	assert.Error(t, err)
}
