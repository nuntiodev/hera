package respository_test

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
	user := user_mock.GetRandomUser(nil)
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, "")
	initialMetadata := user.Metadata
	initialUpdatedAt := user.UpdatedAt
	assert.Nil(t, err)
	// act
	createdUser.Metadata = user_mock.GetMetadata(nil)
	updatedUser, err := users.UpdateMetadata(ctx, createdUser, createdUser, "")
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.NotEmpty(t, updatedUser.Metadata)
	assert.NotEqual(t, initialMetadata, updatedUser.Metadata)
	assert.NotEqual(t, initialUpdatedAt.Nanos, updatedUser.UpdatedAt.Nanos)
	// validate in database
	getUser, err := users.Get(ctx, createdUser, "")
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
}

func TestUpdateMetadataWithEncryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, encryptionKey)
	initialMetadata := user.Metadata
	initialUpdatedAt := user.UpdatedAt
	assert.Nil(t, err)
	// act
	createdUser.Metadata = user_mock.GetMetadata(nil)
	updatedUser, err := users.UpdateMetadata(ctx, createdUser, createdUser, encryptionKey)
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.NotEmpty(t, updatedUser.Metadata)
	assert.NotEqual(t, initialMetadata, updatedUser.Metadata)
	assert.NotEqual(t, initialUpdatedAt.Nanos, updatedUser.UpdatedAt.Nanos)
	// validate in database
	getUser, err := users.Get(ctx, createdUser, encryptionKey)
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
}

func TestUpdateMetadataWithInvalidEncryptionKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, encryptionKey)
	assert.Nil(t, err)
	// act
	createdUser.Metadata = user_mock.GetMetadata(nil)
	_, err = users.UpdateMetadata(ctx, createdUser, createdUser, invalidEncryptionKey)
	// validate
	assert.Error(t, err)
}

func TestUpdateEncryptedMetadataWithoutKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
	})
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, encryptionKey)
	assert.Nil(t, err)
	// act
	createdUser.Metadata = user_mock.GetMetadata(nil)
	_, err = users.UpdateMetadata(ctx, createdUser, createdUser, "")
	// validate
	assert.Error(t, err)
}

func TestUpdateUnencryptedMetadataWithKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
	})
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, "")
	assert.Nil(t, err)
	// act
	createdUser.Metadata = user_mock.GetMetadata(nil)
	_, err = users.UpdateMetadata(ctx, createdUser, createdUser, encryptionKey)
	// validate
	assert.Error(t, err)
}

func TestUpdateMetadataInvalidNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, "")
	assert.Nil(t, err)
	// act
	usersTwo, err := testRepository.Users(ctx, "")
	assert.NoError(t, err)
	createdUser.Metadata = user_mock.GetMetadata(nil)
	_, err = usersTwo.UpdateMetadata(ctx, createdUser, createdUser, "")
	// validate
	assert.Error(t, err)
}
