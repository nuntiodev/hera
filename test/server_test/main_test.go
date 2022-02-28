package server_test

import (
	"context"
	"fmt"
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
	client, pool, container, clientConn, err := server_mock.NewServerMock(context.Background(), zapLog, containerName, 9001)
	if err != nil {
		if pool != nil {
			if err := pool.Purge(container); err != nil {
				zapLog.Error(fmt.Sprintf("failed to purge pool with err: %s", err))
			}
			if err := pool.RemoveContainerByName(containerName); err != nil {
				zapLog.Error(fmt.Sprintf("failed to remove Docker container with err: %s", err))
			}
		}
		zapLog.Fatal(err.Error())
	}
	testClient = client
	code := m.Run()
	// after test
	if err := pool.Purge(container); err != nil {
		zapLog.Error(fmt.Sprintf("failed to purge pool with err: %s", err))
	}
	if err := pool.RemoveContainerByName(containerName); err != nil {
		zapLog.Error(fmt.Sprintf("failed to remove Docker container with err: %s", err))
	}
	clientConn.Close()
	os.Exit(code)
}
