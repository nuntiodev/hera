package server_mock

import (
	"context"
	"errors"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"github.com/softcorp-io/block-user-service/runner"
	"github.com/softcorp-io/block-user-service/test/mocks/repository_mock"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"time"
)

func getClientConn(zapLog *zap.Logger, port int) (*grpc.ClientConn, error) {
	retry := 5
	for i := 0; i < retry; i++ {
		time.Sleep(time.Second * 2)
		userServiceConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			zapLog.Error(err.Error())
			continue
		}
		return userServiceConn, nil
	}
	return nil, errors.New("could not create client conn")
}

func NewServerMock(ctx context.Context, zapLog *zap.Logger, containerName string, port int) (block_user.ServiceClient, *dockertest.Pool, *dockertest.Resource, *grpc.ClientConn, error) {
	zapLog, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	_, pool, container, err := repository_mock.NewRepositoryMock(context.Background(), zapLog, containerName)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	os.Setenv("GRPC_PORT", fmt.Sprintf("%d", port))
	go func() {
		if err := runner.Run(context.Background(), zapLog); err != nil {
			zapLog.Fatal(err.Error())
		}
		if err := pool.Purge(container); err != nil {
			zapLog.Error(fmt.Sprintf("failed to purge pool with err: %s", err))
		}
		if err := pool.RemoveContainerByName(containerName); err != nil {
			zapLog.Error(fmt.Sprintf("failed to remove Docker container with err: %s", err))
		}
	}()
	userServiceConn, err := getClientConn(zapLog, port)
	if err != nil {
		if err := pool.Purge(container); err != nil {
			zapLog.Error(fmt.Sprintf("failed to purge pool with err: %s", err))
		}
		if err := pool.RemoveContainerByName(containerName); err != nil {
			zapLog.Error(fmt.Sprintf("failed to remove Docker container with err: %s", err))
		}
		return nil, nil, nil, nil, err
	}
	testClient := block_user.NewServiceClient(userServiceConn)
	return testClient, pool, container, userServiceConn, nil
}
