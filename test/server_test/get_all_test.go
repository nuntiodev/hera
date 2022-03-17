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

func TestGetAll(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userThree := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	createUserOne, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      userOne,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	createUserTwo, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      userTwo,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	_, err = testClient.Create(ctx, &go_block.UserRequest{
		User:      userThree,
		Namespace: "",
	})
	assert.NoError(t, err)
	// act
	getUsers, err := testClient.GetAll(ctx, &go_block.UserRequest{
		Namespace: namespace,
	})
	// validate
	assert.NoError(t, err)
	assert.NotNil(t, getUsers)
	assert.Equal(t, 2, len(getUsers.Users))
	assert.NoError(t, user_mock.CompareUsers(createUserOne.User, getUsers.Users[0]))
	assert.NoError(t, user_mock.CompareUsers(createUserTwo.User, getUsers.Users[1]))
}

func TestGetAllWithPartialEncryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userThree := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	createUserOne, err := testClient.Create(ctx, &go_block.UserRequest{
		User:          userOne,
		EncryptionKey: encryptionKey,
		Namespace:     namespace,
	})
	assert.NoError(t, err)
	createUserTwo, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      userTwo,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	_, err = testClient.Create(ctx, &go_block.UserRequest{
		User:          userThree,
		EncryptionKey: encryptionKey,
		Namespace:     namespace,
	})
	assert.NoError(t, err)
	// act
	getUsers, err := testClient.GetAll(ctx, &go_block.UserRequest{
		Namespace:     namespace,
		EncryptionKey: encryptionKey,
	})
	// validate
	assert.NoError(t, err)
	assert.NotNil(t, getUsers)
	assert.Equal(t, 3, len(getUsers.Users))
	assert.NoError(t, user_mock.CompareUsers(createUserOne.User, getUsers.Users[0]))
	assert.NoError(t, user_mock.CompareUsers(createUserTwo.User, getUsers.Users[1]))
}

func TestGetAllWithPartialEncryptionNoKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userThree := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	_, err := testClient.Create(ctx, &go_block.UserRequest{
		User:          userOne,
		EncryptionKey: encryptionKey,
		Namespace:     namespace,
	})
	assert.NoError(t, err)
	_, err = testClient.Create(ctx, &go_block.UserRequest{
		User:      userTwo,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	_, err = testClient.Create(ctx, &go_block.UserRequest{
		User:          userThree,
		EncryptionKey: encryptionKey,
		Namespace:     namespace,
	})
	assert.NoError(t, err)
	// act
	getAll, err := testClient.GetAll(ctx, &go_block.UserRequest{
		Namespace: namespace,
	})
	// validate
	assert.NoError(t, err)
	assert.Equal(t, 3, len(getAll.Users))
}

func TestGetAllNoReq(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	// act
	_, err := testClient.Get(ctx, nil)
	assert.Error(t, err)
}
