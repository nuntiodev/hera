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

func TestSearchName(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      "Birthe Hansen",
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Name:      "Birthe Borlund",
		Namespace: namespace,
	})
	userThree := user_mock.GetRandomUser(&block_user.User{
		Name:      "Dan Sommer",
		Namespace: namespace,
	})
	_, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userTwo)
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userThree)
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.Search(ctx, "Birthe", &block_user.UserFilter{
		Namespace: namespace,
	})
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 2, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(userTwo, getUsers[0]))
	assert.Nil(t, user_mock.CompareUsers(userOne, getUsers[1]))
}

func TestSearchNameDifferentNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      "Birthe Hansen",
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Name: "Birthe Borlund",
	})
	userThree := user_mock.GetRandomUser(&block_user.User{
		Name:      "Dan Sommer",
		Namespace: namespace,
	})
	_, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userTwo)
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userThree)
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.Search(ctx, "Birthe", &block_user.UserFilter{
		Namespace: namespace,
	})
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 1, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(userOne, getUsers[0]))
}

func TestSearchEmail(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      "Carl Hansen",
		Email:     "carl@softcorp.io",
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Name:      "Birthe Borlund",
		Email:     "birthe@softcorp.io",
		Namespace: namespace,
	})
	userThree := user_mock.GetRandomUser(&block_user.User{
		Name:      "Dan Sommer",
		Email:     "dan@gmail.com",
		Namespace: namespace,
	})
	_, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userTwo)
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userThree)
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.Search(ctx, "soft", &block_user.UserFilter{
		Namespace: namespace,
	})
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 2, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(userTwo, getUsers[0]))
	assert.Nil(t, user_mock.CompareUsers(userOne, getUsers[1]))
}

func TestSearchId(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      "Carl Hansen",
		Email:     "carl@softcorp.io",
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Name:      "Birthe Borlund",
		Email:     "birthe@softcorp.io",
		Namespace: namespace,
	})
	userThree := user_mock.GetRandomUser(&block_user.User{
		Name:      "Dan Sommer",
		Email:     "dan@gmail.com",
		Namespace: namespace,
	})
	_, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userTwo)
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userThree)
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.Search(ctx, userOne.Id, &block_user.UserFilter{
		Namespace: namespace,
	})
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 1, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(userOne, getUsers[0]))
}

func TestSearchCountry(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      "Carl Hansen",
		Email:     "carl@softcorp.io",
		Country:   "Denmark",
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Name:      "Birthe Borlund",
		Email:     "birthe@softcorp.io",
		Namespace: namespace,
	})
	userThree := user_mock.GetRandomUser(&block_user.User{
		Name:      "Dan Sommer",
		Country:   "Denmark",
		Email:     "dan@gmail.com",
		Namespace: namespace,
	})
	_, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userTwo)
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userThree)
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.Search(ctx, "Denma", &block_user.UserFilter{
		Namespace: namespace,
	})
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, getUsers)
	assert.Equal(t, 2, len(getUsers))
	assert.Nil(t, user_mock.CompareUsers(userOne, getUsers[0]))
	assert.Nil(t, user_mock.CompareUsers(userThree, getUsers[1]))
}

func TestSearchNoResults(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	userOne := user_mock.GetRandomUser(&block_user.User{
		Name:      "Carl Hansen",
		Email:     "carl@softcorp.io",
		Country:   "Denmark",
		Birthdate: ts.Now(),
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Name:      "Birthe Borlund",
		Email:     "birthe@softcorp.io",
		Namespace: namespace,
	})
	userThree := user_mock.GetRandomUser(&block_user.User{
		Name:      "Dan Sommer",
		Country:   "Denmark",
		Email:     "dan@gmail.com",
		Namespace: namespace,
	})
	_, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userTwo)
	assert.Nil(t, err)
	_, err = testRepo.Create(ctx, userThree)
	assert.Nil(t, err)
	// act
	getUsers, err := testRepo.Search(ctx, "Norway", &block_user.UserFilter{
		Namespace: namespace,
	})
	assert.Nil(t, err)
	// validate
	assert.Nil(t, getUsers)
	assert.Equal(t, 0, len(getUsers))
}
