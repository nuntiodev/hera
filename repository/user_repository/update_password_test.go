package user_repository

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/nuntiodev/block-proto/go_block"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUpdatePasswordIEEncryptedById(t *testing.T) {
	// setup available clients
	var clients []*mongodbRepository
	userRepositoryFullEncryption, err := getTestUserRepository(context.Background(), true, true, "")
	assert.NoError(t, err)
	userRepositoryInternalEncryption, err := getTestUserRepository(context.Background(), true, false, "")
	assert.NoError(t, err)
	userRepositoryExternalEncryption, err := getTestUserRepository(context.Background(), false, true, "")
	assert.NoError(t, err)
	userRepositoryNoEncryption, err := getTestUserRepository(context.Background(), false, false, "")
	assert.NoError(t, err)
	clients = []*mongodbRepository{userRepositoryFullEncryption, userRepositoryInternalEncryption, userRepositoryExternalEncryption, userRepositoryNoEncryption}
	for _, userRepository := range clients {
		// create some metadata
		metadata, err := json.Marshal(&CustomMetadata{
			Name:      gofakeit.Name(),
			ClassYear: 3,
		})
		assert.NoError(t, err)
		password := gofakeit.Password(true, true, true, true, true, 20)
		user := &go_block.User{
			Username: uuid.NewV4().String(),
			Email:    gofakeit.Email(),
			Password: password,
			Image:    gofakeit.ImageURL(10, 10),
			Metadata: string(metadata),
		}
		createdUser, err := userRepository.Create(context.Background(), user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		// act
		newPassword := "My@Secure3NewPassword1234!"
		createdUser.Password = newPassword
		createdUser.Username = ""
		createdUser.Email = ""
		updatedUser, err := userRepository.UpdatePassword(context.Background(), createdUser, createdUser)
		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		// validate updated fields
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte(newPassword)))
		// validate change has been updated in db
		getUser, err := userRepository.Get(context.Background(), updatedUser, true)
		assert.NoError(t, err)
		assert.Equal(t, updatedUser.Password, getUser.Password)
		// assert.NoError(t, compareUsers(getUser, updatedUser, true)) todo: return a valid new state of user
	}
}

func TestUpdatePasswordIEEncryptedByEmail(t *testing.T) {
	// setup available clients
	var clients []*mongodbRepository
	userRepositoryFullEncryption, err := getTestUserRepository(context.Background(), true, true, "")
	assert.NoError(t, err)
	userRepositoryInternalEncryption, err := getTestUserRepository(context.Background(), true, false, "")
	assert.NoError(t, err)
	userRepositoryExternalEncryption, err := getTestUserRepository(context.Background(), false, true, "")
	assert.NoError(t, err)
	userRepositoryNoEncryption, err := getTestUserRepository(context.Background(), false, false, "")
	assert.NoError(t, err)
	clients = []*mongodbRepository{userRepositoryFullEncryption, userRepositoryInternalEncryption, userRepositoryExternalEncryption, userRepositoryNoEncryption}
	for _, userRepository := range clients {
		// create some metadata
		metadata, err := json.Marshal(&CustomMetadata{
			Name:      gofakeit.Name(),
			ClassYear: 3,
		})
		assert.NoError(t, err)
		password := gofakeit.Password(true, true, true, true, true, 20)
		user := &go_block.User{
			Username: uuid.NewV4().String(),
			Email:    gofakeit.Email(),
			Password: password,
			Image:    gofakeit.ImageURL(10, 10),
			Metadata: string(metadata),
		}
		createdUser, err := userRepository.Create(context.Background(), user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		// act
		newPassword := "My@Secure3NewPassword1234!"
		createdUser.Password = newPassword
		createdUser.Username = ""
		createdUser.Id = ""
		updatedUser, err := userRepository.UpdatePassword(context.Background(), createdUser, createdUser)
		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		// validate updated fields
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte(newPassword)))
	}
}

