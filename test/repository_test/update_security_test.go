package respository_test

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/softcorp-io/block-proto/go_block/block_user"
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
	user := user_mock.GetRandomUser(&block_user.User{
		Email:   gofakeit.Email(),
		Name:    gofakeit.Name(),
		Image:   gofakeit.ImageURL(10, 10),
		Country: gofakeit.Country(),
		Gender:  block_user.Gender_FEMALE,
	})
	createdUser, err := testRepo.Create(ctx, user, nil)
	assert.Nil(t, err)
	// act
	createdUser.Role = gofakeit.Name()
	createdUser.Blocked = true
	createdUser.Verified = true
	createdUser.DisablePasswordValidation = true
	updatedUser, err := testRepo.UpdateSecurity(ctx, createdUser, createdUser, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.True(t, updatedUser.Encrypted)
	assert.True(t, updatedUser.Verified)
	assert.True(t, updatedUser.Blocked)
	// validate in database
	getUser, err := testRepo.Get(ctx, createdUser, nil)
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
	fmt.Println(getUser)
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
	user := user_mock.GetRandomUser(&block_user.User{
		Email:   gofakeit.Email(),
		Name:    gofakeit.Name(),
		Image:   gofakeit.ImageURL(10, 10),
		Country: gofakeit.Country(),
		Gender:  block_user.Gender_FEMALE,
	})
	createdUser, err := testRepo.Create(ctx, user, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.Nil(t, err)
	// act
	createdUser.Role = gofakeit.Name()
	createdUser.Blocked = true
	createdUser.Verified = true
	createdUser.DisablePasswordValidation = true
	updatedUser, err := testRepo.UpdateSecurity(ctx, createdUser, createdUser, nil)
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.True(t, updatedUser.Encrypted)
	assert.True(t, updatedUser.Verified)
	assert.True(t, updatedUser.Blocked)
	// validate in database
	getUser, err := testRepo.Get(ctx, createdUser, nil)
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
}
