package server_test

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

func TestValidateCredentials(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	user.Id = ""
	namespace := uuid.NewV4().String()
	password := user.Password
	createUser, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      user,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	// act
	createUser.User.Password = password
	_, err = testClient.ValidateCredentials(ctx, &go_block.UserRequest{
		User:      createUser.User,
		Namespace: namespace,
	})
	// validate
	assert.NoError(t, err)
}

func TestValidateCredentialsWithEncryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	user.Id = ""
	namespace := uuid.NewV4().String()
	password := user.Password
	createUser, err := testClient.Create(ctx, &go_block.UserRequest{
		User:          user,
		EncryptionKey: encryptionKey,
		Namespace:     namespace,
	})
	assert.NoError(t, err)
	// act
	createUser.User.Password = password
	_, err = testClient.ValidateCredentials(ctx, &go_block.UserRequest{
		User:          createUser.User,
		EncryptionKey: encryptionKey,
		Namespace:     namespace,
	})
	// validate
	assert.NoError(t, err)
}

func TestValidateCredentialsWithout(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	user.Id = ""
	namespace := uuid.NewV4().String()
	password := user.Password
	createUser, err := testClient.Create(ctx, &go_block.UserRequest{
		User:          user,
		EncryptionKey: encryptionKey,
		Namespace:     namespace,
	})
	assert.NoError(t, err)
	// act
	createUser.User.Password = password
	_, err = testClient.ValidateCredentials(ctx, &go_block.UserRequest{
		User:      createUser.User,
		Namespace: namespace,
	})
	// validate
	assert.NoError(t, err)
}

func TestValidateCredentialsNoPassword(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	user.Password = ""
	namespace := uuid.NewV4().String()
	createUser, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      user,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.ValidateCredentials(ctx, &go_block.UserRequest{
		User:      createUser.User,
		Namespace: namespace,
	})
	// validate
	assert.Error(t, err)
}

func TestValidateCredentialsNoUser(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	namespace := uuid.NewV4().String()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	_, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      user,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.ValidateCredentials(ctx, &go_block.UserRequest{
		Namespace: namespace,
	})
	// validate
	assert.Error(t, err)
}

func TestValidateCredentialsNoReq(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
	})
	user.Id = ""
	namespace := uuid.NewV4().String()
	_, err := testClient.Create(ctx, &go_block.UserRequest{
		User:      user,
		Namespace: namespace,
	})
	assert.NoError(t, err)
	// act
	_, err = testClient.ValidateCredentials(ctx, nil)
	// validate
	assert.Error(t, err)
}
