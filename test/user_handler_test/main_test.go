package user_handler_test

import (
	"context"
	"fmt"
	"github.com/softcorp-io/block-user-service/handler"
	"github.com/softcorp-io/block-user-service/test/mocks/repository_mock"
	"go.uber.org/zap"
	"os"
	"testing"
)

var testHandler handler.Handler

func TestMain(m *testing.M) {
	// before test
	zapLog, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	containerName := "mongodb-user-handler-test"
	repository, pool, container, err := repository_mock.NewRepositoryMock(context.Background(), zapLog, containerName)
	// ensure container is cleaned up.
	if err != nil {
		zapLog.Fatal(err.Error())
	}
	testHandler, err = handler.New(zapLog, repository)
	if err != nil {
		zapLog.Fatal(err.Error())
	}
	code := m.Run()
	// after test
	if err := pool.Purge(container); err != nil {
		zapLog.Fatal(fmt.Sprintf("failed to purge pool with err: %s", err))
	}
	if err := pool.RemoveContainerByName(containerName); err != nil {
		zapLog.Fatal(fmt.Sprintf("failed to remove Docker container with err: %s", err))
	}
	os.Exit(code)
}
