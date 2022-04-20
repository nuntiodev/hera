package token_repository

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/x/cryptox"
	"github.com/nuntiodev/x/mockx/mongo_mock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

var (
	mongodbTestClient *mongo.Client
)

const (
	expiresAfter = time.Second * 1
)

func getTestTokenRepository(ctx context.Context, internal bool, dbName string) (*mongodbRepository, error) {
	// create the repository
	myCrypto, err := cryptox.New()
	if err != nil {
		return nil, err
	}
	var internalKeys []string
	if internal {
		internalKeyOne, err := myCrypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		if err != nil {
			return nil, err
		}
		internalKeyTwo, err := myCrypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		if err != nil {
			return nil, err
		}
		internalKeys = []string{internalKeyOne, internalKeyTwo}
	}
	if dbName == "" {
		dbName = "nuntio-testdb"
	}
	tokenRepository, err := newMongodbTokenRepository(ctx, mongodbTestClient.Database(dbName).Collection("tokens"), myCrypto, internalKeys)
	if err != nil {
		return nil, err
	}
	return tokenRepository, nil
}

func getToken(token *go_block.Token) *go_block.Token {
	if token == nil {
		token = &go_block.Token{}
	}
	if strings.TrimSpace(token.Id) == "" {
		token.Id = uuid.NewV4().String()
	}
	if strings.TrimSpace(token.UserId) == "" {
		token.UserId = uuid.NewV4().String()
	}
	if strings.TrimSpace(token.DeviceInfo) == "" {
		token.DeviceInfo = gofakeit.Phone()
	}
	if token.ExpiresAt == nil || token.ExpiresAt.IsValid() == false {
		token.ExpiresAt = ts.New(time.Now().Add(expiresAfter))
	}
	token.Type = go_block.TokenType_TOKEN_TYPE_ACCESS
	return token
}

func TestMain(m *testing.M) {
	// before test
	containerName := "token-repo-test"
	zapLog, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	myMongoDBClient, cleanup, err := mongo_mock.NewDatabaseMock(context.Background(), zapLog, containerName)
	defer cleanup()
	if err != nil {
		zapLog.Fatal(err.Error())
	}
	mongodbTestClient = myMongoDBClient
	code := m.Run()
	// after test
	os.Exit(code)
}

func TestCreateTokenRepository(t *testing.T) {
	obj, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	assert.NotNil(t, obj)
}
