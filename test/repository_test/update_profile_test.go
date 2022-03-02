package respository_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestUpdateProfile(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	createdUser, err := testRepo.Create(ctx, user)
	initialName := user.Name
	initialImage := user.Image
	initialCountry := user.Country
	initialGender := user.Gender
	initialUpdatedAt := user.UpdatedAt
	initialBlocker := createdUser.Blocked
	assert.Nil(t, err)
	// act
	createdUser.Name = gofakeit.Name()
	createdUser.Image = gofakeit.ImageURL(10, 10)
	createdUser.Country = gofakeit.Country()
	createdUser.Birthdate = ts.Now()
	createdUser.Gender = block_user.Gender_MALE
	updatedUser, err := testRepo.UpdateProfile(ctx, createdUser)
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.NotEmpty(t, updatedUser.Name)
	assert.NotEqual(t, initialName, updatedUser.Name)
	assert.NotEmpty(t, updatedUser.Image)
	assert.NotEqual(t, initialImage, updatedUser.Image)
	assert.NotEmpty(t, updatedUser.Country)
	assert.NotEqual(t, initialCountry, updatedUser.Country)
	assert.NotEqual(t, initialGender, updatedUser.Gender)
	assert.Equal(t, initialBlocker, updatedUser.Blocked)
	assert.True(t, updatedUser.Birthdate.IsValid())
	assert.NotEqual(t, initialUpdatedAt.Nanos, updatedUser.UpdatedAt.Nanos)
	// validate in database
	getUser, err := testRepo.GetById(ctx, createdUser)
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
}

func TestUpdateProfileInvalidNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	createdUser, err := testRepo.Create(ctx, user)
	assert.Nil(t, err)
	// act
	createdUser.Name = gofakeit.Name()
	createdUser.Namespace = uuid.NewV4().String()
	_, err = testRepo.UpdateProfile(ctx, createdUser)
	// validate
	assert.Error(t, err)
}
