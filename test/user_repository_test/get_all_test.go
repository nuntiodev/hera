package user_repository_test

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

func TestGetAll(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Namespace: namespace,
	})
	createdUserOne, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	createdUserTwo, err := testRepo.Create(ctx, userTwo)
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &block_user.UserFilter{
		Namespace: namespace,
	})
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 2, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(createdUserOne, getUsers[0]))
	assert.Nil(t, user_mock.CompareUsers(createdUserTwo, getUsers[1]))
}

func TestGetAllDifferentSortCreatedAt(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Namespace: namespace,
	})
	createdUserOne, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	createdUserTwo, err := testRepo.Create(ctx, userTwo)
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &block_user.UserFilter{
		Namespace: namespace,
		Sort:      block_user.UserFilter_CREATED_AT,
		Order:     block_user.UserFilter_DEC,
	})
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
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Namespace: namespace,
	})
	createdUserOne, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	createdUserTwo, err := testRepo.Create(ctx, userTwo)
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &block_user.UserFilter{
		Namespace: namespace,
		Sort:      block_user.UserFilter_UPDATE_AT,
		Order:     block_user.UserFilter_DEC,
	})
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
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Namespace: namespace,
		Birthdate: ts.Now(),
	})
	createdUserOne, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	createdUserTwo, err := testRepo.Create(ctx, userTwo)
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &block_user.UserFilter{
		Namespace: namespace,
		Sort:      block_user.UserFilter_BIRTHDATE,
		Order:     block_user.UserFilter_DEC,
	})
	assert.Nil(t, err)
	// validate
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
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      "Ole Hansen",
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Name:      "Andreas Simonsen",
		Namespace: namespace,
		Birthdate: ts.Now(),
	})
	userThree := user_mock.GetRandomUser(&block_user.User{
		Name:      "Birthe Hansen",
		Namespace: namespace,
		Birthdate: ts.Now(),
	})
	createdUserOne, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	createdUserTwo, err := testRepo.Create(ctx, userTwo)
	assert.Nil(t, err)
	createdUserThree, err := testRepo.Create(ctx, userThree)
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &block_user.UserFilter{
		Namespace: namespace,
		Sort:      block_user.UserFilter_NAME,
		Order:     block_user.UserFilter_INC,
	})
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 3, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(createdUserTwo, getUsers[0]))
	assert.Nil(t, user_mock.CompareUsers(createdUserThree, getUsers[1]))
	assert.Nil(t, user_mock.CompareUsers(createdUserOne, getUsers[2]))
}

func TestGetAllDifferentNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{})
	createdUserOne, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userTwo)
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &block_user.UserFilter{
		Namespace: namespace,
	})
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
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Namespace: namespace,
	})
	createdUserOne, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userTwo)
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &block_user.UserFilter{
		Namespace: namespace,
		From:      0,
		To:        1,
	})
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
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Namespace: namespace,
	})
	_, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	createdUserTwo, err := testRepo.Create(ctx, userTwo)
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.GetAll(ctx, &block_user.UserFilter{
		Namespace: namespace,
		From:      1,
		To:        2,
	})
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 1, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(createdUserTwo, getUsers[0]))
}
