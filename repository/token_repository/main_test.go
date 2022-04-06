package token_repository

import (
	"context"
	"github.com/softcorp-io/block-user-service/mockx/mongo_mock"
	"github.com/softcorp-io/x/cryptox"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"os"
	"testing"
)

var (
	mongodbTestClient *mongo.Client
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
		dbName = "softcorp-testdb"
	}
	tokenRepository, err := newMongodbTokenRepository(ctx, mongodbTestClient.Database(dbName).Collection("tokens"), myCrypto, internalKeys)
	if err != nil {
		return nil, err
	}
	return tokenRepository, nil
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
