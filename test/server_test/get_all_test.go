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

func TestGetAll(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userThree := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	createUserOne, err := testClient.Create(ctx, &block_user.UserRequest{
		User: userOne,
	})
	assert.NoError(t, err)
	createUserTwo, err := testClient.Create(ctx, &block_user.UserRequest{
		User: userTwo,
	})
	assert.NoError(t, err)
	_, err = testClient.Create(ctx, &block_user.UserRequest{
		User: userThree,
	})
	assert.NoError(t, err)
	// act
	getUsers, err := testClient.GetAll(ctx, &block_user.UserRequest{
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
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userThree := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	createUserOne, err := testClient.Create(ctx, &block_user.UserRequest{
		User:           userOne,
		WithEncryption: true,
		EncryptionKey:  encryptionKey,
	})
	assert.NoError(t, err)
	createUserTwo, err := testClient.Create(ctx, &block_user.UserRequest{
		User: userTwo,
	})
	assert.NoError(t, err)
	_, err = testClient.Create(ctx, &block_user.UserRequest{
		User:           userThree,
		WithEncryption: true,
		EncryptionKey:  encryptionKey,
	})
	assert.NoError(t, err)
	// act
	getUsers, err := testClient.GetAll(ctx, &block_user.UserRequest{
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
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userThree := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	_, err := testClient.Create(ctx, &block_user.UserRequest{
		User:           userOne,
		WithEncryption: true,
		EncryptionKey:  encryptionKey,
	})
	assert.NoError(t, err)
	_, err = testClient.Create(ctx, &block_user.UserRequest{
		User: userTwo,
	})
	assert.NoError(t, err)
	_, err = testClient.Create(ctx, &block_user.UserRequest{
		User:           userThree,
		WithEncryption: true,
		EncryptionKey:  encryptionKey,
	})
	assert.NoError(t, err)
	// act
	getAll, err := testClient.GetAll(ctx, &block_user.UserRequest{
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
