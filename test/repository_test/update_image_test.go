package respository_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"github.com/softcorp-io/block-user-service/repository/user_repository"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUpdateImage(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Email: gofakeit.Email(),
		Image: gofakeit.ImageURL(10, 10),
	})
	createdUser, err := testRepo.Create(ctx, user, nil)
	initialImage := user.Image
	initialUpdatedAt := user.UpdatedAt
	assert.Nil(t, err)
	// act
	createdUser.Image = gofakeit.ImageURL(20, 10)
	updatedUser, err := testRepo.UpdateImage(ctx, createdUser, createdUser, nil)
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.NotEmpty(t, updatedUser.Image)
	assert.NotEqual(t, initialImage, updatedUser.Image)
	assert.NotEqual(t, initialUpdatedAt.Nanos, updatedUser.UpdatedAt.Nanos)
	// validate in database
	getUser, err := testRepo.Get(ctx, createdUser, nil)
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
}

func TestUpdateImageWithEncryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Email: gofakeit.Email(),
		Image: gofakeit.ImageURL(10, 10),
	})
	createdUser, err := testRepo.Create(ctx, user, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	initialImage := user.Image
	initialUpdatedAt := user.UpdatedAt
	assert.Nil(t, err)
	// act
	createdUser.Image = gofakeit.ImageURL(20, 10)
	updatedUser, err := testRepo.UpdateImage(ctx, createdUser, createdUser, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.NotEmpty(t, updatedUser.Image)
	assert.NotEqual(t, initialImage, updatedUser.Image)
	assert.NotEqual(t, initialUpdatedAt.Nanos, updatedUser.UpdatedAt.Nanos)
	// validate in database
	getUser, err := testRepo.Get(ctx, createdUser, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
}

func TestUpdateImageWithInvalidEncryptionKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Email: gofakeit.Email(),
		Image: gofakeit.ImageURL(10, 10),
	})
	createdUser, err := testRepo.Create(ctx, user, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.Nil(t, err)
	// act
	createdUser.Image = gofakeit.ImageURL(20, 10)
	_, err = testRepo.UpdateImage(ctx, createdUser, createdUser, &user_repository.EncryptionOptions{
		Key: invalidEncryptionKey,
	})
	// validate
	assert.Error(t, err)
}

func TestUpdateEncryptedImageWithoutKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Email: gofakeit.Email(),
		Image: gofakeit.ImageURL(10, 10),
	})
	createdUser, err := testRepo.Create(ctx, user, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.Nil(t, err)
	// act
	createdUser.Image = gofakeit.ImageURL(20, 10)
	_, err = testRepo.UpdateImage(ctx, createdUser, createdUser, &user_repository.EncryptionOptions{})
	// validate
	assert.Error(t, err)
}

func TestUpdateUnencryptedImageWithKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Email: gofakeit.Email(),
		Image: gofakeit.ImageURL(10, 10),
	})
	createdUser, err := testRepo.Create(ctx, user, nil)
	assert.Nil(t, err)
	// act
	createdUser.Image = gofakeit.ImageURL(20, 10)
	_, err = testRepo.UpdateImage(ctx, createdUser, createdUser, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	// validate
	assert.Error(t, err)
}

func TestUpdateImageInvalidNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	createdUser, err := testRepo.Create(ctx, user, nil)
	assert.Nil(t, err)
	// act
	createdUser.Image = gofakeit.ImageURL(20, 10)
	createdUser.Namespace = uuid.NewV4().String()
	_, err = testRepo.UpdateImage(ctx, createdUser, createdUser, nil)
	// validate
	assert.Error(t, err)
}
