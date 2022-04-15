package user_repository

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/io-nuntio/block-proto/go_block"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateIEEncrypted(t *testing.T) {
	// setup user client
	userRepository, err := getTestUserRepository(context.Background(), true, true, "")
	assert.NoError(t, err)
	// create some metadata
	metadata, err := json.Marshal(&CustomMetadata{
		Name:      gofakeit.Name(),
		ClassYear: 3,
	})
	assert.NoError(t, err)
	password := gofakeit.Password(true, true, true, true, true, 20)
	user := &go_block.User{
		OptionalId: uuid.NewV4().String(),
		Email:      gofakeit.Email(),
		Password:   password,
		Image:      gofakeit.ImageURL(10, 10),
		Metadata:   string(metadata),
	}
	// act
	createdUser, err := userRepository.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	// assert new fields are present
	assert.NotEmpty(t, createdUser.Id)
	assert.True(t, createdUser.InternalEncrypted)
	assert.True(t, createdUser.ExternalEncrypted)
	assert.Equal(t, createdUser.ExternalEncryptionLevel, int32(1))
	assert.Equal(t, createdUser.InternalEncryptionLevel, int32(2))
	// assert that old fields are the same
	assert.Equal(t, createdUser.Email, user.Email)
	assert.NotEqual(t, password, user.Password)
	assert.Equal(t, createdUser.Image, user.Image)
	assert.Equal(t, createdUser.Metadata, user.Metadata)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte(password)))
}

func TestCreateIEncrypted(t *testing.T) {
	// setup user client
	userRepository, err := getTestUserRepository(context.Background(), true, false, "")
	assert.NoError(t, err)
	// create some metadata
	metadata, err := json.Marshal(&CustomMetadata{
		Name:      gofakeit.Name(),
		ClassYear: 3,
	})
	assert.NoError(t, err)
	password := gofakeit.Password(true, true, true, true, true, 20)
	user := &go_block.User{
		OptionalId: uuid.NewV4().String(),
		Email:      gofakeit.Email(),
		Password:   password,
		Image:      gofakeit.ImageURL(10, 10),
		Metadata:   string(metadata),
	}
	createdUser, err := userRepository.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	// assert new fields are present
	assert.NotEmpty(t, createdUser.Id)
	assert.True(t, createdUser.InternalEncrypted)
	assert.False(t, createdUser.ExternalEncrypted)
	assert.Equal(t, createdUser.ExternalEncryptionLevel, int32(0))
	assert.Equal(t, createdUser.InternalEncryptionLevel, int32(2))
	// assert that old fields are the same
	assert.Equal(t, createdUser.Email, user.Email)
	assert.NotEqual(t, password, user.Password)
	assert.Equal(t, createdUser.Image, user.Image)
	assert.Equal(t, createdUser.Metadata, user.Metadata)
}

func TestCreateEEncrypted(t *testing.T) {
	// setup user client
	userRepository, err := getTestUserRepository(context.Background(), false, true, "")
	assert.NoError(t, err)
	// create some metadata
	metadata, err := json.Marshal(&CustomMetadata{
		Name:      gofakeit.Name(),
		ClassYear: 3,
	})
	assert.NoError(t, err)
	password := gofakeit.Password(true, true, true, true, true, 20)
	user := &go_block.User{
		OptionalId: uuid.NewV4().String(),
		Email:      gofakeit.Email(),
		Password:   password,
		Image:      gofakeit.ImageURL(10, 10),
		Metadata:   string(metadata),
	}
	// act
	createdUser, err := userRepository.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	// assert new fields are present
	assert.NotEmpty(t, createdUser.Id)
	assert.False(t, createdUser.InternalEncrypted)
	assert.True(t, createdUser.ExternalEncrypted)
	assert.Equal(t, createdUser.ExternalEncryptionLevel, int32(1))
	assert.Equal(t, createdUser.InternalEncryptionLevel, int32(0))
	// assert that old fields are the same
	assert.Equal(t, createdUser.Email, user.Email)
	assert.NotEqual(t, password, user.Password)
	assert.Equal(t, createdUser.Image, user.Image)
	assert.Equal(t, createdUser.Metadata, user.Metadata)
}

