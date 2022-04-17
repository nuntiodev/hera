package server

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/handler"
	"github.com/nuntiodev/nuntio-user-block/interceptor"
	"github.com/nuntiodev/nuntio-user-block/repository"
	"github.com/nuntiodev/nuntio-user-block/server/grpc_server"
	"github.com/nuntiodev/nuntio-user-block/token"
	database "github.com/nuntiodev/x/repositoryx"
	"github.com/nuntiodev/x/cryptox"
	"go.uber.org/zap"
	"os"
	"strings"
)

type Server struct {
	GrpcServer *grpc_server.Server
}

var (
	encryptionKeys []string
)

func initialize() error {
	encryptionKeysString, _ := os.LookupEnv("ENCRYPTION_KEYS")
	encryptionKeys = strings.Fields(encryptionKeysString)
	for i, key := range encryptionKeys {
		encryptionKeys[i] = strings.TrimSpace(key)
	}
	return nil
}

func New(ctx context.Context, zapLog *zap.Logger) (*Server, error) {
	if err := initialize(); err != nil {
		return nil, err
	}
	myDatabase, err := database.CreateDatabase(zapLog)
	if err != nil {
		return nil, err
	}
	mongoClient, err := myDatabase.CreateMongoClient(ctx)
	if err != nil {
		return nil, err
	}
	myCrypto, err := cryptox.New()
	if err != nil {
		return nil, err
	}
	myToken, err := token.New()
	if err != nil {
		return nil, err
	}
	myRepository, err := repository.New(mongoClient, myCrypto, encryptionKeys, zapLog)
	if err != nil {
		return nil, err
	}
	myHandler, err := handler.New(zapLog, myRepository, myCrypto, myToken)
	if err != nil {
		return nil, err
	}
	myInterceptor, err := interceptor.New(zapLog)
	if err != nil {
		return nil, err
	}
	grpcServer, err := grpc_server.New(zapLog, myHandler, myInterceptor)
	if err != nil {
		return nil, err
	}
	return &Server{
		GrpcServer: grpcServer,
	}, nil
}
