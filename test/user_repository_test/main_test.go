package user_repository_test

import (
	"context"
	"fmt"
	"github.com/softcorp-io/block-user-service/repository/user_repository"
	"github.com/softcorp-io/block-user-service/test/mocks/repository_mock"
	"go.uber.org/zap"
	"os"
	"testing"
)

var testRepo user_repository.UserRepository

func TestMain(m *testing.M) {
	// before test
	zapLog, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	containerName := "mongodb-user-repo-test"
	repository, pool, container, err := repository_mock.NewRepositoryMock(context.Background(), zapLog, containerName)
	// ensure container is cleaned up.
	if err != nil {
		zapLog.Fatal(err.Error())
	}
	testRepo = repository.UserRepository
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
