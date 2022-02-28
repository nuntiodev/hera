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

func TestUpdateNamespace(t *testing.T) {
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
	createUser, err := testClient.Create(ctx, &block_user.UserRequest{
		User: user,
	})
	assert.NoError(t, err)
	// act
	newNamespace := uuid.NewV4().String()
	createUser.User.Namespace = newNamespace
	updateUser, err := testClient.UpdateNamespace(ctx, &block_user.UserRequest{
		Update: createUser.User,
	})
	assert.NoError(t, err)
	// validate
	assert.Equal(t, newNamespace, updateUser.User.Namespace)
}

func TestUpdateNamespaceNoUser(t *testing.T) {
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
	_, err := testClient.Create(ctx, &block_user.UserRequest{
		User: user,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.UpdateNamespace(ctx, &block_user.UserRequest{})
	// validate
	assert.Error(t, err)
}

func TestUpdateNamespaceNoReq(t *testing.T) {
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
	_, err := testClient.Create(ctx, &block_user.UserRequest{
		User: user,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.UpdateNamespace(ctx, nil)
	// validate
	assert.Error(t, err)
}
