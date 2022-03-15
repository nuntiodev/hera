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

func TestGetStream(t *testing.T) {
	//t.Skipf("Only run this test with a mongodb deployed as a Replicaset")
	// setup
	namespace := uuid.NewV4().String()
	stream, err := testClient.GetStream(context.Background(), &go_block.UserRequest{
		Namespace: namespace,
	})
	defer stream.CloseSend()
	assert.NoError(t, err)
	user := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
		Image:     gofakeit.ImageURL(10, 10),
	})
	user.Id = ""
	// act
	userOne, err := testClient.Create(context.Background(), &go_block.UserRequest{
		User: user,
	})
	assert.NoError(t, err)
	userOne.User.Email = gofakeit.Email()
	_, err = testClient.UpdateEmail(context.Background(), &go_block.UserRequest{
		Update: userOne.User,
		User:   userOne.User,
	})
	_, err = testClient.Delete(context.Background(), &go_block.UserRequest{
		User: userOne.User,
	})
	// validate
	assert.NoError(t, err)
	time.Sleep(time.Second * 5)
}
