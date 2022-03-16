package server_test

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
	"time"
)

func skipStream(t *testing.T) {
	if os.Getenv("SKIP_STREAM") != "" {
		t.Skip("Skipping testing stream test... these can only be run with a clustered MongoDB")
	}
}

func TestGetCreateStreamWithEncryption(t *testing.T) {
	skipStream(t)
	// setup create stream
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	namespace := uuid.NewV4().String()
	stream, err := testClient.GetStream(ctx, &go_block.UserRequest{
		Namespace: namespace,
		StreamType: []go_block.StreamType{
			go_block.StreamType_CREATE,
		},
		EncryptionKey:    encryptionKey,
		AutoFollowStream: false,
	})
	assert.NoError(t, err)
	defer stream.CloseSend()
	user := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
		Id:        uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
		Email:     gofakeit.Email(),
	})
	// act
	go func() {
		_, err = testClient.Create(context.Background(), &go_block.UserRequest{
			User:          user,
			EncryptionKey: encryptionKey,
		})
		assert.NoError(t, err)
	}()
	// validate
	streamResp, err := stream.Recv()
	assert.NoError(t, err)
	assert.NotNil(t, streamResp)
	assert.Equal(t, go_block.StreamType_CREATE, streamResp.StreamType)
	assert.NotNil(t, streamResp.User)
	assert.Equal(t, streamResp.User.Image, user.Image)
	assert.Equal(t, streamResp.User.Email, user.Email)
	// act two
	go func() {
		user.Email = gofakeit.Email()
		_, err = testClient.UpdateEmail(context.Background(), &go_block.UserRequest{
			User:          user,
			EncryptionKey: encryptionKey,
			Update:        user,
		})
		assert.NoError(t, err)
	}()
	// validate two
	stream.Recv()
	assert.Error(t, ctx.Err())
}

func TestGetCreateStreamWithoutEncryption(t *testing.T) {
	skipStream(t)
	// setup create stream
	namespace := uuid.NewV4().String()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	stream, err := testClient.GetStream(ctx, &go_block.UserRequest{
		Namespace: namespace,
		StreamType: []go_block.StreamType{
			go_block.StreamType_CREATE,
		},
		AutoFollowStream: false,
	})
	assert.NoError(t, err)
	defer stream.CloseSend()
	user := user_mock.GetRandomUser(&go_block.User{
		Namespace: namespace,
		Id:        uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
		Email:     gofakeit.Email(),
	})
	// act one
	go func() {
		_, err = testClient.Create(context.Background(), &go_block.UserRequest{
			User:          user,
			EncryptionKey: encryptionKey,
		})
		assert.NoError(t, err)
	}()
	// validate one
	streamResp, err := stream.Recv()
	assert.NoError(t, err)
	assert.NotNil(t, streamResp)
	assert.Equal(t, go_block.StreamType_CREATE, streamResp.StreamType)
	assert.NotNil(t, streamResp.User)
	assert.NotEqual(t, streamResp.User.Image, user.Image)
	assert.NotEqual(t, streamResp.User.Email, user.Email)
	// act two
	go func() {
		_, err = testClient.Delete(context.Background(), &go_block.UserRequest{
			User: user,
		})
		assert.NoError(t, err)
	}()
	// validate two
	stream.Recv()
	assert.Error(t, ctx.Err())
}

func TestGetCreateStreamInvalidNamespace(t *testing.T) {
	skipStream(t)
	// setup create stream
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	namespace := uuid.NewV4().String()
	stream, err := testClient.GetStream(ctx, &go_block.UserRequest{
		Namespace: namespace,
		StreamType: []go_block.StreamType{
			go_block.StreamType_CREATE,
		},
		AutoFollowStream: false,
	})
	assert.NoError(t, err)
	defer stream.CloseSend()
	user := user_mock.GetRandomUser(&go_block.User{
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
		Email:     gofakeit.Email(),
	})
	// act
	go func() {
		_, err = testClient.Create(context.Background(), &go_block.UserRequest{
			User:          user,
			EncryptionKey: encryptionKey,
		})
		assert.NoError(t, err)
	}()
	// validate
	stream.Recv()
	assert.Error(t, ctx.Err())
}

func TestGetUpdateDeleteStreamWithEncryptionWithoutAutoFollow(t *testing.T) {
	skipStream(t)
	// setup 50 users
	namespace := uuid.NewV4().String()
	var userBatch []*go_block.User
	for i := 0; i < 20; i++ {
		user := user_mock.GetRandomUser(&go_block.User{
			Namespace: namespace,
			Image:     gofakeit.ImageURL(10, 10),
			Email:     gofakeit.Email(),
			Id:        uuid.NewV4().String(),
		})
		userResp, err := testClient.Create(context.Background(), &go_block.UserRequest{
			User:          user,
			EncryptionKey: encryptionKey,
		})
		assert.NoError(t, err)
		userBatch = append(userBatch, userResp.User)
	}
	// setup update/delete stream
	stream, err := testClient.GetStream(context.Background(), &go_block.UserRequest{
		Namespace: namespace,
		UserBatch: userBatch,
		StreamType: []go_block.StreamType{
			go_block.StreamType_UPDATE,
			go_block.StreamType_DELETE,
		},
		EncryptionKey:    encryptionKey,
		AutoFollowStream: false,
	})
	assert.NoError(t, err)
	defer stream.CloseSend()
	// act one - update user
	go func() {
		userBatch[0].Email = gofakeit.Email()
		_, err = testClient.UpdateEmail(context.Background(), &go_block.UserRequest{
			User:          userBatch[0],
			EncryptionKey: encryptionKey,
			Update:        userBatch[0],
		})
		assert.NoError(t, err)
	}()
	// validate one
	streamResp, err := stream.Recv()
	assert.NoError(t, err)
	assert.NotNil(t, streamResp)
	assert.Equal(t, go_block.StreamType_UPDATE, streamResp.StreamType)
	assert.NotNil(t, streamResp.User)
	assert.Equal(t, userBatch[0].Email, streamResp.User.Email)
	// act two - delete user
	go func() {
		_, err = testClient.Delete(context.Background(), &go_block.UserRequest{
			User: userBatch[0],
		})
		assert.NoError(t, err)
	}()
	// validate two
	streamResp, err = stream.Recv()
	assert.NoError(t, err)
	assert.Equal(t, go_block.StreamType_DELETE, streamResp.StreamType)
}

