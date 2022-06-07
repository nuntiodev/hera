package token_repository

/*
import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/stretchr/testify/assert"
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
		token := getToken(&go_block.Token{
			DeviceInfo: device,
		})
		// act
		createdToken, err := tokenRepository.Create(context.Background(), token)
		assert.NoError(t, err)
		assert.NotNil(t, createdToken)
		// assert new fields are present
		assert.NotEmpty(t, token.CreatedAt.String())
		assert.NotEmpty(t, token.UsedAt.String())
		assert.Equal(t, len(tokenRepository.internalEncryptionKeys), int(token.InternalEncryptionLevel))
		assert.Equal(t, token.Id, createdToken.Id)
		assert.Equal(t, token.UserId, createdToken.UserId)
		assert.Equal(t, token.ExpiresAt, createdToken.ExpiresAt)
		assert.NotEmpty(t, createdToken.DeviceInfo)
		assert.Equal(t, device, createdToken.DeviceInfo)
	}
}

func TestCreateNoId(t *testing.T) {
	// setup user client
	tokenRepository, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	token := getToken(nil)
	token.Id = ""
	// act
	createdToken, err := tokenRepository.Create(context.Background(), token)
	assert.Error(t, err)
	assert.Nil(t, createdToken)
}

func TestCreateNoUserId(t *testing.T) {
	// setup user client
	tokenRepository, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	token := getToken(nil)
	token.UserId = ""
	// act
	createdToken, err := tokenRepository.Create(context.Background(), token)
	assert.Error(t, err)
	assert.Nil(t, createdToken)
}

func TestCreateNoUserExpiresAt(t *testing.T) {
	// setup user client
	tokenRepository, err := getTestTokenRepository(context.Background(), true, "")
	assert.NoError(t, err)
	token := getToken(nil)
	token.ExpiresAt = nil
	// act
	createdToken, err := tokenRepository.Create(context.Background(), token)
	assert.Error(t, err)
	assert.Nil(t, createdToken)
}
*/
