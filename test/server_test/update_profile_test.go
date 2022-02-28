package server_test

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

func TestUpdateProfile(t *testing.T) {
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
	user.Id = ""
	createUser, err := testClient.Create(ctx, &block_user.UserRequest{
		User: user,
	})
	assert.NoError(t, err)
	// act
	newName := gofakeit.Name()
	newBirthdate := ts.Now()
	newCountry := gofakeit.Country()
	newImage := gofakeit.ImageURL(10, 10)
	createUser.User.Name = newName
	createUser.User.Birthdate = newBirthdate
	createUser.User.Country = newCountry
	createUser.User.Image = newImage
	updateUser, err := testClient.UpdateProfile(ctx, &block_user.UserRequest{
		Update: createUser.User,
	})
	assert.NoError(t, err)
	// validate
	assert.Equal(t, newName, updateUser.User.Name)
	assert.Equal(t, newBirthdate.String(), updateUser.User.Birthdate.String())
	assert.Equal(t, newImage, updateUser.User.Image)
	assert.Equal(t, newCountry, updateUser.User.Country)
}

func TestUpdateProfileNoUser(t *testing.T) {
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
	_, err := testClient.Create(ctx, &block_user.UserRequest{
		User: user,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.UpdateProfile(ctx, &block_user.UserRequest{})
	// validate
	assert.Error(t, err)
}

func TestUpdateProfileNoReq(t *testing.T) {
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
	user.Id = ""
	_, err := testClient.Create(ctx, &block_user.UserRequest{
		User: user,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.UpdateProfile(ctx, nil)
	// validate
	assert.Error(t, err)
}
