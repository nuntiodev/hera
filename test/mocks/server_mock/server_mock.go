package server_mock

import (
	"context"
	"errors"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/runner"
	"github.com/softcorp-io/block-user-service/test/mocks/repository_mock"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"time"
)

type ServerTest struct {
	pool          *dockertest.Pool
	resource      *dockertest.Resource
	conn          *grpc.ClientConn
	containerName string
	Client        go_block.UserServiceClient
}

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

func (st *ServerTest) Purge() error {
	if st.pool != nil {
		if err := st.pool.Purge(st.resource); err != nil {
			return err
		}
		if err := st.pool.RemoveContainerByName(st.containerName); err != nil {
			return err
		}
	}
	if st.conn != nil {
		if err := st.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

func NewServerMock(ctx context.Context, zapLog *zap.Logger, containerName string, port int) (*ServerTest, error) {
	_, pool, container, err := repository_mock.NewRepositoryMock(context.Background(), zapLog, containerName)
	if err != nil {
		return nil, err
	}
	os.Setenv("GRPC_PORT", fmt.Sprintf("%d", port))
	serverTest := &ServerTest{
		pool:          pool,
		resource:      container,
		containerName: containerName,
	}
	go func() {
		defer serverTest.Purge()
		if err := runner.Run(context.Background(), zapLog); err != nil {
			zapLog.Fatal(err.Error())
		}
	}()
	userServiceConn, err := getClientConn(zapLog, port)
	if err != nil {
		serverTest.Purge()
		return nil, err
	}
	testClient := go_block.NewUserServiceClient(userServiceConn)
	serverTest.Client = testClient
	serverTest.conn = userServiceConn
	return serverTest, nil
}
