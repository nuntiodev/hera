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
	userOne := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{})
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUserOne, err := users.Create(ctx, userOne, "")
	assert.Nil(t, err)
	createdUserTwo, err := users.Create(ctx, userTwo, "")
	assert.Nil(t, err)
	// act
	getUsers, err := users.GetAll(ctx, &go_block.UserFilter{}, "")
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
	userOne := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{})
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUserOne, err := users.Create(ctx, userOne, "")
	assert.Nil(t, err)
	createdUserTwo, err := users.Create(ctx, userTwo, encryptionKey)
	assert.Nil(t, err)
	// act
	getUsers, err := users.GetAll(ctx, &go_block.UserFilter{}, encryptionKey)
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
	userOne := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{})
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	_, err = users.Create(ctx, userOne, "")
	assert.Nil(t, err)
	_, err = users.Create(ctx, userTwo, encryptionKey)
	assert.Nil(t, err)
	// act
	_, err = users.GetAll(ctx, &go_block.UserFilter{}, invalidEncryptionKey)
	// validate
	assert.Error(t, err)
}

func TestGetAllWithPartialEncryptionNoDecryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUserOne, err := users.Create(ctx, userOne, "")
	assert.Nil(t, err)
	createdUserTwo, err := users.Create(ctx, userTwo, encryptionKey)
	assert.Nil(t, err)
	// act
	getUsers, err := users.GetAll(ctx, &go_block.UserFilter{}, "")
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
	userOne := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{})
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUserOne, err := users.Create(ctx, userOne, "")
	assert.Nil(t, err)
	createdUserTwo, err := users.Create(ctx, userTwo, "")
	assert.Nil(t, err)
	// act
	getUsers, err := users.GetAll(ctx, &go_block.UserFilter{
		Sort:  go_block.UserFilter_CREATED_AT,
		Order: go_block.UserFilter_DEC,
	}, "")
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
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	userOne := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{})
	createdUserOne, err := users.Create(ctx, userOne, "")
	assert.Nil(t, err)
	createdUserTwo, err := users.Create(ctx, userTwo, "")
	assert.Nil(t, err)
	// act
	getUsers, err := users.GetAll(ctx, &go_block.UserFilter{
		Sort:  go_block.UserFilter_UPDATE_AT,
		Order: go_block.UserFilter_DEC,
	}, "")
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
	userOne := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{})
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUserOne, err := users.Create(ctx, userOne, "")
	assert.Nil(t, err)
	createdUserTwo, err := users.Create(ctx, userTwo, "")
	assert.Nil(t, err)
	// act
	getUsers, err := users.GetAll(ctx, &go_block.UserFilter{
		Order: go_block.UserFilter_DEC,
	}, "")
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
	userOne := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{})
	userThree := user_mock.GetRandomUser(&go_block.User{})
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUserOne, err := users.Create(ctx, userOne, "")
	assert.Nil(t, err)
	createdUserTwo, err := users.Create(ctx, userTwo, "")
	assert.Nil(t, err)
	createdUserThree, err := users.Create(ctx, userThree, "")
	assert.Nil(t, err)
	// act
	getUsers, err := users.GetAll(ctx, &go_block.UserFilter{
		Order: go_block.UserFilter_DEC,
		Sort:  go_block.UserFilter_CREATED_AT,
	}, "")
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
	userOne := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{})
	usersOne, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUserOne, err := usersOne.Create(ctx, userOne, "")
	assert.Nil(t, err)
	usersTwo, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	_, err = usersTwo.Create(ctx, userTwo, "")
	assert.Nil(t, err)
	// act
	getUsers, err := usersOne.GetAll(ctx, &go_block.UserFilter{}, "")
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
	userOne := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{})
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUserOne, err := users.Create(ctx, userOne, "")
	assert.Nil(t, err)
	_, err = users.Create(ctx, userTwo, "")
	assert.Nil(t, err)
	// act
	getUsers, err := users.GetAll(ctx, &go_block.UserFilter{
		From: 0,
		To:   1,
	}, "")
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
	userOne := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{})
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	_, err = users.Create(ctx, userOne, "")
	assert.Nil(t, err)
	createdUserTwo, err := users.Create(ctx, userTwo, "")
	assert.Nil(t, err)
	// act
	getUsers, err := users.GetAll(ctx, &go_block.UserFilter{
		From: 1,
		To:   2,
	}, "")
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 1, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(createdUserTwo, getUsers[0]))
}