func TestUpdatePasswordIEEncryptedByUsername(t *testing.T) {
	// setup available clients
	var clients []*mongodbRepository
	userRepositoryFullEncryption, err := getTestUserRepository(context.Background(), true, true, "")
	assert.NoError(t, err)
	userRepositoryInternalEncryption, err := getTestUserRepository(context.Background(), true, false, "")
	assert.NoError(t, err)
	userRepositoryExternalEncryption, err := getTestUserRepository(context.Background(), false, true, "")
	assert.NoError(t, err)
	userRepositoryNoEncryption, err := getTestUserRepository(context.Background(), false, false, "")
	assert.NoError(t, err)
	clients = []*mongodbRepository{userRepositoryFullEncryption, userRepositoryInternalEncryption, userRepositoryExternalEncryption, userRepositoryNoEncryption}
	for _, userRepository := range clients {
		// create some metadata
		metadata, err := json.Marshal(&CustomMetadata{
			Name:      gofakeit.Name(),
			ClassYear: 3,
		})
		assert.NoError(t, err)
		password := gofakeit.Password(true, true, true, true, true, 20)
		user := &go_block.User{
			Username: uuid.NewV4().String(),
			Email:    gofakeit.Email(),
			Password: password,
			Image:    gofakeit.ImageURL(10, 10),
			Metadata: string(metadata),
		}
		createdUser, err := userRepository.Create(context.Background(), user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		// act
		newPassword := "My@Secure3NewPassword1234!"
		createdUser.Password = newPassword
		createdUser.Email = ""
		createdUser.Id = ""
		updatedUser, err := userRepository.UpdatePassword(context.Background(), createdUser, createdUser)
		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		// validate updated fields
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte(newPassword)))
	}
}

