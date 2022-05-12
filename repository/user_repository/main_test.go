package user_repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/x/cryptox"
	"github.com/nuntiodev/x/mockx/mongo_mock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type CustomMetadata struct {
	ClassYear int    `json:"class_year"`
	Name      string `json:"name"`
}

var (
	mongoTestClient *mongo.Client
)

func getTestUserRepository(ctx context.Context, internal, external bool, dbName string) (*mongodbRepository, error) {
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
		dbName = "nuntio-testdb"
	}
	userRepository, err := newMongodbUserRepository(ctx, mongoTestClient.Database(dbName).Collection("users"), myCrypto, internalKeys, externalKey, true, time.Minute*5)
	if err != nil {
		return nil, err
	}
	return userRepository, nil
}

func compareUsers(one, two *go_block.User, validateLevel bool) error {
	if one == nil {
		return errors.New("one is nil")
	} else if two == nil {
		return errors.New("two is nil")
	} else if one.Email != two.Email {
		return errors.New("emails are different")
	} else if one.Id != two.Id {
		return errors.New("ids are different")
	} else if one.Image != two.Image {
		return errors.New("images are different")
	} else if one.Metadata != two.Metadata {
		return errors.New("metadata is different")
	} else if one.OptionalId != two.OptionalId {
		return errors.New("optionalIds are different")
	} else if one.Password != two.Password {
		return errors.New("password are different")
	} else if validateLevel && one.ExternalEncryptionLevel != two.ExternalEncryptionLevel {
		return fmt.Errorf("external encryption levels are different: %d/%d", one.ExternalEncryptionLevel, two.ExternalEncryptionLevel)
	} else if validateLevel && one.InternalEncryptionLevel != two.InternalEncryptionLevel {
		return errors.New("internal encryption levels are different")
	}
	return nil
}

func TestMain(m *testing.M) {
	// before test
	containerName := "user-repo-test"
	zapLog, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	myMongodbClient, cleanup, err := mongo_mock.NewDatabaseMock(context.Background(), zapLog, containerName)
	defer cleanup()
	if err != nil {
		zapLog.Fatal(err.Error())
	}
	mongoTestClient = myMongodbClient
	code := m.Run()
	// after test
	os.Exit(code)
}

func TestCreateUserRepository(t *testing.T) {
	obj, err := getTestUserRepository(context.Background(), true, true, "")
	assert.NoError(t, err)
	assert.NotNil(t, obj)
}
