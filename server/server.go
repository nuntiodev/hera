package server

import (
	"context"
	"github.com/nuntiodev/hera/authenticator"
	"github.com/nuntiodev/hera/email"
	"github.com/nuntiodev/hera/handler"
	"github.com/nuntiodev/hera/interceptor"
	"github.com/nuntiodev/hera/repository"
	"github.com/nuntiodev/hera/server/grpc_server"
	"github.com/nuntiodev/hera/text"
	"github.com/nuntiodev/hera/token"
	"github.com/nuntiodev/x/pointerx"
	database "github.com/nuntiodev/x/repositoryx"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

type Server struct {
	GrpcServer *grpc_server.Server
}

var (
	encryptionKeys          []string
	maxEmailVerificationAge = time.Minute * 5
)

func initialize() error {
	encryptionKeysString, _ := os.LookupEnv("ENCRYPTION_KEYS")
	encryptionKeys = strings.Fields(encryptionKeysString)
	for i, key := range encryptionKeys {
		encryptionKeys[i] = strings.TrimSpace(key)
	}
	// MAX_EMAIL_VERIFICATION_AGE
	maxEmailVerificationAgeString, ok := os.LookupEnv("MAX_EMAIL_VERIFICATION_AGE")
	if ok && maxEmailVerificationAgeString == "" {
		t, err := time.ParseDuration(maxEmailVerificationAgeString)
		if err != nil {
			maxEmailVerificationAge = t
		}
	}
	return nil
}

func New(ctx context.Context, logger *zap.Logger) (*Server, error) {
	if err := initialize(); err != nil {
		return nil, err
	}
	myDatabase, err := database.CreateDatabase(logger)
	if err != nil {
		return nil, err
	}
	mongoClient, err := myDatabase.CreateMongoClient(ctx, pointerx.IntPtr(5))
	if err != nil {
		return nil, err
	}
	myToken, err := token.New()
	if err != nil {
		return nil, err
	}
	myRepository, err := repository.New(mongoClient, encryptionKeys, logger, maxEmailVerificationAge)
	if err != nil {
		return nil, err
	}
	// it should be okay to spin up a service without email provider
	myEmail, err := email.New()
	if err != nil {
		logger.Warn("Email is not enabled. If you want to enable the email interface, override the EmailSender interface.")
	}
	myText, err := text.New()
	if err != nil {
		logger.Warn("Text is not enabled. If you want to enable the text interface, override the TextSender interface.")
	}
	myHandler, err := handler.New(logger, myRepository, myToken, myEmail, myText, maxEmailVerificationAge)
	if err != nil {
		return nil, err
	}
	myAuthenticator := authenticator.New()
	myInterceptor, err := interceptor.New(logger, myAuthenticator)
	if err != nil {
		return nil, err
	}
	grpcServer, err := grpc_server.New(logger, myHandler, myInterceptor)
	if err != nil {
		return nil, err
	}
	return &Server{
		GrpcServer: grpcServer,
	}, nil
}
