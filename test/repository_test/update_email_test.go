package respository_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUpdateEmail(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	createdUser, err := testRepo.Create(ctx, user)
	initialEmail := user.Email
	initialUpdatedAt := user.UpdatedAt
	assert.Nil(t, err)
	// act
	newEmail := gofakeit.Email()
	createdUser.Email = newEmail
	updatedUser, err := testRepo.UpdateEmail(ctx, createdUser)
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.NotEmpty(t, updatedUser.Email)
	assert.NotEqual(t, initialEmail, updatedUser.Email)
	assert.True(t, updatedUser.UpdatedAt.IsValid())
	assert.NotEqual(t, initialUpdatedAt.Nanos, updatedUser.UpdatedAt.Nanos)
	// validate in database
	getUser, err := testRepo.GetById(ctx, createdUser)
	assert.Nil(t, err)
	assert.Nil(t, user_mock.CompareUsers(getUser, updatedUser))
}

func TestUpdateInvalidEmail(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	createdUser, err := testRepo.Create(ctx, user)
	assert.Nil(t, err)
	// act
	newEmail := "email@@softcorp.io"
	createdUser.Email = newEmail
	_, err = testRepo.UpdateEmail(ctx, createdUser)
	// validate
	assert.Error(t, err)
}

func TestUpdateEmailInvalidNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	createdUser, err := testRepo.Create(ctx, user)
	assert.Nil(t, err)
	// act
	createdUser.Email = gofakeit.Email()
	createdUser.Namespace = uuid.NewV4().String()
	_, err = testRepo.UpdateProfile(ctx, createdUser)
	// validate
	assert.Error(t, err)
}
