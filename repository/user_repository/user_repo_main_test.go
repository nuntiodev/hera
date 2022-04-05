package user_repository

import (
	"context"
	"fmt"
	"github.com/softcorp-io/x/cryptox"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"os"
	"testing"
)

type CustomMetadata struct {
	ClassYear int    `json:"class_year"`
	Name      string `json:"name"`
}

var (
	mongoTestClient *mongo.Client
)

func getTestUserRepository(ctx context.Context, internal, external bool, dbName string) (*mongoRepository, error) {
	// create the repository
	myCrypto, err := cryptox.New()
	if err != nil {
		return nil, err
	}
	var internalKeys []string
	externalKey := ""
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
	if external {
		externalKey, err = myCrypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		if err != nil {
			return nil, err
		}
	}
	if dbName == "" {
		dbName = "softcorp-users-testdb"
	}
	userRepository, err := newMongoUserRepository(ctx, mongoTestClient.Database(dbName).Collection("users"), myCrypto, internalKeys, externalKey, true)
	if err != nil {
		return nil, err
	}
	return userRepository, nil
}

func TestMain(m *testing.M) {
	// before test
	containerName := "user-repo-test"
	zapLog, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	myMongoClient, pool, container, err := NewDatabaseMock(context.Background(), zapLog, containerName)
	defer func() {
		if pool != nil {
			if err := pool.Purge(container); err != nil {
				zapLog.Error(fmt.Sprintf("failed to purge pool with err: %s", err))
			}
			if err := pool.RemoveContainerByName(containerName); err != nil {
				zapLog.Error(fmt.Sprintf("failed to remove Docker container with err: %s", err))
			}
		}
	}()
	if err != nil {
		zapLog.Fatal(err.Error())
	}
	mongoTestClient = myMongoClient
	code := m.Run()
	// after test
	if err := pool.Purge(container); err != nil {
		zapLog.Error(fmt.Sprintf("failed to purge pool with err: %s", err))
	}
	os.Exit(code)
}
