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

func TestUpdateOptionalId(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Email:      gofakeit.Email(),
		OptionalId: uuid.NewV4().String(),
	})
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, "")
	initialOptionalId := user.OptionalId
	initialUpdatedAt := user.UpdatedAt
	assert.Nil(t, err)
	// act
	newOptionalId := uuid.NewV4().String()
	createdUser.OptionalId = newOptionalId
	updatedUser, err := users.UpdateOptionalId(ctx, createdUser, createdUser)
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.NotEmpty(t, updatedUser.Email)
	assert.NotEqual(t, initialOptionalId, updatedUser.OptionalId)
	assert.Equal(t, newOptionalId, updatedUser.OptionalId)
	assert.NotEqual(t, initialUpdatedAt.Nanos, updatedUser.UpdatedAt.Nanos)
	// validate in database
	_, err = users.Get(ctx, &go_block.User{
		OptionalId: user.Email,
	}, "")
	assert.Error(t, err)
	getUser, err := users.Get(ctx, &go_block.User{
		OptionalId: updatedUser.OptionalId,
	}, "")
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
}

func TestOptionalIdInvalidNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	usersOne, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser, err := usersOne.Create(ctx, user, "")
	assert.Nil(t, err)
	// act
	usersTwo, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser.OptionalId = uuid.NewV4().String()
	_, err = usersTwo.UpdateOptionalId(ctx, createdUser, createdUser)
	// validate
	assert.Error(t, err)
}
