package token_repository

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/x/cryptox"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsBlockedIEncrypted(t *testing.T) {
	// setup user client
	tokenRepositoryWithEncryption, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	tokenRepositoryNoEncryption, err := getTestTokenRepository(context.Background(), false, "")
	assert.NoError(t, err)
	clients := []*mongodbRepository{tokenRepositoryWithEncryption, tokenRepositoryNoEncryption}
	for _, tokenRepository := range clients {
		userId := uuid.NewV4().String()
		device := gofakeit.Phone()
		token := getToken(&go_block.Token{
			UserId:     userId,
			DeviceInfo: device,
		})
		// create tokens
		createdToken, err := tokenRepository.Create(context.Background(), token)
		assert.NoError(t, err)
		assert.NotNil(t, createdToken)
		// insert new key
		newKey, err := tokenRepository.crypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		tokenRepository.internalEncryptionKeys = append(tokenRepository.internalEncryptionKeys, newKey)
		blockedToken, err := tokenRepository.Block(context.Background(), &go_block.Token{
			Id:     token.Id,
			UserId: token.UserId,
		})
		assert.NoError(t, err)
		assert.NotNil(t, blockedToken)
		assert.True(t, blockedToken.Blocked)
		// act
		isBlocked, err := tokenRepository.IsBlocked(context.Background(), &go_block.Token{
			Id:     token.Id,
			UserId: token.UserId,
		})
		assert.NoError(t, err)
		assert.True(t, isBlocked)
	}
}

func TestIsNotBlockedIEncrypted(t *testing.T) {
	// setup user client
	tokenRepositoryWithEncryption, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	tokenRepositoryNoEncryption, err := getTestTokenRepository(context.Background(), false, "")
	assert.NoError(t, err)
	clients := []*mongodbRepository{tokenRepositoryWithEncryption, tokenRepositoryNoEncryption}
	for _, tokenRepository := range clients {
		userId := uuid.NewV4().String()
		device := gofakeit.Phone()
		token := getToken(&go_block.Token{
			UserId:     userId,
			DeviceInfo: device,
		})
		// create tokens
		createdToken, err := tokenRepository.Create(context.Background(), token)
		assert.NoError(t, err)
		assert.NotNil(t, createdToken)
		// insert new key
		newKey, err := tokenRepository.crypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		tokenRepository.internalEncryptionKeys = append(tokenRepository.internalEncryptionKeys, newKey)
		// act
		isBlocked, err := tokenRepository.IsBlocked(context.Background(), &go_block.Token{
			Id:     token.Id,
			UserId: token.UserId,
		})
		assert.NoError(t, err)
		assert.False(t, isBlocked)
	}
}

func TestIsBlockedEmptyId(t *testing.T) {
	// setup user client
	tokenRepositoryWithEncryption, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	tokenRepositoryNoEncryption, err := getTestTokenRepository(context.Background(), false, "")
	assert.NoError(t, err)
	clients := []*mongodbRepository{tokenRepositoryWithEncryption, tokenRepositoryNoEncryption}
	for _, tokenRepository := range clients {
		userId := uuid.NewV4().String()
		device := gofakeit.Phone()
		token := getToken(&go_block.Token{
			UserId:     userId,
			DeviceInfo: device,
		})
		// create tokens
		createdToken, err := tokenRepository.Create(context.Background(), token)
		assert.NoError(t, err)
		assert.NotNil(t, createdToken)
		// insert new key
		newKey, err := tokenRepository.crypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		tokenRepository.internalEncryptionKeys = append(tokenRepository.internalEncryptionKeys, newKey)
		blockedToken, err := tokenRepository.Block(context.Background(), &go_block.Token{
			Id:     token.Id,
			UserId: token.UserId,
		})
		assert.NoError(t, err)
		assert.NotNil(t, blockedToken)
		assert.True(t, blockedToken.Blocked)
		// act
		isBlocked, err := tokenRepository.IsBlocked(context.Background(), &go_block.Token{})
		assert.Error(t, err)
		assert.False(t, isBlocked)
	}
}

func TestIsBlockedNil(t *testing.T) {
	// setup user client
	tokenRepositoryWithEncryption, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	tokenRepositoryNoEncryption, err := getTestTokenRepository(context.Background(), false, "")
	assert.NoError(t, err)
	clients := []*mongodbRepository{tokenRepositoryWithEncryption, tokenRepositoryNoEncryption}
	for _, tokenRepository := range clients {
		userId := uuid.NewV4().String()
		device := gofakeit.Phone()
		token := getToken(&go_block.Token{
			UserId:     userId,
			DeviceInfo: device,
		})
		// create tokens
		createdToken, err := tokenRepository.Create(context.Background(), token)
		assert.NoError(t, err)
		assert.NotNil(t, createdToken)
		// insert new key
		newKey, err := tokenRepository.crypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		tokenRepository.internalEncryptionKeys = append(tokenRepository.internalEncryptionKeys, newKey)
		blockedToken, err := tokenRepository.Block(context.Background(), &go_block.Token{
			Id:     token.Id,
			UserId: token.UserId,
		})
		assert.NoError(t, err)
		assert.NotNil(t, blockedToken)
		assert.True(t, blockedToken.Blocked)
		// act
		isBlocked, err := tokenRepository.IsBlocked(context.Background(), nil)
		assert.Error(t, err)
		assert.False(t, isBlocked)
	}
}
