package token_repository

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateIEncrypted(t *testing.T) {
	// setup user client
	tokenRepositoryWithEncryption, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	tokenRepositoryNoEncryption, err := getTestTokenRepository(context.Background(), false, "")
	assert.NoError(t, err)
	clients := []*mongodbRepository{tokenRepositoryWithEncryption, tokenRepositoryNoEncryption}
	for _, tokenRepository := range clients {
		device := gofakeit.Phone()
		token := &Token{
			Id:        uuid.NewV4().String(),
			UserId:    uuid.NewV4().String(),
			Device:    device,
			ExpiresAt: time.Second.Milliseconds() * 1000,
		}
		// act
		createdToken, err := tokenRepository.Create(context.Background(), token)
		assert.NoError(t, err)
		assert.NotNil(t, createdToken)
		// assert new fields are present
		assert.NotEmpty(t, token.CreatedAt.String())
		assert.NotEmpty(t, token.UsedAt.String())
		assert.Equal(t, len(tokenRepository.internalEncryptionKeys), token.InternalEncryptionLevel)
		assert.Equal(t, token.Id, createdToken.Id)
		assert.Equal(t, token.UserId, createdToken.UserId)
		assert.Equal(t, token.ExpiresAt, createdToken.ExpiresAt)
		assert.NotEmpty(t, createdToken.Device)
		if token.Encrypted {
			assert.NotEqual(t, device, createdToken.Device)
		} else {
			assert.Equal(t, device, createdToken.Device)
		}
	}
}

func TestCreateNoId(t *testing.T) {
	// setup user client
	tokenRepository, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	user := &Token{
		UserId:    uuid.NewV4().String(),
		Device:    gofakeit.Phone(),
		ExpiresAt: time.Second.Milliseconds() * 1000,
	}
	// act
	token, err := tokenRepository.Create(context.Background(), user)
	assert.Error(t, err)
	assert.Nil(t, token)
}

func TestCreateNoUserId(t *testing.T) {
	// setup user client
	tokenRepository, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	user := &Token{
		Id:        uuid.NewV4().String(),
		Device:    gofakeit.Phone(),
		ExpiresAt: time.Second.Milliseconds() * 1000,
	}
	// act
	token, err := tokenRepository.Create(context.Background(), user)
	assert.Error(t, err)
	assert.Nil(t, token)
}

func TestCreateNoUserExpiresAt(t *testing.T) {
	// setup user client
	tokenRepository, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	user := &Token{
		Id:     uuid.NewV4().String(),
		UserId: uuid.NewV4().String(),
		Device: gofakeit.Phone(),
	}
	// act
	token, err := tokenRepository.Create(context.Background(), user)
	assert.Error(t, err)
	assert.Nil(t, token)
}
