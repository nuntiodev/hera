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

func TestUpdateMetadata(t *testing.T) {
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
	newMetadata := user_mock.GetMetadata(nil)
	createUser.User.Metadata = newMetadata
	updateUser, err := testClient.UpdateMetadata(ctx, &go_block.UserRequest{
		Update:    createUser.User,
		User:      createUser.User,
		Namespace: namespace,
	})
	// validate
	assert.NoError(t, err)
	assert.NotNil(t, updateUser)
	assert.NotNil(t, updateUser.User)
	assert.Equal(t, updateUser.User.Metadata, newMetadata)
}

func TestUpdateMetadataWithEncryption(t *testing.T) {
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
	newMetadata := user_mock.GetMetadata(nil)
	createUser.User.Metadata = newMetadata
	updateUser, err := testClient.UpdateMetadata(ctx, &go_block.UserRequest{
		Update:        createUser.User,
		User:          createUser.User,
		EncryptionKey: encryptionKey,
		Namespace:     namespace,
	})
	// validate
	assert.NoError(t, err)
	assert.NotNil(t, updateUser)
	assert.NotNil(t, updateUser.User)
	assert.Equal(t, updateUser.User.Metadata, newMetadata)
}

func TestUpdateMetadataNoUser(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	_, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      user,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.UpdateMetadata(ctx, &go_block.UserRequest{})
	// validate
	assert.Error(t, err)
}

func TestUpdateMetadataNoReq(t *testing.T) {
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
	_, err = testClient.UpdateMetadata(ctx, nil)
	// validate
	assert.Error(t, err)
}
