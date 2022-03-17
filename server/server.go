package server

import (
	"context"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/crypto"
	"github.com/softcorp-io/block-user-service/handler"
	"github.com/softcorp-io/block-user-service/interceptor"
	"github.com/softcorp-io/block-user-service/repository"
	"github.com/softcorp-io/block-user-service/server/grpc_server"
	"github.com/softcorp-io/softcorp_db_helper"
	"go.uber.org/zap"
)

var (
	METADATA_TYPE = go_block.MetadataType_METADATA_TYPE_JSON
)

type Server struct {
	GrpcServer *grpc_server.Server
}

func New(ctx context.Context, zapLog *zap.Logger) (*Server, error) {
	myDatabase, err := database.CreateDatabase(zapLog)
	if err != nil {
		return nil, err
	}
	mongoClient, err := myDatabase.CreateMongoClient(ctx)
	if err != nil {
		return nil, err
	}
	myCrypto, err := crypto.New()
	if err != nil {
		return nil, err
	}
	myRepository, err := repository.New(mongoClient, myCrypto, zapLog)
	if err != nil {
		return nil, err
	}
	myHandler, err := handler.New(zapLog, myRepository, myCrypto)
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
