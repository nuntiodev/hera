package respository_test

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block"
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
		Image: gofakeit.ImageURL(10, 10),
		Email: gofakeit.Email(),
	})
	password := user.Password
	// act
	users, err := testRepository.Users(ctx, namespace)
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, "")
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, createdUser)
	assert.NotEmpty(t, createdUser.Email)
	assert.NotEmpty(t, createdUser.Id)
	assert.NotEmpty(t, createdUser.Image)
	assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte(password)))
	assert.True(t, createdUser.UpdatedAt.IsValid())
	assert.True(t, createdUser.CreatedAt.IsValid())
	// validate in database
	getUser, err := users.Get(ctx, createdUser, "")
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Nil(t, user_mock.CompareUsers(getUser, createdUser))
}

func TestCreateWithEncryption(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&go_block.User{
		Image: gofakeit.ImageURL(10, 10),
		Email: gofakeit.Email(),
	})
	password := user.Password
	// act
	users, err := testRepository.Users(ctx, namespace)
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, encryptionKey)
	assert.NoError(t, err)
	// validate
	assert.NotNil(t, createdUser)
	assert.NotEmpty(t, createdUser.Email)
	assert.NotEmpty(t, createdUser.Id)
	assert.NotEmpty(t, createdUser.Image)
	assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte(password)))
	assert.True(t, createdUser.UpdatedAt.IsValid())
	assert.True(t, createdUser.CreatedAt.IsValid())
	// validate in database
	getUser, err := users.Get(ctx, createdUser, encryptionKey)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Nil(t, user_mock.CompareUsers(getUser, createdUser))
}

func TestCreateWithEmptyFields(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	password := user.Password
	// act
	users, err := testRepository.Users(ctx, namespace)
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, "")
	assert.NoError(t, err)
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
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, "")
	assert.NoError(t, err)
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
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	user := user_mock.GetRandomUser(nil)
	createdUser, err := users.Create(ctx, user, "")
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	if _, err := users.Create(ctx, user, ""); mongo.IsDuplicateKeyError(err) == false {
		t.Fatal(errors.New("creating a user with the same email is not allowed"))
	}
}

func TestCreateDuplicateEmailSameNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	userOne := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Email: userOne.Email,
	})
	createdUser, err := users.Create(ctx, userOne, "")
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	if _, err := users.Create(ctx, userTwo, ""); mongo.IsDuplicateKeyError(err) == false {
		t.Fatal(errors.New("creating a user with the same email in the same namespace is not allowed"))
	}
}

func TestCreateDuplicateEmailDifferentNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	usersOne, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	usersTwo, err := testRepository.Users(ctx, "")
	assert.NoError(t, err)
	userOne := user_mock.GetRandomUser(&go_block.User{
		Email: gofakeit.Email(),
	})
	userTwo := user_mock.GetRandomUser(&go_block.User{
		Email: userOne.Email,
	})
	createdUser, err := usersOne.Create(ctx, userOne, "")
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	if _, err := usersTwo.Create(ctx, userTwo, ""); err != nil {
		t.Fatal(errors.New("creating a user with the same email in different namespaces are allowed"))
	}
}

func TestCreateDuplicateEmptyEmail(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	userOne := user_mock.GetRandomUser(&go_block.User{})
	userTwo := user_mock.GetRandomUser(&go_block.User{})
	createdUser, err := users.Create(ctx, userOne, "")
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	if _, err := users.Create(ctx, userTwo, ""); err != nil {
		t.Fatal(errors.New("creating two users with empty emails are allowed"))
	}
}

func TestCreateDuplicateOptionalIdSameNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	user.OptionalId = uuid.NewV4().String()
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	createdUser, err := users.Create(ctx, user, "")
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	newUser := user_mock.GetRandomUser(nil)
	newUser.Id = uuid.NewV4().String()
	newUser.OptionalId = user.OptionalId
	if _, err := users.Create(ctx, newUser, ""); mongo.IsDuplicateKeyError(err) == false {
		t.Fatal(errors.New("creating a user with the same optional id and same namespace is not allowed"))
	}
}

func TestCreateDuplicateDifferentNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	usersOne, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	usersTwo, err := testRepository.Users(ctx, "")
	assert.NoError(t, err)
	createdUser, err := usersOne.Create(ctx, user, "")
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	user.Id = uuid.NewV4().String()
	if _, err := usersTwo.Create(ctx, user, ""); err != nil {
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
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	_, err = users.Create(ctx, user, "")
	// validate
	assert.Error(t, err)
}

func TestCreateInvalidPassword(t *testing.T) {
	t.Skipf("currently not validating password")
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	// act
	user.Password = "Test1234"
	users, err := testRepository.Users(ctx, uuid.NewV4().String())
	assert.NoError(t, err)
	_, err = users.Create(ctx, user, "")
	// validate
	assert.Error(t, err)
}
