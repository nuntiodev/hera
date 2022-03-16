package respository_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/softcorp-io/block-proto/go_block"
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
	createdUser, err := testRepo.Create(ctx, user, "")
	assert.Nil(t, err)
	// act
	createdUser.Role = gofakeit.Name()
	updatedUser, err := testRepo.UpdateSecurity(ctx, createdUser, createdUser, encryptionKey)
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.True(t, updatedUser.Encrypted)
	// validate in database
	_, err = testRepo.Get(ctx, createdUser, encryptionKey)
	assert.NoError(t, err)
	getUser, err := testRepo.Get(ctx, createdUser, encryptionKey)
	assert.NoError(t, err)
	createdUser.Encrypted = true
	assert.NoError(t, user_mock.CompareUsers(getUser, createdUser))
}

func TestUpdateEncryptedSecurity(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
		Image: gofakeit.ImageURL(10, 10),
	})
	createdUser, err := testRepo.Create(ctx, user, encryptionKey)
	assert.Nil(t, err)
	// act
	createdUser.Role = gofakeit.Name()
	updatedUser, err := testRepo.UpdateSecurity(ctx, createdUser, createdUser, encryptionKey)
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.False(t, updatedUser.Encrypted)
	// validate in database
	_, err = testRepo.Get(ctx, createdUser, "")
	assert.NoError(t, err)
}

func TestUpdateSecurityWithInvalidEncryptionKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
		Image: gofakeit.ImageURL(10, 10),
	})
	createdUser, err := testRepo.Create(ctx, user, encryptionKey)
	assert.Nil(t, err)
	// act
	createdUser.Role = gofakeit.Name()
	_, err = testRepo.UpdateSecurity(ctx, createdUser, createdUser, invalidEncryptionKey)
	// validate
	assert.Error(t, err)
}
