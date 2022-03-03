package respository_test

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"github.com/softcorp-io/block-user-service/repository/user_repository"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestUpdateProfile(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Email:   gofakeit.Email(),
		Name:    gofakeit.Name(),
		Image:   gofakeit.ImageURL(10, 10),
		Country: gofakeit.Country(),
		Gender:  block_user.Gender_FEMALE,
	})
	createdUser, err := testRepo.Create(ctx, user, nil)
	initialName := user.Name
	initialImage := user.Image
	initialCountry := user.Country
	initialGender := user.Gender
	initialUpdatedAt := user.UpdatedAt
	initialEmail := user.Email
	initialBlocker := createdUser.Blocked
	assert.Nil(t, err)
	// act
	createdUser.Name = gofakeit.Name()
	createdUser.Image = uuid.NewV4().String()
	createdUser.Country = gofakeit.Country()
	createdUser.Birthdate = ts.Now()
	createdUser.Email = gofakeit.Email()
	createdUser.Gender = block_user.Gender_MALE
	updatedUser, err := testRepo.UpdateProfile(ctx, createdUser, createdUser, nil)
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.NotEmpty(t, updatedUser.Name)
	assert.NotEqual(t, initialName, updatedUser.Name)
	assert.NotEmpty(t, updatedUser.Email)
	assert.NotEqual(t, initialEmail, updatedUser.Email)
	assert.NotEmpty(t, updatedUser.Image)
	assert.NotEqual(t, initialImage, updatedUser.Image)
	assert.NotEmpty(t, updatedUser.Country)
	assert.NotEqual(t, initialCountry, updatedUser.Country)
	assert.NotEqual(t, initialGender, updatedUser.Gender)
	assert.Equal(t, initialBlocker, updatedUser.Blocked)
	assert.True(t, updatedUser.Birthdate.IsValid())
	assert.NotEqual(t, initialUpdatedAt.Nanos, updatedUser.UpdatedAt.Nanos)
	// validate in database
	getUser, err := testRepo.Get(ctx, createdUser, nil)
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
}

func TestUpdateProfileWithEncryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Email:   gofakeit.Email(),
		Name:    gofakeit.Name(),
		Image:   gofakeit.ImageURL(10, 10),
		Country: gofakeit.Country(),
		Gender:  block_user.Gender_FEMALE,
	})
	createdUser, err := testRepo.Create(ctx, user, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	initialName := user.Name
	initialImage := user.Image
	initialCountry := user.Country
	initialGender := user.Gender
	initialUpdatedAt := user.UpdatedAt
	initialEmail := user.Email
	initialBlocker := createdUser.Blocked
	assert.Nil(t, err)
	// act
	createdUser.Name = gofakeit.Name()
	createdUser.Image = uuid.NewV4().String()
	createdUser.Country = gofakeit.Country()
	createdUser.Birthdate = ts.Now()
	createdUser.Email = gofakeit.Email()
	createdUser.Gender = block_user.Gender_MALE
	updatedUser, err := testRepo.UpdateProfile(ctx, createdUser, createdUser, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.NotEmpty(t, updatedUser.Name)
	assert.NotEqual(t, initialName, updatedUser.Name)
	assert.NotEmpty(t, updatedUser.Email)
	assert.NotEqual(t, initialEmail, updatedUser.Email)
	assert.NotEmpty(t, updatedUser.Image)
	assert.NotEqual(t, initialImage, updatedUser.Image)
	assert.NotEmpty(t, updatedUser.Country)
	assert.NotEqual(t, initialCountry, updatedUser.Country)
	assert.NotEqual(t, initialGender, updatedUser.Gender)
	assert.Equal(t, initialBlocker, updatedUser.Blocked)
	assert.True(t, updatedUser.Birthdate.IsValid())
	assert.NotEqual(t, initialUpdatedAt.Nanos, updatedUser.UpdatedAt.Nanos)
	// validate in database
	getUser, err := testRepo.Get(ctx, &block_user.User{
		Id:        createdUser.Id,
		Namespace: createdUser.Namespace,
	}, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
}

func TestUpdateEncryptedProfileWithoutKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Email:   gofakeit.Email(),
		Name:    gofakeit.Name(),
		Image:   gofakeit.ImageURL(10, 10),
		Country: gofakeit.Country(),
		Gender:  block_user.Gender_FEMALE,
	})
	createdUser, err := testRepo.Create(ctx, user, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.Nil(t, err)
	// act
	createdUser.Name = gofakeit.Name()
	createdUser.Image = uuid.NewV4().String()
	createdUser.Country = gofakeit.Country()
	createdUser.Birthdate = ts.Now()
	createdUser.Email = gofakeit.Email()
	createdUser.Gender = block_user.Gender_MALE
	_, err = testRepo.UpdateProfile(ctx, createdUser, createdUser, &user_repository.EncryptionOptions{})
	assert.Error(t, err)
}

