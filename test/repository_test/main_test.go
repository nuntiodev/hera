package respository_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/softcorp-io/block-user-service/repository/user_repository"
	"github.com/softcorp-io/block-user-service/test/mocks/repository_mock"
	"go.uber.org/zap"
	"os"
	"testing"
)

var testRepo user_repository.UserRepository
var encryptionKey = "VmYq3t6w9z$C&F)J@McQfTjWnZr4u7x!"
var invalidEncryptionKey = "kpLå3t6w9z$C&F)J@McQfTjWnZr4u7x!"

func TestMain(m *testing.M) {
	// before test
	zapLog, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	encryptionKey = hex.EncodeToString([]byte(encryptionKey))
	containerName := "mongodb-user-repo-test"
	repository, pool, container, err := repository_mock.NewRepositoryMock(context.Background(), zapLog, containerName)
	if err != nil {
		zapLog.Fatal(err.Error())
	}
	testRepo = repository.UserRepository
	code := m.Run()
	// after test
	if err := pool.Purge(container); err != nil {
		zapLog.Error(fmt.Sprintf("failed to purge pool with err: %s", err))
	}
	if err := pool.RemoveContainerByName(containerName); err != nil {
		zapLog.Error(fmt.Sprintf("failed to remove Docker container with err: %s", err))
	}
	os.Exit(code)
}
