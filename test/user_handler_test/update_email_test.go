package user_handler_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestUpdateEmail(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
	})
	_, err := testHandler.Create(ctx, &block_user.UserRequest{
		User: user,
	})
	assert.Nil(t, err)
	// act
	user.Email = gofakeit.Email()
	resp, err := testHandler.UpdateEmail(ctx, &block_user.UserRequest{
		Update: user,
	})
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, resp)
}
