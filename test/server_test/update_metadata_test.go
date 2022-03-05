package server_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUpdateMetadata(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
	})
	user.Id = ""
	createUser, err := testClient.Create(ctx, &block_user.UserRequest{
		User: user,
	})
	assert.NoError(t, err)
	// act
	newMetadata := user_mock.GetMetadata(nil)
	createUser.User.Metadata = newMetadata
	updateUser, err := testClient.UpdateMetadata(ctx, &block_user.UserRequest{
		Update: createUser.User,
		User:   createUser.User,
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
	user := user_mock.GetRandomUser(&block_user.User{
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
	})
	user.Id = ""
	createUser, err := testClient.Create(ctx, &block_user.UserRequest{
		User:          user,
		EncryptionKey: encryptionKey,
	})
	assert.NoError(t, err)
	// act
	newMetadata := user_mock.GetMetadata(nil)
	createUser.User.Metadata = newMetadata
	updateUser, err := testClient.UpdateMetadata(ctx, &block_user.UserRequest{
		Update:        createUser.User,
		User:          createUser.User,
		EncryptionKey: encryptionKey,
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
	user := user_mock.GetRandomUser(&block_user.User{
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
	})
	_, err := testClient.Create(ctx, &block_user.UserRequest{
		User: user,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.UpdateMetadata(ctx, &block_user.UserRequest{})
	// validate
	assert.Error(t, err)
}

func TestUpdateMetadataNoReq(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
	})
	user.Id = ""
	_, err := testClient.Create(ctx, &block_user.UserRequest{
		User: user,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.UpdateMetadata(ctx, nil)
	// validate
	assert.Error(t, err)
}
