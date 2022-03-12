package respository_test

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block"
	"github.com/softcorp-io/block-user-service/repository/user_repository"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
		Email:     gofakeit.Email(),
		Role:      gofakeit.Name(),
	})
	password := user.Password
	user.Id = ""
	// act
	createdUser, err := testRepo.Create(ctx, user, nil)
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, createdUser)
	assert.NotEmpty(t, createdUser.Email)
	assert.NotEmpty(t, createdUser.Id)
	assert.NotEmpty(t, createdUser.Namespace)
	assert.NotEmpty(t, createdUser.Image)
	assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte(password)))
	assert.True(t, createdUser.UpdatedAt.IsValid())
	assert.True(t, createdUser.CreatedAt.IsValid())
	// validate in database
	getUser, err := testRepo.Get(ctx, createdUser, nil)
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	assert.Nil(t, user_mock.CompareUsers(getUser, createdUser))
}

func TestCreateWithEncryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
		Email:     gofakeit.Email(),
		Role:      gofakeit.Name(),
	})
	password := user.Password
	user.Id = ""
	// act
	createdUser, err := testRepo.Create(ctx, user, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, createdUser)
	assert.NotEmpty(t, createdUser.Email)
	assert.NotEmpty(t, createdUser.Id)
	assert.NotEmpty(t, createdUser.Namespace)
	assert.NotEmpty(t, createdUser.Image)
	assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte(password)))
	assert.True(t, createdUser.UpdatedAt.IsValid())
	assert.True(t, createdUser.CreatedAt.IsValid())
	// validate in database
	getUser, err := testRepo.Get(ctx, createdUser, &user_repository.EncryptionOptions{
		Key: encryptionKey,
	})
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	assert.Nil(t, user_mock.CompareUsers(getUser, createdUser))
}

func TestCreateWithEmptyFields(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	password := user.Password
	user.Id = ""
	// act
	createdUser, err := testRepo.Create(ctx, user, nil)
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, createdUser)
	assert.NotEmpty(t, createdUser.Id)
	assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte(password)))
	assert.True(t, createdUser.UpdatedAt.IsValid())
	assert.True(t, createdUser.CreatedAt.IsValid())
}

func TestCreateWithEmptyPasswordDisableAuth(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	user.Password = ""
	// act
	createdUser, err := testRepo.Create(ctx, user, nil)
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, createdUser)
	assert.NotEmpty(t, createdUser.Id)
	assert.True(t, createdUser.UpdatedAt.IsValid())
	assert.True(t, createdUser.CreatedAt.IsValid())
}

func TestCreateDuplicateIdSameNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	createdUser, err := testRepo.Create(ctx, user, nil)
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	if _, err := testRepo.Create(ctx, user, nil); mongo.IsDuplicateKeyError(err) == false {
		t.Fatal(errors.New("creating a user with the same email is not allowed"))
	}
}

func TestCreateDuplicateEmailSameNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Email:     gofakeit.Email(),
		Namespace: uuid.NewV4().String(),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Email:     userOne.Email,
		Namespace: userOne.Namespace,
	})
	createdUser, err := testRepo.Create(ctx, userOne, nil)
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	if _, err := testRepo.Create(ctx, userTwo, nil); mongo.IsDuplicateKeyError(err) == false {
		t.Fatal(errors.New("creating a user with the same email in the same namespace is not allowed"))
	}
}

func TestCreateDuplicateEmailDifferentNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	userOne := user_mock.GetRandomUser(&go_block.User{
		Email:     gofakeit.Email(),
		Namespace: uuid.NewV4().String(),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Email:     userOne.Email,
		Namespace: uuid.NewV4().String(),
	})
	createdUser, err := testRepo.Create(ctx, userOne, nil)
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	if _, err := testRepo.Create(ctx, userTwo, nil); err != nil {
		t.Fatal(errors.New("creating a user with the same email in different namespaces are allowed"))
	}
}

func TestCreateDuplicateEmptyEmail(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	userOne := user_mock.GetRandomUser(&go_block.User{})
	userTwo := user_mock.GetRandomUser(&go_block.User{})
	createdUser, err := testRepo.Create(ctx, userOne, nil)
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	if _, err := testRepo.Create(ctx, userTwo, nil); err != nil {
		t.Fatal(errors.New("creating two users with empty emails are allowed"))
	}
}

func TestCreateDuplicateOptionalIdSameNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	user.OptionalId = uuid.NewV4().String()
	createdUser, err := testRepo.Create(ctx, user, nil)
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	newUser := user_mock.GetRandomUser(nil)
	newUser.Id = uuid.NewV4().String()
	newUser.OptionalId = user.OptionalId
	newUser.Namespace = user.Namespace
	if _, err := testRepo.Create(ctx, newUser, nil); mongo.IsDuplicateKeyError(err) == false {
		t.Fatal(errors.New("creating a user with the same optional id and same namespace is not allowed"))
	}
}

func TestCreateDuplicateDifferentNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	createdUser, err := testRepo.Create(ctx, user, nil)
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	user.Namespace = uuid.NewV4().String()
	user.Id = uuid.NewV4().String()
	if _, err := testRepo.Create(ctx, user, nil); err != nil {
		t.Fatal(errors.New("creating users with the same email in two different namespaces is allowed"))
	}
}

func TestCreateInvalidEmail(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	// act
	user.Email = "softcorp@@test.io"
	_, err := testRepo.Create(ctx, user, nil)
	// validate
	assert.Error(t, err)
}

func TestCreateInvalidPassword(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	// act
	user.Password = "Test1234"
	_, err := testRepo.Create(ctx, user, nil)
	// validate
	assert.Error(t, err)
}

func TestCreateOverwriteId(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	user.Id = uuid.NewV4().String()
	id := user.Id
	// act
	createdUser, err := testRepo.Create(ctx, user, nil)
	// validate
	assert.Nil(t, err)
	assert.NotEmpty(t, createdUser.Id)
	assert.Equal(t, id, createdUser.Id)
}
