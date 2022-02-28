package server_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	password := user.Password
	user.Id = ""
	// act
	createUser, err := testClient.Create(ctx, &block_user.UserRequest{
		User: user,
	})
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, createUser)
	assert.NotNil(t, createUser.User)
	assert.NotEmpty(t, createUser.User.Name)
	assert.NotEmpty(t, createUser.User.Email)
	assert.NotEmpty(t, createUser.User.Id)
	assert.NotEmpty(t, createUser.User.Namespace)
	assert.NotEmpty(t, createUser.User.Image)
	assert.NotEqual(t, block_user.Gender_INVALID_GENDER, createUser.User.Gender)
	assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(createUser.User.Password), []byte(password)))
	assert.True(t, createUser.User.Birthdate.IsValid())
	assert.True(t, createUser.User.UpdatedAt.IsValid())
	assert.True(t, createUser.User.CreatedAt.IsValid())
	// validate in database
	getUser, err := testClient.Get(ctx, &block_user.UserRequest{
		User: createUser.User,
	})
	assert.Nil(t, err)
	assert.NotNil(t, getUser)
	assert.NotNil(t, getUser.User)
	assert.Nil(t, user_mock.CompareUsers(getUser.User, createUser.User))
}

func TestCreateNoUser(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	// act
	_, err := testClient.Create(ctx, &block_user.UserRequest{})
	assert.Error(t, err)
}

func TestValidateNoUReq(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	// act
	_, err := testClient.Create(ctx, nil)
	assert.Error(t, err)
}
