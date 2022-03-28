package server_test

import (
	"context"
	"encoding/hex"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/test/mocks/server_mock"
	"go.uber.org/zap"
	"os"
	"testing"
	"time"
)

var testClient go_block.UserServiceClient
var encryptionKey = "VmYq3t6w9z$C&F)J@McQfTjWnZr4u7x!"
var accessTokenExpiresAt = time.Second * 10
var refreshTokenExpiresAt = time.Second * 10

func TestMain(m *testing.M) {
	// before test
	zapLog, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	encryptionKey = hex.EncodeToString([]byte(encryptionKey))
	containerName := "mongodb-user-server-test"
	os.Setenv("ACCESS_TOKEN_EXPIRY", accessTokenExpiresAt.String())
	os.Setenv("REFRESH_TOKEN_EXPIRY", refreshTokenExpiresAt.String())
	serverTest, err := server_mock.NewServerMock(context.Background(), zapLog, containerName, 9001)
	if err != nil {
		zapLog.Fatal(err.Error())
	}
	testClient = serverTest.Client
	code := m.Run()
	// after test
	if err := serverTest.Purge(); err != nil {
		zapLog.Fatal(err.Error())
	}
	os.Exit(code)
}
