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

func TestUpdateEmail(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
	})
	users, err := testRepository.Users(ctx, uuid.NewV4().String(), "")
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user)
	initialEmail := user.Email
	initialUpdatedAt := user.UpdatedAt
	assert.Nil(t, err)
	// act
	updateEmail := "info@softcorp.io"
	createdUser.Email = updateEmail
	updatedUser, err := users.UpdateEmail(ctx, createdUser, createdUser)
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.NotEmpty(t, updatedUser.Email)
	assert.NotEqual(t, initialEmail, updatedUser.Email)
	assert.Equal(t, updateEmail, updatedUser.Email)
	assert.NotEqual(t, initialUpdatedAt.Nanos, updatedUser.UpdatedAt.Nanos)
	// validate in database
	_, err = users.Get(ctx, &go_block.User{
		Email: user.Email,
	})
	assert.Error(t, err)
	getUser, err := users.Get(ctx, &go_block.User{
		Email: updatedUser.Email,
	})
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
}

func TestUpdateEmailWithEncryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
	})
	users, err := testRepository.Users(ctx, uuid.NewV4().String(), encryptionKey)
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user)
	initialEmail := user.Email
	initialUpdatedAt := user.UpdatedAt
	assert.Nil(t, err)
	// act
	updateEmail := "info@softcorp.io"
	createdUser.Email = updateEmail
	updatedUser, err := users.UpdateEmail(ctx, createdUser, createdUser)
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.NotEmpty(t, updatedUser.Email)
	assert.NotEqual(t, initialEmail, updatedUser.Email)
	assert.Equal(t, updateEmail, updatedUser.Email)
	assert.NotEqual(t, initialUpdatedAt.Nanos, updatedUser.UpdatedAt.Nanos)
	// validate in database
	_, err = users.Get(ctx, &go_block.User{
		Email: createdUser.Email,
	})
	assert.Error(t, err)
	getUser, err := users.Get(ctx, &go_block.User{
		Email: updatedUser.Email,
	})
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
}

func TestUpdateEmailWithInvalidEncryptionKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
	})
	usersOne, err := testRepository.Users(ctx, uuid.NewV4().String(), encryptionKey)
	assert.NoError(t, err)
	createdUser, err := usersOne.Create(ctx, user)
	assert.Nil(t, err)
	// act
	createdUser.Email = gofakeit.Email()
	usersTwo, err := testRepository.Users(ctx, uuid.NewV4().String(), invalidEncryptionKey)
	assert.NoError(t, err)
	_, err = usersTwo.UpdateEmail(ctx, createdUser, createdUser)
	// validate
	assert.Error(t, err)
}

func TestUpdateEncryptedEmailWithoutKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
	})
	usersOne, err := testRepository.Users(ctx, uuid.NewV4().String(), encryptionKey)
	assert.NoError(t, err)
	createdUser, err := usersOne.Create(ctx, user)
	assert.Nil(t, err)
	// act
	createdUser.Email = gofakeit.Email()
	usersTwo, err := testRepository.Users(ctx, uuid.NewV4().String(), "")
	assert.NoError(t, err)
	_, err = usersTwo.UpdateEmail(ctx, createdUser, createdUser)
	// validate
	assert.Error(t, err)
}

func TestUpdateUnencryptedEmailWithKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
	})
	usersOne, err := testRepository.Users(ctx, uuid.NewV4().String(), "")
	assert.NoError(t, err)
	createdUser, err := usersOne.Create(ctx, user)
	assert.Nil(t, err)
	// act
	createdUser.Email = gofakeit.Email()
	usersTwo, err := testRepository.Users(ctx, uuid.NewV4().String(), encryptionKey)
	assert.NoError(t, err)
	_, err = usersTwo.UpdateEmail(ctx, createdUser, createdUser)
	// validate
	assert.Error(t, err)
}

func TestEmailInvalidNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	users, err := testRepository.Users(ctx, namespace, "")
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user)
	assert.Nil(t, err)
	// act
	createdUser.Email = gofakeit.Email()
	usersTwo, err := testRepository.Users(ctx, uuid.NewV4().String(), "")
	assert.NoError(t, err)
	_, err = usersTwo.UpdateEmail(ctx, createdUser, createdUser)
	// validate
	assert.Error(t, err)
}
