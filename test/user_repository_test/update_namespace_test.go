package user_repository_test

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUpdateNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	user.Namespace = ""
	createdUser, err := testRepo.Create(ctx, user)
	assert.Nil(t, err)
	assert.Empty(t, createdUser.Namespace)
	// act
	createdUser.Namespace = uuid.NewV4().String()
	updatedUser, err := testRepo.UpdateNamespace(ctx, createdUser)
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.NotEmpty(t, updatedUser.Namespace)
	// validate in database
	getUser, err := testRepo.GetById(ctx, createdUser)
	assert.Nil(t, err)
	assert.Nil(t, user_mock.CompareUsers(getUser, updatedUser))
}
