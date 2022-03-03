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

func TestSecurityProfile(t *testing.T) {
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
	newRole := gofakeit.MacAddress()
	isBlocked := true
	isVerified := true
	createUser.User.Blocked = isBlocked
	createUser.User.Verified = isVerified
	createUser.User.Role = newRole
	updateUser, err := testClient.UpdateSecurity(ctx, &block_user.UserRequest{
		Update: createUser.User,
		User:   createUser.User,
	})
	assert.NoError(t, err)
	// validate
	assert.Equal(t, newRole, updateUser.User.Role)
	assert.Equal(t, isBlocked, updateUser.User.Blocked)
	assert.Equal(t, isVerified, updateUser.User.Verified)
}

func TestUpdateSecurityNoUser(t *testing.T) {
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
	_, err = testClient.UpdateSecurity(ctx, &block_user.UserRequest{})
	// validate
	assert.Error(t, err)
}

func TestUpdateSecurityNoReq(t *testing.T) {
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
	_, err = testClient.UpdateSecurity(ctx, nil)
	// validate
	assert.Error(t, err)
}
