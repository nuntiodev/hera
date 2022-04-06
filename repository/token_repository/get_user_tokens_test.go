package token_repository

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/x/cryptox"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetUserTokensIEncrypted(t *testing.T) {
	// setup user client
	tokenRepositoryWithEncryption, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	tokenRepositoryNoEncryption, err := getTestTokenRepository(context.Background(), false, "")
	assert.NoError(t, err)
	clients := []*mongodbRepository{tokenRepositoryWithEncryption, tokenRepositoryNoEncryption}
	expiresAfter := time.Second * 1
	for _, tokenRepository := range clients {
		userId := uuid.NewV4().String()
		device := gofakeit.Phone()
		tokenOne := &Token{
			Id:        uuid.NewV4().String(),
			UserId:    userId,
			Device:    device,
			ExpiresAt: time.Now().UTC().Add(expiresAfter),
		}
		tokenTwo := &Token{
			Id:        uuid.NewV4().String(),
			UserId:    userId,
			Device:    device,
			ExpiresAt: time.Now().UTC().Add(expiresAfter),
		}
		tokenThree := &Token{
			Id:        uuid.NewV4().String(),
			UserId:    uuid.NewV4().String(),
			Device:    device,
			ExpiresAt: time.Now().UTC().Add(expiresAfter),
		}
		// create tokens
		createdTokenOne, err := tokenRepository.Create(context.Background(), tokenOne)
		assert.NoError(t, err)
		assert.NotNil(t, createdTokenOne)
		createdTokenTwo, err := tokenRepository.Create(context.Background(), tokenTwo)
		assert.NoError(t, err)
		assert.NotNil(t, createdTokenTwo)
		_, err = tokenRepository.Create(context.Background(), tokenThree)
		assert.NoError(t, err)
		assert.NotNil(t, createdTokenTwo)
		// insert new key
		newKey, err := tokenRepository.crypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		tokenRepository.internalEncryptionKeys = append(tokenRepository.internalEncryptionKeys, newKey)
		// act
		userTokens, err := tokenRepository.GetUserTokens(context.Background(), &Token{
			UserId: userId,
		})
		// validate response
		assert.NoError(t, err)
		assert.NotNil(t, userTokens)
		assert.Equal(t, 2, len(userTokens))
	}
}

func TestCannotGetUserTokensExpired(t *testing.T) {
	// setup user client
	tokenRepositoryWithEncryption, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	tokenRepositoryNoEncryption, err := getTestTokenRepository(context.Background(), false, "")
	assert.NoError(t, err)
	clients := []*mongodbRepository{tokenRepositoryWithEncryption, tokenRepositoryNoEncryption}
	expiresAfter := time.Second * 3
	for _, tokenRepository := range clients {
		userId := uuid.NewV4().String()
		device := gofakeit.Phone()
		tokenOne := &Token{
			Id:        uuid.NewV4().String(),
			UserId:    userId,
			Device:    device,
			ExpiresAt: time.Now().Add(expiresAfter),
		}
		tokenTwo := &Token{
			Id:        uuid.NewV4().String(),
			UserId:    userId,
			Device:    device,
			ExpiresAt: time.Now().Add(expiresAfter),
		}
		tokenThree := &Token{
			Id:        uuid.NewV4().String(),
			UserId:    uuid.NewV4().String(),
			Device:    device,
			ExpiresAt: time.Now().Add(expiresAfter),
		}
		// create tokens
		createdTokenOne, err := tokenRepository.Create(context.Background(), tokenOne)
		assert.NoError(t, err)
		assert.NotNil(t, createdTokenOne)
		createdTokenTwo, err := tokenRepository.Create(context.Background(), tokenTwo)
		assert.NoError(t, err)
		assert.NotNil(t, createdTokenTwo)
		_, err = tokenRepository.Create(context.Background(), tokenThree)
		assert.NoError(t, err)
		assert.NotNil(t, createdTokenTwo)
		// insert new key
		newKey, err := tokenRepository.crypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		tokenRepository.internalEncryptionKeys = append(tokenRepository.internalEncryptionKeys, newKey)
		// act
		time.Sleep(time.Second * 60) //mongodb background task that deletes documents runs every 60s
		userTokens, err := tokenRepository.GetUserTokens(context.Background(), &Token{
			UserId: userId,
		})
		// validate response
		// token should expire
		assert.NoError(t, err)
		assert.Equal(t, 0, len(userTokens))
	}
}

