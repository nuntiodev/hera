package server_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	user.Id = ""
	// act
	createUser, err := testClient.Create(ctx, &block_user.UserRequest{
		User: user,
	})
	assert.NoError(t, err)
	// validate
	getUser, err := testClient.Get(ctx, &block_user.UserRequest{
		User: createUser.User,
	})
	assert.NoError(t, err)
	assert.NotNil(t, getUser)
	assert.NoError(t, user_mock.CompareUsers(getUser.User, createUser.User))
}

func TestGetWithEncryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	user.Id = ""
	// act
	createUser, err := testClient.Create(ctx, &block_user.UserRequest{
		User:           user,
		WithEncryption: true,
		EncryptionKey:  encryptionKey,
	})
	assert.NoError(t, err)
	// validate
	getUser, err := testClient.Get(ctx, &block_user.UserRequest{
		User:          createUser.User,
		EncryptionKey: encryptionKey,
	})
	assert.NoError(t, err)
	assert.NotNil(t, getUser)
	assert.NoError(t, user_mock.CompareUsers(getUser.User, createUser.User))
}

func TestGetNoUser(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	// act
	_, err := testClient.Get(ctx, &block_user.UserRequest{})
	assert.Error(t, err)
}

func TestGetNoUReq(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	// act
	_, err := testClient.Get(ctx, nil)
	assert.Error(t, err)
}
