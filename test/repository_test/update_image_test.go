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

func TestUpdateImage(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
		Image: gofakeit.ImageURL(10, 10),
	})
	users, err := testRepository.Users(ctx, uuid.NewV4().String(), "")
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user)
	initialImage := user.Image
	initialUpdatedAt := user.UpdatedAt
	assert.Nil(t, err)
	// act
	createdUser.Image = gofakeit.ImageURL(20, 10)
	updatedUser, err := users.UpdateImage(ctx, createdUser, createdUser)
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.NotEmpty(t, updatedUser.Image)
	assert.NotEqual(t, initialImage, updatedUser.Image)
	assert.NotEqual(t, initialUpdatedAt.Nanos, updatedUser.UpdatedAt.Nanos)
	// validate in database
	getUser, err := users.Get(ctx, createdUser)
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
}

func TestUpdateImageWithEncryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
		Image: gofakeit.ImageURL(10, 10),
	})
	users, err := testRepository.Users(ctx, uuid.NewV4().String(), encryptionKey)
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user)
	initialImage := user.Image
	initialUpdatedAt := user.UpdatedAt
	assert.Nil(t, err)
	// act
	createdUser.Image = gofakeit.ImageURL(20, 10)
	updatedUser, err := users.UpdateImage(ctx, createdUser, createdUser)
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.NotEmpty(t, updatedUser.Image)
	assert.NotEqual(t, initialImage, updatedUser.Image)
	assert.NotEqual(t, initialUpdatedAt.Nanos, updatedUser.UpdatedAt.Nanos)
	// validate in database
	getUser, err := users.Get(ctx, createdUser)
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
}

func TestUpdateImageWithInvalidEncryptionKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
		Image: gofakeit.ImageURL(10, 10),
	})
	usersOne, err := testRepository.Users(ctx, uuid.NewV4().String(), encryptionKey)
	assert.NoError(t, err)
	createdUser, err := usersOne.Create(ctx, user)
	assert.Nil(t, err)
	// act
	createdUser.Image = gofakeit.ImageURL(20, 10)
	usersTwo, err := testRepository.Users(ctx, uuid.NewV4().String(), invalidEncryptionKey)
	assert.NoError(t, err)
	_, err = usersTwo.UpdateImage(ctx, createdUser, createdUser)
	// validate
	assert.Error(t, err)
}

func TestUpdateEncryptedImageWithoutKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
		Image: gofakeit.ImageURL(10, 10),
	})
	usersOne, err := testRepository.Users(ctx, uuid.NewV4().String(), encryptionKey)
	assert.NoError(t, err)
	createdUser, err := usersOne.Create(ctx, user)
	assert.Nil(t, err)
	// act
	createdUser.Image = gofakeit.ImageURL(20, 10)
	usersTwo, err := testRepository.Users(ctx, uuid.NewV4().String(), "")
	assert.NoError(t, err)
	_, err = usersTwo.UpdateImage(ctx, createdUser, createdUser)
	// validate
	assert.Error(t, err)
}

func TestUpdateUnencryptedImageWithKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
		Image: gofakeit.ImageURL(10, 10),
	})
	usersOne, err := testRepository.Users(ctx, uuid.NewV4().String(), "")
	assert.NoError(t, err)
	createdUser, err := usersOne.Create(ctx, user)
	assert.Nil(t, err)
	// act
	createdUser.Image = gofakeit.ImageURL(20, 10)
	usersTwo, err := testRepository.Users(ctx, uuid.NewV4().String(), encryptionKey)
	assert.NoError(t, err)
	_, err = usersTwo.UpdateImage(ctx, createdUser, createdUser)
	// validate
	assert.Error(t, err)
}

func TestUpdateImageInvalidNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	users, err := testRepository.Users(ctx, uuid.NewV4().String(), "")
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user)
	assert.Nil(t, err)
	// act
	createdUser.Image = gofakeit.ImageURL(20, 10)
	usersTwo, err := testRepository.Users(ctx, uuid.NewV4().String(), "")
	assert.NoError(t, err)
	_, err = usersTwo.UpdateImage(ctx, createdUser, createdUser)
	// validate
	assert.Error(t, err)
}
