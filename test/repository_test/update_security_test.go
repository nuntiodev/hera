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

func TestUpdateUnencryptedSecurity(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
		Image: gofakeit.ImageURL(10, 10),
	})
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, "")
	assert.Nil(t, err)
	// act
	updatedUser, err := users.UpdateSecurity(ctx, createdUser, encryptionKey)
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.True(t, updatedUser.Encrypted)
	// validate in database
	_, err = users.Get(ctx, createdUser, encryptionKey)
	assert.NoError(t, err)
	getUser, err := users.Get(ctx, createdUser, encryptionKey)
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
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, encryptionKey)
	assert.Nil(t, err)
	// act
	updatedUser, err := users.UpdateSecurity(ctx, createdUser, encryptionKey)
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.False(t, updatedUser.Encrypted)
	// validate in database
	_, err = users.Get(ctx, createdUser, "")
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
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, encryptionKey)
	assert.Nil(t, err)
	// act
	_, err = users.UpdateSecurity(ctx, createdUser, invalidEncryptionKey)
	// validate
	assert.Error(t, err)
}
