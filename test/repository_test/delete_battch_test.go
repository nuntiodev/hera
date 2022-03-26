package respository_test

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDeleteBatch(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	userOne := user_mock.GetRandomUser(&go_block.User{})
	userTwo := user_mock.GetRandomUser(&go_block.User{})
	users, err := testRepository.Users(ctx, uuid.NewV4().String(), "")
	assert.NoError(t, err)
	_, err = users.Create(ctx, userOne)
	assert.NoError(t, err)
	_, err = users.Create(ctx, userTwo)
	assert.NoError(t, err)
	// act
	err = users.DeleteBatch(ctx, []*go_block.User{userOne, userTwo})
	assert.NoError(t, err)
	// validate
	getUsersDeletedNamespace, err := users.GetAll(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(getUsersDeletedNamespace))
}