func TestGetUserTokensNoUserId(t *testing.T) {
	// setup user client
	tokenRepositoryWithEncryption, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	tokenRepositoryNoEncryption, err := getTestTokenRepository(context.Background(), false, "")
	assert.NoError(t, err)
	clients := []*mongodbRepository{tokenRepositoryWithEncryption, tokenRepositoryNoEncryption}
	expiresAfter := time.Second * 1
	for _, tokenRepository := range clients {
		userId := uuid.NewV4().String()
		device := gofakeit.Phone()
		tokenOne := &Token{
			Id:        uuid.NewV4().String(),
			UserId:    userId,
			Device:    device,
			ExpiresAt: time.Now().Add(expiresAfter),
		}
		tokenTwo := &Token{
			Id:        uuid.NewV4().String(),
			UserId:    userId,
			Device:    device,
			ExpiresAt: time.Now().Add(expiresAfter),
		}
		tokenThree := &Token{
			Id:        uuid.NewV4().String(),
			UserId:    uuid.NewV4().String(),
			Device:    device,
			ExpiresAt: time.Now().Add(expiresAfter),
		}
		// create tokens
		createdTokenOne, err := tokenRepository.Create(context.Background(), tokenOne)
		assert.NoError(t, err)
		assert.NotNil(t, createdTokenOne)
		createdTokenTwo, err := tokenRepository.Create(context.Background(), tokenTwo)
		assert.NoError(t, err)
		assert.NotNil(t, createdTokenTwo)
		_, err = tokenRepository.Create(context.Background(), tokenThree)
		assert.NoError(t, err)
		assert.NotNil(t, createdTokenTwo)
		// insert new key
		newKey, err := tokenRepository.crypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		tokenRepository.internalEncryptionKeys = append(tokenRepository.internalEncryptionKeys, newKey)
		// act
		userTokens, err := tokenRepository.GetUserTokens(context.Background(), &Token{})
		// validate response
		assert.Error(t, err)
		assert.Nil(t, userTokens)
	}
}

func TestGetUserTokensNil(t *testing.T) {
	// setup user client
	tokenRepositoryWithEncryption, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	tokenRepositoryNoEncryption, err := getTestTokenRepository(context.Background(), false, "")
	assert.NoError(t, err)
	clients := []*mongodbRepository{tokenRepositoryWithEncryption, tokenRepositoryNoEncryption}
	expiresAfter := time.Second * 1
	for _, tokenRepository := range clients {
		userId := uuid.NewV4().String()
		device := gofakeit.Phone()
		tokenOne := &Token{
			Id:        uuid.NewV4().String(),
			UserId:    userId,
			Device:    device,
			ExpiresAt: time.Now().Add(expiresAfter),
		}
		tokenTwo := &Token{
			Id:        uuid.NewV4().String(),
			UserId:    userId,
			Device:    device,
			ExpiresAt: time.Now().Add(expiresAfter),
		}
		tokenThree := &Token{
			Id:        uuid.NewV4().String(),
			UserId:    uuid.NewV4().String(),
			Device:    device,
			ExpiresAt: time.Now().Add(expiresAfter),
		}
		// create tokens
		createdTokenOne, err := tokenRepository.Create(context.Background(), tokenOne)
		assert.NoError(t, err)
		assert.NotNil(t, createdTokenOne)
		createdTokenTwo, err := tokenRepository.Create(context.Background(), tokenTwo)
		assert.NoError(t, err)
		assert.NotNil(t, createdTokenTwo)
		_, err = tokenRepository.Create(context.Background(), tokenThree)
		assert.NoError(t, err)
		assert.NotNil(t, createdTokenTwo)
		// insert new key
		newKey, err := tokenRepository.crypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		tokenRepository.internalEncryptionKeys = append(tokenRepository.internalEncryptionKeys, newKey)
		// act
		userTokens, err := tokenRepository.GetUserTokens(context.Background(), nil)
		// validate response
		assert.Error(t, err)
		assert.Nil(t, userTokens)
	}
}
