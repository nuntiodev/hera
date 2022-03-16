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

func TestGetAll(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
	})
	createdUserOne, err := testRepo.Create(ctx, userOne, "")
	assert.Nil(t, err)
	createdUserTwo, err := testRepo.Create(ctx, userTwo, "")
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &go_block.UserFilter{}, namespace, "")
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 2, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(createdUserOne, getUsers[0]))
	assert.Nil(t, user_mock.CompareUsers(createdUserTwo, getUsers[1]))
}

func TestGetAllWithPartialEncryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
	})
	createdUserOne, err := testRepo.Create(ctx, userOne, "")
	assert.Nil(t, err)
	createdUserTwo, err := testRepo.Create(ctx, userTwo, encryptionKey)
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &go_block.UserFilter{}, namespace, encryptionKey)
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 2, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(createdUserOne, getUsers[0]))
	assert.Nil(t, user_mock.CompareUsers(createdUserTwo, getUsers[1]))
}

func TestGetAllWithInvalidEncryptionKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
	})
	_, err := testRepo.Create(ctx, userOne, "")
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userTwo, encryptionKey)
	assert.Nil(t, err)
	// act
	_, err = testRepo.GetAll(ctx, &go_block.UserFilter{}, namespace, invalidEncryptionKey)
	// validate
	assert.Error(t, err)
}

func TestGetAllWithPartialEncryptionNoDecryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
	})
	createdUserOne, err := testRepo.Create(ctx, userOne, "")
	assert.Nil(t, err)
	createdUserTwo, err := testRepo.Create(ctx, userTwo, encryptionKey)
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &go_block.UserFilter{}, namespace, "")
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 2, len(getUsers))
	assert.NoError(t, user_mock.CompareUsers(createdUserOne, getUsers[0]))
	assert.Equal(t, createdUserTwo.Id, getUsers[1].Id)
	assert.Error(t, user_mock.CompareUsers(createdUserTwo, getUsers[1]))
}

func TestGetAllDifferentSortCreatedAt(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
	})
	createdUserOne, err := testRepo.Create(ctx, userOne, "")
	assert.Nil(t, err)
	createdUserTwo, err := testRepo.Create(ctx, userTwo, "")
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &go_block.UserFilter{
		Sort:  go_block.UserFilter_CREATED_AT,
		Order: go_block.UserFilter_DEC,
	}, namespace, "")
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 2, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(createdUserOne, getUsers[1]))
	assert.Nil(t, user_mock.CompareUsers(createdUserTwo, getUsers[0]))
}

func TestGetAllDifferentSortUpdatedAt(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
	})
	createdUserOne, err := testRepo.Create(ctx, userOne, "")
	assert.Nil(t, err)
	createdUserTwo, err := testRepo.Create(ctx, userTwo, "")
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &go_block.UserFilter{
		Sort:  go_block.UserFilter_UPDATE_AT,
		Order: go_block.UserFilter_DEC,
	}, namespace, "")
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 2, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(createdUserOne, getUsers[1]))
	assert.Nil(t, user_mock.CompareUsers(createdUserTwo, getUsers[0]))
}

func TestGetAllDifferentSortBirthdate(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
	})
	createdUserOne, err := testRepo.Create(ctx, userOne, "")
	assert.Nil(t, err)
	createdUserTwo, err := testRepo.Create(ctx, userTwo, "")
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &go_block.UserFilter{
		Order: go_block.UserFilter_DEC,
	}, namespace, "")
	// validate
	assert.NoError(t, err)
	assert.NotNil(t, getUsers)
	assert.Equal(t, 2, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(createdUserOne, getUsers[1]))
	assert.Nil(t, user_mock.CompareUsers(createdUserTwo, getUsers[0]))
}

func TestGetAllDifferentSortName(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
	})
	userThree := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
	})
	createdUserOne, err := testRepo.Create(ctx, userOne, "")
	assert.Nil(t, err)
	createdUserTwo, err := testRepo.Create(ctx, userTwo, "")
	assert.Nil(t, err)
	createdUserThree, err := testRepo.Create(ctx, userThree, "")
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &go_block.UserFilter{
		Order: go_block.UserFilter_DEC,
		Sort:  go_block.UserFilter_CREATED_AT,
	}, namespace, "")
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 3, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(createdUserOne, getUsers[2]))
	assert.Nil(t, user_mock.CompareUsers(createdUserTwo, getUsers[1]))
	assert.Nil(t, user_mock.CompareUsers(createdUserThree, getUsers[0]))
}

func TestGetAllDifferentNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{})
	createdUserOne, err := testRepo.Create(ctx, userOne, "")
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userTwo, "")
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &go_block.UserFilter{}, namespace, "")
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 1, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(createdUserOne, getUsers[0]))
}

func TestGetAllSetFromTo(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
	})
	createdUserOne, err := testRepo.Create(ctx, userOne, "")
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userTwo, "")
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &go_block.UserFilter{
		From: 0,
		To:   1,
	}, namespace, "")
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 1, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(createdUserOne, getUsers[0]))
}

func TestGetAllSetFromToWithSkip(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
	})
	_, err := testRepo.Create(ctx, userOne, "")
	assert.Nil(t, err)
	createdUserTwo, err := testRepo.Create(ctx, userTwo, "")
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &go_block.UserFilter{
		From: 1,
		To:   2,
	}, namespace, "")
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 1, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(createdUserTwo, getUsers[0]))
}
