package server_test

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"github.com/softcorp-io/block-user-service/test/mocks/server_mock"
	"go.uber.org/zap"
	"os"
	"testing"
)

var testClient block_user.ServiceClient

func TestMain(m *testing.M) {
	// before test
	zapLog, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	containerName := "mongodb-user-server-test"
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
