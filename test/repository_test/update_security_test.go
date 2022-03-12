package respository_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/repository/user_repository"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUpdateUnencryptedSecurity(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
		Image: gofakeit.ImageURL(10, 10),
	})
	createdUser, err := testRepo.Create(ctx, user, nil)
	assert.Nil(t, err)
	// act
	createdUser.Role = gofakeit.Name()
	updatedUser, err := testRepo.UpdateSecurity(ctx, createdUser, createdUser, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.True(t, updatedUser.Encrypted)
	// validate in database
	getUser, err := testRepo.Get(ctx, createdUser, nil)
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
	getUser, err = testRepo.Get(ctx, createdUser, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.NoError(t, err)
	assert.Error(t, user_mock.CompareUsers(getUser, updatedUser))
}

func TestUpdateEncryptedSecurity(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
		Image: gofakeit.ImageURL(10, 10),
	})
	createdUser, err := testRepo.Create(ctx, user, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.Nil(t, err)
	// act
	createdUser.Role = gofakeit.Name()
	updatedUser, err := testRepo.UpdateSecurity(ctx, createdUser, createdUser, nil)
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.True(t, updatedUser.Encrypted)
	// validate in database
	getUser, err := testRepo.Get(ctx, createdUser, nil)
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
}

func TestUpdateSecurityWithInvalidEncryptionKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
		Image: gofakeit.ImageURL(10, 10),
	})
	createdUser, err := testRepo.Create(ctx, user, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.Nil(t, err)
	// act
	createdUser.Role = gofakeit.Name()
	_, err = testRepo.UpdateSecurity(ctx, createdUser, createdUser, &user_repository.EncryptionOptions{
		Key: invalidEncryptionKey,
	})
	// validate
	assert.Error(t, err)
}
