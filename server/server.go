package server

import (
	"context"
	"errors"
	"github.com/nuntiodev/nuntio-user-block/authenticator"
	"github.com/nuntiodev/nuntio-user-block/email"
	"github.com/nuntiodev/nuntio-user-block/handler"
	"github.com/nuntiodev/nuntio-user-block/interceptor"
	"github.com/nuntiodev/nuntio-user-block/repository"
	"github.com/nuntiodev/nuntio-user-block/server/grpc_server"
	"github.com/nuntiodev/nuntio-user-block/token"
	"github.com/nuntiodev/x/cryptox"
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

func getSender() (email.Sender, error) {
	// first check if postmark is present
	postmarkServerToken := os.Getenv("POSTMARK_SERVER_TOKEN")
	postmarkAccountToken := os.Getenv("POSTMARK_ACCOUNT_TOKEN")
	if postmarkServerToken != "" && postmarkAccountToken != "" {
		postmarkSender, err := email.NewPostmarkSender(postmarkServerToken, postmarkAccountToken)
		if err != nil {
			return nil, err
		}
		return postmarkSender, nil
	}
	return nil, errors.New("no email provider is available")
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
	myRepository, err := repository.New(mongoClient, myCrypto, encryptionKeys, zapLog, maxEmailVerificationAge)
	if err != nil {
		return nil, err
	}
	// it should be okay to spin up a service without email provider
	sender, err := getSender()
	if err != nil {
		zapLog.Warn("no valid sender available witn err: " + err.Error())
	}
	myEmail, err := email.New(sender)
	if err != nil {
		zapLog.Warn("email is not enabled with err: " + err.Error())
	}
	myHandler, err := handler.New(zapLog, myRepository, myCrypto, myToken, myEmail, maxEmailVerificationAge)
	if err != nil {
		return nil, err
	}
	myAuthenticator := authenticator.New()
	myInterceptor, err := interceptor.New(zapLog, myAuthenticator)
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
