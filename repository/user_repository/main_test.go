package user_repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"testing"
	"time"

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

func compareUsers(one, two *models.User) error {
	if one == nil {
		return errors.New("one is nil")
	} else if two == nil {
		return errors.New("two is nil")
	} else if one.Id != two.Id {
		return fmt.Errorf("different ids: %s and %s", one.Id, two.Id)
	} else if one.Username.Body != two.Username.Body {
		return fmt.Errorf("different usernames: %s and %s", one.Username.Body, two.Username.Body)
	} else if one.Email.Body != two.Email.Body {
		return fmt.Errorf("different emails: %s and %s", one.Email.Body, two.Email.Body)
	} else if one.Image.Body != two.Image.Body {
		return fmt.Errorf("different images: %s and %s", one.Image.Body, two.Image.Body)
	} else if one.Metadata.Body != two.Metadata.Body {
		return fmt.Errorf("different metadata: %s and %s", one.Metadata.Body, two.Metadata.Body)
	} else if one.FirstName.Body != two.FirstName.Body {
		return fmt.Errorf("different firstnames: %s and %s", one.FirstName.Body, two.FirstName.Body)
	} else if one.LastName.Body != two.LastName.Body {
		return fmt.Errorf("different lastnames: %s and %s", one.LastName.Body, two.LastName.Body)
	} else if one.PhoneNumber.Body != two.PhoneNumber.Body {
		return fmt.Errorf("different phone numbers: %s and %s", one.PhoneNumber.Body, two.PhoneNumber.Body)
	}
	return nil
}

func getTestUser() go_block.User {
	metadata, _ := json.Marshal(&CustomMetadata{
		Name:      gofakeit.Name(),
		ClassYear: 3,
	})
	return go_block.User{
		Id:          gofakeit.UUID(),
		Username:    gofakeit.Username(),
		Email:       gofakeit.Email(),
		Password:    gofakeit.Password(true, true, true, true, true, 30),
		Image:       gofakeit.ImageURL(10, 10),
		Metadata:    string(metadata),
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		Birthdate:   ts.Now(),
		PhoneNumber: gofakeit.Phone(),
	}
}

func getUserRepositories() ([]*mongodbRepository, error) {
	// setup available clients
	var clients []*mongodbRepository
	ns := uuid.NewString()
	userRepositoryFullEncryption, err := getTestUserRepository(context.Background(), true, true, ns)
	if err != nil {
		return nil, err
	}
	userRepositoryInternalEncryption, err := getTestUserRepository(context.Background(), true, false, ns)
	if err != nil {
		return nil, err
	}
	userRepositoryExternalEncryption, err := getTestUserRepository(context.Background(), false, true, ns)
	if err != nil {
		return nil, err
	}
	userRepositoryNoEncryption, err := getTestUserRepository(context.Background(), false, false, ns)
	if err != nil {
		return nil, err
	}
	clients = []*mongodbRepository{userRepositoryFullEncryption, userRepositoryInternalEncryption, userRepositoryExternalEncryption, userRepositoryNoEncryption}
	return clients, nil
}

func getTestUserRepository(ctx context.Context, internal, external bool, dbName string) (*mongodbRepository, error) {
	externalKey := ""
	var internalKeys []string
	if internal {
		internalKeyOne, err := cryptox.GenerateSymmetricKey(32, cryptox.AlphaNum)
		if err != nil {
			return nil, err
		}
		internalKeyTwo, err := cryptox.GenerateSymmetricKey(32, cryptox.AlphaNum)
		if err != nil {
			return nil, err
		}
		internalKeys = []string{internalKeyOne, internalKeyTwo}
	}
	if external {
		var err error
		externalKey, err = cryptox.GenerateSymmetricKey(32, cryptox.AlphaNum)
		if err != nil {
			return nil, err
		}
	}
	if dbName == "" {
		dbName = "nuntio-testdb"
	}
	// create the repository
	myCrypto, err := cryptox.New(internalKeys, []string{externalKey})
	if err != nil {
		return nil, err
	}
	userRepository, err := newMongodbUserRepository(ctx, mongoTestClient.Database(dbName).Collection("users"), myCrypto, true, time.Minute*5)
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