func TestUpdatePasswordWeakPassword(t *testing.T) {
	// setup available clients
	var clients []*mongodbRepository
	userRepositoryFullEncryption, err := getTestUserRepository(context.Background(), true, true, "")
	assert.NoError(t, err)
	userRepositoryInternalEncryption, err := getTestUserRepository(context.Background(), true, false, "")
	assert.NoError(t, err)
	userRepositoryExternalEncryption, err := getTestUserRepository(context.Background(), false, true, "")
	assert.NoError(t, err)
	userRepositoryNoEncryption, err := getTestUserRepository(context.Background(), false, false, "")
	assert.NoError(t, err)
	clients = []*mongodbRepository{userRepositoryFullEncryption, userRepositoryInternalEncryption, userRepositoryExternalEncryption, userRepositoryNoEncryption}
	for _, userRepository := range clients {
		// create some metadata
		metadata, err := json.Marshal(&CustomMetadata{
			Name:      gofakeit.Name(),
			ClassYear: 3,
		})
		assert.NoError(t, err)
		password := gofakeit.Password(true, true, true, true, true, 20)
		user := &go_block.User{
			Username: uuid.NewV4().String(),
			Email:    gofakeit.Email(),
			Password: password,
			Image:    gofakeit.ImageURL(10, 10),
			Metadata: string(metadata),
		}
		createdUser, err := userRepository.Create(context.Background(), user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		// act
		newPassword := "newpassword"
		createdUser.Password = newPassword
		updatedUser, err := userRepository.UpdatePassword(context.Background(), createdUser, createdUser)
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	}
}

func TestUpdatePasswordNoUpdate(t *testing.T) {
	// setup available clients
	var clients []*mongodbRepository
	userRepositoryFullEncryption, err := getTestUserRepository(context.Background(), true, true, "")
	assert.NoError(t, err)
	userRepositoryInternalEncryption, err := getTestUserRepository(context.Background(), true, false, "")
	assert.NoError(t, err)
	userRepositoryExternalEncryption, err := getTestUserRepository(context.Background(), false, true, "")
	assert.NoError(t, err)
	userRepositoryNoEncryption, err := getTestUserRepository(context.Background(), false, false, "")
	assert.NoError(t, err)
	clients = []*mongodbRepository{userRepositoryFullEncryption, userRepositoryInternalEncryption, userRepositoryExternalEncryption, userRepositoryNoEncryption}
	for _, userRepository := range clients {
		// create some metadata
		metadata, err := json.Marshal(&CustomMetadata{
			Name:      gofakeit.Name(),
			ClassYear: 3,
		})
		assert.NoError(t, err)
		password := gofakeit.Password(true, true, true, true, true, 20)
		user := &go_block.User{
			Username: uuid.NewV4().String(),
			Email:    gofakeit.Email(),
			Password: password,
			Image:    gofakeit.ImageURL(10, 10),
			Metadata: string(metadata),
		}
		createdUser, err := userRepository.Create(context.Background(), user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		// act
		createdUser.Id = ""
		createdUser.Email = ""
		createdUser.Username = ""
		updatedUser, err := userRepository.UpdatePassword(context.Background(), createdUser, createdUser)
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	}
}

func TestUpdatePasswordNilUpdate(t *testing.T) {
	// setup available clients
	var clients []*mongodbRepository
	userRepositoryFullEncryption, err := getTestUserRepository(context.Background(), true, true, "")
	assert.NoError(t, err)
	userRepositoryInternalEncryption, err := getTestUserRepository(context.Background(), true, false, "")
	assert.NoError(t, err)
	userRepositoryExternalEncryption, err := getTestUserRepository(context.Background(), false, true, "")
	assert.NoError(t, err)
	userRepositoryNoEncryption, err := getTestUserRepository(context.Background(), false, false, "")
	assert.NoError(t, err)
	clients = []*mongodbRepository{userRepositoryFullEncryption, userRepositoryInternalEncryption, userRepositoryExternalEncryption, userRepositoryNoEncryption}
	for _, userRepository := range clients {
		// create some metadata
		metadata, err := json.Marshal(&CustomMetadata{
			Name:      gofakeit.Name(),
			ClassYear: 3,
		})
		assert.NoError(t, err)
		password := gofakeit.Password(true, true, true, true, true, 20)
		user := &go_block.User{
			Username: uuid.NewV4().String(),
			Email:    gofakeit.Email(),
			Password: password,
			Image:    gofakeit.ImageURL(10, 10),
			Metadata: string(metadata),
		}
		createdUser, err := userRepository.Create(context.Background(), user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		// act
		updatedUser, err := userRepository.UpdatePassword(context.Background(), createdUser, nil)
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	}
}

func TestUpdatePasswordNilGet(t *testing.T) {
	// setup available clients
	var clients []*mongodbRepository
	userRepositoryFullEncryption, err := getTestUserRepository(context.Background(), true, true, "")
	assert.NoError(t, err)
	userRepositoryInternalEncryption, err := getTestUserRepository(context.Background(), true, false, "")
	assert.NoError(t, err)
	userRepositoryExternalEncryption, err := getTestUserRepository(context.Background(), false, true, "")
	assert.NoError(t, err)
	userRepositoryNoEncryption, err := getTestUserRepository(context.Background(), false, false, "")
	assert.NoError(t, err)
	clients = []*mongodbRepository{userRepositoryFullEncryption, userRepositoryInternalEncryption, userRepositoryExternalEncryption, userRepositoryNoEncryption}
	for _, userRepository := range clients {
		// create some metadata
		metadata, err := json.Marshal(&CustomMetadata{
			Name:      gofakeit.Name(),
			ClassYear: 3,
		})
		assert.NoError(t, err)
		password := gofakeit.Password(true, true, true, true, true, 20)
		user := &go_block.User{
			Username: uuid.NewV4().String(),
			Email:    gofakeit.Email(),
			Password: password,
			Image:    gofakeit.ImageURL(10, 10),
			Metadata: string(metadata),
		}
		createdUser, err := userRepository.Create(context.Background(), user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		// act
		updatedUser, err := userRepository.UpdatePassword(context.Background(), nil, createdUser)
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	}
}
