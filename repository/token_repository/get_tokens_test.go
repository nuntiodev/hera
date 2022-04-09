package token_repository

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/x/cryptox"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetTokensIEncrypted(t *testing.T) {
	// setup user client
	tokenRepositoryWithEncryption, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	tokenRepositoryNoEncryption, err := getTestTokenRepository(context.Background(), false, "")
	assert.NoError(t, err)
	clients := []*mongodbRepository{tokenRepositoryWithEncryption, tokenRepositoryNoEncryption}
	for _, tokenRepository := range clients {
		userId := uuid.NewV4().String()
		device := gofakeit.Phone()
		tokenOne := getToken(&go_block.Token{
			UserId:     userId,
			DeviceInfo: device,
		})
		tokenTwo := getToken(&go_block.Token{
			UserId:     userId,
			DeviceInfo: device,
		})
		tokenThree := getToken(&go_block.Token{
			UserId:     uuid.NewV4().String(),
			DeviceInfo: device,
		})
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
		userTokens, err := tokenRepository.GetTokens(context.Background(), &go_block.Token{
			UserId: userId,
		})
		// validate response
		assert.NoError(t, err)
		assert.NotNil(t, userTokens)
		assert.Equal(t, 2, len(userTokens))
		for _, token := range userTokens {
			assert.False(t, token.Blocked)
			assert.NotEmpty(t, token.Id)
			assert.NotEmpty(t, token.UserId)
			assert.NotEmpty(t, token.DeviceInfo)
			assert.NotEmpty(t, token.CreatedAt.String())
			assert.NotEmpty(t, token.UsedAt.String())
		}
	}
}

func TestCannotGetUserTokensExpired(t *testing.T) {
	// setup user client
	tokenRepositoryWithEncryption, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	tokenRepositoryNoEncryption, err := getTestTokenRepository(context.Background(), false, "")
	assert.NoError(t, err)
	clients := []*mongodbRepository{tokenRepositoryWithEncryption, tokenRepositoryNoEncryption}
	for _, tokenRepository := range clients {
		userId := uuid.NewV4().String()
		device := gofakeit.Phone()
		tokenOne := getToken(&go_block.Token{
			UserId:     userId,
			DeviceInfo: device,
		})
		tokenTwo := getToken(&go_block.Token{
			UserId:     userId,
			DeviceInfo: device,
		})
		tokenThree := getToken(&go_block.Token{
			DeviceInfo: device,
		})
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
		userTokens, err := tokenRepository.GetTokens(context.Background(), &go_block.Token{
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
	for _, tokenRepository := range clients {
		userId := uuid.NewV4().String()
		device := gofakeit.Phone()
		tokenOne := getToken(&go_block.Token{
			UserId:     userId,
			DeviceInfo: device,
		})
		tokenTwo := getToken(&go_block.Token{
			UserId:     userId,
			DeviceInfo: device,
		})
		tokenThree := getToken(&go_block.Token{
			DeviceInfo: device,
		})
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
		userTokens, err := tokenRepository.GetTokens(context.Background(), &go_block.Token{})
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
	for _, tokenRepository := range clients {
		userId := uuid.NewV4().String()
		device := gofakeit.Phone()
		tokenOne := getToken(&go_block.Token{
			UserId:     userId,
			DeviceInfo: device,
		})
		tokenTwo := getToken(&go_block.Token{
			UserId:     userId,
			DeviceInfo: device,
		})
		tokenThree := getToken(&go_block.Token{
			DeviceInfo: device,
		})
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
		userTokens, err := tokenRepository.GetTokens(context.Background(), nil)
		// validate response
		assert.Error(t, err)
		assert.Nil(t, userTokens)
	}
}