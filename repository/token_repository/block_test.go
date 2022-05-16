package token_repository

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/x/cryptox"
	"github.com/stretchr/testify/assert"
)

/*
	TestBlockTokenIEncrypted - positive service test.
	this test checks that we are actually able to block a token; both with encryption and without.
*/
func TestBlockTokenIEncrypted(t *testing.T) {
	// setup user client
	tokenRepositoryWithEncryption, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	tokenRepositoryNoEncryption, err := getTestTokenRepository(context.Background(), false, "")
	assert.NoError(t, err)
	clients := []*mongodbRepository{tokenRepositoryWithEncryption, tokenRepositoryNoEncryption}
	for _, tokenRepository := range clients {
		userId := uuid.NewString()
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
		blockedToken, err := tokenRepository.Block(context.Background(), &go_block.Token{
			Id:     createdToken.Id,
			UserId: token.UserId,
		})
		// validate response
		assert.NoError(t, err)
		assert.NotNil(t, blockedToken)
		assert.True(t, blockedToken.Blocked)
		// validate in database
		getToken, err := tokenRepository.Get(context.Background(), &go_block.Token{
			Id: token.Id,
		})
		assert.NoError(t, err)
		assert.NotNil(t, getToken)
		assert.True(t, getToken.Blocked)
	}
}

/*
	TestBlockTokenIEncrypted - exploratory service test.
	this test makes sure that we cannot block a token without an id.
*/
func TestBlockTokenNoId(t *testing.T) {
	// setup user client
	tokenRepositoryWithEncryption, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	tokenRepositoryNoEncryption, err := getTestTokenRepository(context.Background(), false, "")
	assert.NoError(t, err)
	clients := []*mongodbRepository{tokenRepositoryWithEncryption, tokenRepositoryNoEncryption}
	for _, tokenRepository := range clients {
		userId := uuid.NewString()
		device := gofakeit.Phone()
		token := getToken(&go_block.Token{
			Id:         uuid.NewString(),
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
		blockedToken, err := tokenRepository.Block(context.Background(), &go_block.Token{})
		// validate response
		assert.Error(t, err)
		assert.Nil(t, blockedToken)
	}
}

/*
	TestBlockTokenIEncrypted - exploratory service test.
	this test makes sure we throw a valid err when token is nil.
*/
func TestBlockTokenNil(t *testing.T) {
	// setup user client
	tokenRepositoryWithEncryption, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	tokenRepositoryNoEncryption, err := getTestTokenRepository(context.Background(), false, "")
	assert.NoError(t, err)
	clients := []*mongodbRepository{tokenRepositoryWithEncryption, tokenRepositoryNoEncryption}
	for _, tokenRepository := range clients {
		userId := uuid.NewString()
		device := gofakeit.Phone()
		token := getToken(&go_block.Token{
			Id:         uuid.NewString(),
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
		blockedToken, err := tokenRepository.Block(context.Background(), nil)
		// validate response
		assert.Error(t, err)
		assert.Nil(t, blockedToken)
	}
}
