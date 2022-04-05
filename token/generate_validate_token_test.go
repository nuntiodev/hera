package token

import (
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/x/cryptox"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/*
	TestGenerateValidateToken generates and validates a JWT signed and validated by a public/private keypair
*/
func TestGenerateValidateToken(t *testing.T) {
	// setup crypto
	c, err := cryptox.New()
	assert.NoError(t, err)
	// setup token
	to, err := New()
	assert.NoError(t, err)
	// generate rsa keys
	privateKey, publicKey, err := c.GenerateRsaKeyPair(2048)
	assert.NoError(t, err)
	// data to validate
	userId := uuid.NewV4().String()
	refreshId := uuid.NewV4().String()
	expiresAfter := time.Second * 10
	// act one - generate token
	token, claims, err := to.GenerateToken(privateKey, userId, refreshId, TokenTypeAccess, expiresAfter)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.NotEmpty(t, token)
	// validate one - check data matches
	assert.Equal(t, claims.UserId, userId)
	assert.Equal(t, claims.RefreshTokenId, refreshId)
	// act two - validate token
	validatedClaims, err := to.ValidateToken(publicKey, token)
	assert.NoError(t, err)
	assert.NotNil(t, validatedClaims)
	assert.Equal(t, validatedClaims.Id, claims.Id)
	assert.Equal(t, validatedClaims.UserId, claims.UserId)
	assert.Equal(t, validatedClaims.RefreshTokenId, claims.RefreshTokenId)
}

func TestGenerateValidateTokenInvalidKey(t *testing.T) {
	// setup crypto
	c, err := cryptox.New()
	assert.NoError(t, err)
	// setup token
	to, err := New()
	assert.NoError(t, err)
	// generate rsa keys
	privateKey, _, err := c.GenerateRsaKeyPair(2048)
	assert.NoError(t, err)
	// data to validate
	userId := uuid.NewV4().String()
	refreshId := uuid.NewV4().String()
	expiresAfter := time.Second * 10
	// act one - generate token
	token, claims, err := to.GenerateToken(privateKey, userId, refreshId, TokenTypeAccess, expiresAfter)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.NotEmpty(t, token)
	// validate one - check data matches
	assert.Equal(t, claims.UserId, userId)
	assert.Equal(t, claims.RefreshTokenId, refreshId)
	// act two - validate token with invalid public key
	_, invalidPublicKey, err := c.GenerateRsaKeyPair(2048)
	assert.NoError(t, err)
	validatedClaims, err := to.ValidateToken(invalidPublicKey, token)
	assert.Error(t, err)
	assert.Nil(t, validatedClaims)
}

func TestGenerateTokenEmptyRefreshId(t *testing.T) {
	// setup crypto
	c, err := cryptox.New()
	assert.NoError(t, err)
	// setup token
	to, err := New()
	assert.NoError(t, err)
	// generate rsa keys
	privateKey, _, err := c.GenerateRsaKeyPair(2048)
	assert.NoError(t, err)
	// data to validate
	expiresAfter := time.Second * 10
	// act one - generate token
	token, claims, err := to.GenerateToken(privateKey, uuid.NewV4().String(), "", TokenTypeAccess, expiresAfter)
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Empty(t, token)
}

func TestGenerateTokenInvalidType(t *testing.T) {
	// setup crypto
	c, err := cryptox.New()
	assert.NoError(t, err)
	// setup token
	to, err := New()
	assert.NoError(t, err)
	// generate rsa keys
	privateKey, _, err := c.GenerateRsaKeyPair(2048)
	assert.NoError(t, err)
	// data to validate
	expiresAfter := time.Second * 10
	// act one - generate token
	token, claims, err := to.GenerateToken(privateKey, uuid.NewV4().String(), uuid.NewV4().String(), "invalid", expiresAfter)
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Empty(t, token)
}
