package respository_test

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDeleteNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&block_user.User{
		Namespace: namespace,
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Namespace: namespace,
	})
	userThree := user_mock.GetRandomUser(&block_user.User{
		Namespace: uuid.NewV4().String(),
	})
	_, err := testRepo.Create(ctx, userOne, nil)
	assert.NoError(t, err)
	_, err = testRepo.Create(ctx, userTwo, nil)
	assert.NoError(t, err)
	_, err = testRepo.Create(ctx, userThree, nil)
	assert.NoError(t, err)
	// act
	err = testRepo.DeleteNamespace(ctx, namespace)
	assert.NoError(t, err)
	// validate
	getUsersDeleteNamespace, err := testRepo.GetAll(ctx, nil, namespace, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(getUsersDeleteNamespace))
	getUsersAliveNamespace, err := testRepo.GetAll(ctx, nil, userThree.Namespace, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(getUsersAliveNamespace))
	assert.NoError(t, user_mock.CompareUsers(getUsersAliveNamespace[0], userThree))
}
