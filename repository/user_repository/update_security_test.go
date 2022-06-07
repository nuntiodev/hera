package user_repository

/*
import (
	"context"
	"github.com/google/uuid"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/x/cryptox"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateSecurityIEEncrypted(t *testing.T) {
	// setup available clients
	var clients []*mongodbRepository
	ns := uuid.NewString()
	userRepositoryFullEncryption, err := getTestUserRepository(context.Background(), true, true, ns)
	assert.NoError(t, err)
	userRepositoryInternalEncryption, err := getTestUserRepository(context.Background(), true, false, ns)
	assert.NoError(t, err)
	userRepositoryExternalEncryption, err := getTestUserRepository(context.Background(), false, true, ns)
	assert.NoError(t, err)
	clients = []*mongodbRepository{userRepositoryFullEncryption, userRepositoryInternalEncryption, userRepositoryExternalEncryption}
	assert.NoError(t, err)
	for index, userRepository := range clients {
		userOne := getTestUser()
		dbUserOne, err := userRepository.Create(context.Background(), &userOne)
		assert.NoError(t, err)
		assert.NotNil(t, dbUserOne)
		// set new internal and external encryption key
		encryptionKey, err := cryptox.GenerateSymmetricKey(32, cryptox.AlphaNum)
		assert.NoError(t, err)
		// internal
		internalKeys, _ := userRepository.crypto.GetInternalEncryptionKeys()
		internalKeys = append(internalKeys, encryptionKey)
		assert.NoError(t, userRepository.crypto.SetInternalEncryptionKeys(internalKeys))
		// external
		externalKeys, _ := userRepository.crypto.GetExternalEncryptionKeys()
		externalKeys = append(externalKeys, encryptionKey)
		assert.NoError(t, userRepository.crypto.SetExternalEncryptionKeys(externalKeys))
		// act 1
		updatedUser, err := userRepository.UpdateSecurity(context.Background(), models.UserToProtoUser(dbUserOne))
		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		// assert that update has been propagated correctly to database
		getUser, err := userRepository.Get(context.Background(), models.UserToProtoUser(updatedUser))
		assert.NoError(t, err, index)
		assert.NotNil(t, getUser)
		// validate encryption levels
		internal, external := userRepository.crypto.EncryptionLevel(getUser)
		assert.Equal(t, int32(0), internal, getUser)
		assert.Equal(t, int32(0), external, index)
		// act 2
		updatedUser, err = userRepository.UpdateSecurity(context.Background(), models.UserToProtoUser(dbUserOne))
		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		// assert that update has been propagated correctly to database
		getUser, err = userRepository.Get(context.Background(), models.UserToProtoUser(updatedUser))
		assert.NoError(t, err, index)
		assert.NotNil(t, getUser)
		// validate encryption levels
		internal, external = userRepository.crypto.EncryptionLevel(getUser)
		assert.Equal(t, int32(len(internalKeys)), internal)
		assert.Equal(t, int32(len(externalKeys)), external)
}
}
*/
