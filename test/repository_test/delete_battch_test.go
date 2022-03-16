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
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
	})
	userThree := user_mock.GetRandomUser(&go_block.User{
		Namespace: uuid.NewV4().String(),
	})
	_, err := testRepo.Create(ctx, userOne, "")
	assert.NoError(t, err)
	_, err = testRepo.Create(ctx, userTwo, "")
	assert.NoError(t, err)
	_, err = testRepo.Create(ctx, userThree, "")
	assert.NoError(t, err)
	// act
	err = testRepo.DeleteBatch(ctx, []*go_block.User{userOne, userTwo, userThree}, namespace)
	assert.NoError(t, err)
	// validate
	getUsersDeletedNamespace, err := testRepo.GetAll(ctx, nil, namespace, "")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(getUsersDeletedNamespace))
	getUsersAliveNamespace, err := testRepo.GetAll(ctx, nil, userThree.Namespace, "")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(getUsersAliveNamespace))
	assert.NoError(t, user_mock.CompareUsers(getUsersAliveNamespace[0], userThree))
}