func TestGetStreamDeleteBatch(t *testing.T) {
	skipStream(t)
	// setup 50 users
	namespace := uuid.NewV4().String()
	var userBatch []*go_block.User
	for i := 0; i < 20; i++ {
		user := user_mock.GetRandomUser(&go_block.User{
			Namespace: namespace,
			Image:     gofakeit.ImageURL(10, 10),
			Email:     gofakeit.Email(),
			Id:        uuid.NewV4().String(),
		})
		userResp, err := testClient.Create(context.Background(), &go_block.UserRequest{
			User:          user,
			EncryptionKey: encryptionKey,
		})
		assert.NoError(t, err)
		userBatch = append(userBatch, userResp.User)
	}
	// setup update/delete stream
	stream, err := testClient.GetStream(context.Background(), &go_block.UserRequest{
		Namespace: namespace,
		UserBatch: userBatch,
		StreamType: []go_block.StreamType{
			go_block.StreamType_UPDATE,
			go_block.StreamType_DELETE,
		},
		EncryptionKey:    encryptionKey,
		AutoFollowStream: false,
	})
	assert.NoError(t, err)
	defer stream.CloseSend()
	// act one - update user
	go func() {
		userBatch[0].Email = gofakeit.Email()
		_, err = testClient.DeleteBatch(context.Background(), &go_block.UserRequest{
			UserBatch: userBatch,
			Namespace: namespace,
		})
		assert.NoError(t, err)
	}()
	// validate one
	streamResp, err := stream.Recv()
	assert.NoError(t, err)
	assert.NotNil(t, streamResp)
	assert.Equal(t, go_block.StreamType_DELETE, streamResp.StreamType)
	fmt.Println(streamResp)
}

func TestGetUpdateDeleteStreamWithEncryptionWithAutoFollow(t *testing.T) {
	skipStream(t)
	// setup 50 users
	namespace := uuid.NewV4().String()
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		func() {
			user := user_mock.GetRandomUser(&go_block.User{
				Namespace: namespace,
				Image:     gofakeit.ImageURL(10, 10),
				Email:     gofakeit.Email(),
				Id:        uuid.NewV4().String(),
			})
			_, err := testClient.Create(context.Background(), &go_block.UserRequest{
				User:          user,
				EncryptionKey: encryptionKey,
			})
			assert.NoError(t, err)
		}()
		wg.Done()
	}
	wg.Wait()
	userResp, err := testClient.GetAll(context.Background(), &go_block.UserRequest{
		Namespace: namespace,
		Filter: &go_block.UserFilter{
			From: 0,
			To:   10,
		},
	})
	assert.NoError(t, err)
	userBatch := userResp.Users
	// setup update/delete stream
	stream, err := testClient.GetStream(context.Background(), &go_block.UserRequest{
		Namespace: namespace,
		UserBatch: userBatch,
		StreamType: []go_block.StreamType{
			go_block.StreamType_UPDATE,
			go_block.StreamType_DELETE,
			go_block.StreamType_CREATE,
		},
		AutoFollowStream: true,
		EncryptionKey:    encryptionKey,
	})
	assert.NoError(t, err)
	defer stream.CloseSend()
	// act massively
	var newUserBatch []*go_block.User
	for i := 0; i < 16; i++ {
		user := user_mock.GetRandomUser(&go_block.User{
			Namespace: namespace,
			Image:     gofakeit.ImageURL(10, 10),
			Email:     gofakeit.Email(),
			Id:        uuid.NewV4().String(),
		})
		userResp, err := testClient.Create(context.Background(), &go_block.UserRequest{
			User:          user,
			EncryptionKey: encryptionKey,
		})
		assert.NoError(t, err)
		streamResp, err := stream.Recv()
		assert.NoError(t, err)
		assert.Equal(t, go_block.StreamType_CREATE, streamResp.StreamType)
		newUserBatch = append(newUserBatch, userResp.User)
	}
	// validate we can still get stream values
	go func() {
		newUserBatch[len(newUserBatch)-1].Email = gofakeit.Email()
		_, err = testClient.UpdateEmail(context.Background(), &go_block.UserRequest{
			User:          newUserBatch[len(newUserBatch)-1],
			EncryptionKey: encryptionKey,
			Update:        newUserBatch[len(newUserBatch)-1],
		})
		assert.NoError(t, err)
	}()
	streamResp, err := stream.Recv()
	assert.NoError(t, err)
	assert.Equal(t, go_block.StreamType_UPDATE, streamResp.StreamType)
	assert.Equal(t, newUserBatch[len(newUserBatch)-1].Email, streamResp.User.Email)
}
