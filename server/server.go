package server

import (
	"context"
	"errors"
	"github.com/softcorp-io/block-user-service/crypto"
	"github.com/softcorp-io/block-user-service/handler"
	"github.com/softcorp-io/block-user-service/interceptor"
	"github.com/softcorp-io/block-user-service/repository"
	"github.com/softcorp-io/block-user-service/server/grpc_server"
	"github.com/softcorp-io/softcorp_db_helper"
	"go.uber.org/zap"
	"os"
	"time"
)

var (
	accessTokenExpiry  = time.Minute * 30
	refreshTokenExpiry = time.Hour * 24 * 30
	jwtPublicKey       = ""
	jwtPrivateKey      = ""
)

type Server struct {
	GrpcServer *grpc_server.Server
}

func initialize() error {
	accessTokenExpiryString, ok := os.LookupEnv("ACCESS_TOKEN_EXPIRY")
	if ok {
		dur, err := time.ParseDuration(accessTokenExpiryString)
		if err == nil {
			accessTokenExpiry = dur
		}
	}
	refreshTokenExpiryString, ok := os.LookupEnv("REFRESH_TOKEN_EXPIRY")
	if ok {
		dur, err := time.ParseDuration(refreshTokenExpiryString)
		if err == nil {
			refreshTokenExpiry = dur
		}
	}
	jwtPublicKey, ok = os.LookupEnv("JWT_PUBLIC_KEY")
	if !ok || jwtPublicKey == "" {
		return errors.New("missing required JWT_PUBLIC_KEY")
	}
	jwtPrivateKey, ok = os.LookupEnv("JWT_PRIVATE_KEY")
	if !ok || jwtPrivateKey == "" {
		return errors.New("missing required JWT_PRIVATE_KEY")
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
	myCrypto, err := crypto.New([]byte(jwtPrivateKey), []byte(jwtPublicKey))
	if err != nil {
		return nil, err
	}
	myRepository, err := repository.New(mongoClient, myCrypto, zapLog)
	if err != nil {
		return nil, err
	}
	myHandler, err := handler.New(zapLog, myRepository, myCrypto, accessTokenExpiry, refreshTokenExpiry, []byte(jwtPublicKey))
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