func TestCreateNoEncryption(t *testing.T) {
	// setup user client
	userRepository, err := getTestUserRepository(context.Background(), false, false, "")
	assert.NoError(t, err)
	// create some metadata
	metadata, err := json.Marshal(&CustomMetadata{
		Name:      gofakeit.Name(),
		ClassYear: 3,
	})
	assert.NoError(t, err)
	password := gofakeit.Password(true, true, true, true, true, 20)
	user := &go_block.User{
		OptionalId: uuid.NewV4().String(),
		Email:      gofakeit.Email(),
		Password:   password,
		Image:      gofakeit.ImageURL(10, 10),
		Metadata:   string(metadata),
	}
	// act
	createdUser, err := userRepository.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	// assert new fields are present
	assert.NotEmpty(t, createdUser.Id)
	assert.False(t, createdUser.InternalEncrypted)
	assert.False(t, createdUser.ExternalEncrypted)
	assert.Equal(t, createdUser.ExternalEncryptionLevel, int32(0))
	assert.Equal(t, createdUser.InternalEncryptionLevel, int32(0))
	// assert that old fields are the same
	assert.Equal(t, createdUser.Email, user.Email)
	assert.NotEqual(t, password, user.Password)
	assert.Equal(t, createdUser.Image, user.Image)
	assert.Equal(t, createdUser.Metadata, user.Metadata)
}

func TestCreateInvalidPassword(t *testing.T) {
	// setup user client
	userRepository, err := getTestUserRepository(context.Background(), false, true, "")
	assert.NoError(t, err)
	// create some metadata
	metadata, err := json.Marshal(&CustomMetadata{
		Name:      gofakeit.Name(),
		ClassYear: 3,
	})
	assert.NoError(t, err)
	user := &go_block.User{
		OptionalId: uuid.NewV4().String(),
		Email:      gofakeit.Email(),
		Password:   "Test1234",
		Image:      gofakeit.ImageURL(10, 10),
		Metadata:   string(metadata),
	}
	// act
	createdUser, err := userRepository.Create(context.Background(), user)
	assert.Error(t, err)
	assert.Nil(t, createdUser)
}

func TestCreateInvalidEmail(t *testing.T) {
	// setup user client
	userRepository, err := getTestUserRepository(context.Background(), false, true, "")
	assert.NoError(t, err)
	// create some metadata
	metadata, err := json.Marshal(&CustomMetadata{
		Name:      gofakeit.Name(),
		ClassYear: 3,
	})
	assert.NoError(t, err)
	user := &go_block.User{
		OptionalId: uuid.NewV4().String(),
		Email:      "info@@softcorp.io",
		Password:   gofakeit.Password(true, true, true, true, true, 20),
		Image:      gofakeit.ImageURL(10, 10),
		Metadata:   string(metadata),
	}
	// act
	createdUser, err := userRepository.Create(context.Background(), user)
	assert.Error(t, err)
	assert.Nil(t, createdUser)
}

func TestCreateInvalidMetadata(t *testing.T) {
	// setup user client
	userRepository, err := getTestUserRepository(context.Background(), false, true, "")
	assert.NoError(t, err)
	user := &go_block.User{
		OptionalId: uuid.NewV4().String(),
		Email:      "info@softcorp.io",
		Password:   gofakeit.Password(true, true, true, true, true, 20),
		Image:      gofakeit.ImageURL(10, 10),
		Metadata:   "invalid metadata",
	}
	// act
	createdUser, err := userRepository.Create(context.Background(), user)
	assert.Error(t, err)
	assert.Nil(t, createdUser)
}

func TestCreateNilUser(t *testing.T) {
	// setup user client
	userRepository, err := getTestUserRepository(context.Background(), false, true, "")
	assert.NoError(t, err)
	createdUser, err := userRepository.Create(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, createdUser)
}