func TestUpdateUnencryptedProfileWithKey(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Email:   gofakeit.Email(),
		Name:    gofakeit.Name(),
		Image:   gofakeit.ImageURL(10, 10),
		Country: gofakeit.Country(),
		Gender:  block_user.Gender_FEMALE,
	})
	createdUser, err := testRepo.Create(ctx, user, nil)
	assert.Nil(t, err)
	// act
	createdUser.Name = gofakeit.Name()
	createdUser.Image = uuid.NewV4().String()
	createdUser.Country = gofakeit.Country()
	createdUser.Birthdate = ts.Now()
	createdUser.Email = gofakeit.Email()
	createdUser.Gender = block_user.Gender_MALE
	_, err = testRepo.UpdateProfile(ctx, createdUser, createdUser, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.Error(t, err)
}

func TestUpdateProfileWith(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Email:   gofakeit.Email(),
		Name:    gofakeit.Name(),
		Image:   gofakeit.ImageURL(10, 10),
		Country: gofakeit.Country(),
		Gender:  block_user.Gender_FEMALE,
	})
	createdUser, err := testRepo.Create(ctx, user, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	initialName := user.Name
	initialImage := user.Image
	initialCountry := user.Country
	initialGender := user.Gender
	initialUpdatedAt := user.UpdatedAt
	initialEmail := user.Email
	initialBlocker := createdUser.Blocked
	assert.Nil(t, err)
	// act
	createdUser.Name = gofakeit.Name()
	createdUser.Image = uuid.NewV4().String()
	createdUser.Country = gofakeit.Country()
	createdUser.Birthdate = ts.Now()
	createdUser.Email = gofakeit.Email()
	createdUser.Gender = block_user.Gender_MALE
	updatedUser, err := testRepo.UpdateProfile(ctx, createdUser, createdUser, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	fmt.Println(updatedUser.Id)
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, updatedUser)
	assert.NotEmpty(t, updatedUser.Name)
	assert.NotEqual(t, initialName, updatedUser.Name)
	assert.NotEmpty(t, updatedUser.Email)
	assert.NotEqual(t, initialEmail, updatedUser.Email)
	assert.NotEmpty(t, updatedUser.Image)
	assert.NotEqual(t, initialImage, updatedUser.Image)
	assert.NotEmpty(t, updatedUser.Country)
	assert.NotEqual(t, initialCountry, updatedUser.Country)
	assert.NotEqual(t, initialGender, updatedUser.Gender)
	assert.Equal(t, initialBlocker, updatedUser.Blocked)
	assert.True(t, updatedUser.Birthdate.IsValid())
	assert.NotEqual(t, initialUpdatedAt.Nanos, updatedUser.UpdatedAt.Nanos)
	// validate in database
	getUser, err := testRepo.Get(ctx, &block_user.User{
		Id:        createdUser.Id,
		Namespace: createdUser.Namespace,
	}, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.NoError(t, err)
	assert.NoError(t, user_mock.CompareUsers(getUser, updatedUser))
}

func TestUpdateProfileInvalidNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	createdUser, err := testRepo.Create(ctx, user, nil)
	assert.Nil(t, err)
	// act
	createdUser.Name = gofakeit.Name()
	createdUser.Namespace = uuid.NewV4().String()
	_, err = testRepo.UpdateProfile(ctx, createdUser, createdUser, nil)
	// validate
	assert.Error(t, err)
}
