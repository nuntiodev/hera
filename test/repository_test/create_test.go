package respository_test

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"github.com/softcorp-io/block-user-service/test/mocks/user_mock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(&block_user.User{
		Name:      gofakeit.Name(),
		Birthdate: ts.Now(),
		Namespace: uuid.NewV4().String(),
		Image:     gofakeit.ImageURL(10, 10),
		Gender:    user_mock.GetRandomGender(),
		Email:     gofakeit.Email(),
	})
	password := user.Password
	user.Id = ""
	// act
	createdUser, err := testRepo.Create(ctx, user)
	assert.Nil(t, err)
	// validate
	assert.NotNil(t, createdUser)
	assert.NotEmpty(t, createdUser.Name)
	assert.NotEmpty(t, createdUser.Email)
	assert.NotEmpty(t, createdUser.Id)
	assert.NotEmpty(t, createdUser.Namespace)
	assert.NotEmpty(t, createdUser.Image)
	assert.NotEqual(t, block_user.Gender_INVALID_GENDER, createdUser.Gender)
	assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte(password)))
	assert.NotEmpty(t, createdUser.Gender)
	assert.True(t, createdUser.Birthdate.IsValid())
	assert.True(t, createdUser.UpdatedAt.IsValid())
	assert.True(t, createdUser.CreatedAt.IsValid())
	// validate in database
	getUser, err := testRepo.GetById(ctx, createdUser)
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
	createdUser, err := testRepo.Create(ctx, user)
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
	user := user_mock.GetRandomUser(&block_user.User{
		DisableAuthentication: true,
	})
	user.Password = ""
	// act
	createdUser, err := testRepo.Create(ctx, user)
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
	createdUser, err := testRepo.Create(ctx, user)
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	if _, err := testRepo.Create(ctx, user); mongo.IsDuplicateKeyError(err) == false {
		t.Fatal(errors.New("creating a user with the same email is not allowed"))
	}
}

func TestCreateDuplicateEmailSameNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	userOne := user_mock.GetRandomUser(&block_user.User{
		Email:     gofakeit.Email(),
		Namespace: uuid.NewV4().String(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Email:     userOne.Email,
		Namespace: userOne.Namespace,
	})
	createdUser, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	if _, err := testRepo.Create(ctx, userTwo); mongo.IsDuplicateKeyError(err) == false {
		t.Fatal(errors.New("creating a user with the same email in the same namespace is not allowed"))
	}
}

func TestCreateDuplicateEmailDifferentNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	userOne := user_mock.GetRandomUser(&block_user.User{
		Email:     gofakeit.Email(),
		Namespace: uuid.NewV4().String(),
	})
	userTwo := user_mock.GetRandomUser(&block_user.User{
		Email:     userOne.Email,
		Namespace: uuid.NewV4().String(),
	})
	createdUser, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	if _, err := testRepo.Create(ctx, userTwo); err != nil {
		t.Fatal(errors.New("creating a user with the same email in different namespaces are allowed"))
	}
}

func TestCreateDuplicateEmptyEmail(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	userOne := user_mock.GetRandomUser(&block_user.User{})
	userTwo := user_mock.GetRandomUser(&block_user.User{})
	createdUser, err := testRepo.Create(ctx, userOne)
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	if _, err := testRepo.Create(ctx, userTwo); err != nil {
		t.Fatal(errors.New("creating two users with empty emails are allowed"))
	}
}

func TestCreateDuplicateOptionalIdSameNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	user.OptionalId = uuid.NewV4().String()
	createdUser, err := testRepo.Create(ctx, user)
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	newUser := user_mock.GetRandomUser(nil)
	newUser.Id = uuid.NewV4().String()
	newUser.OptionalId = user.OptionalId
	newUser.Namespace = user.Namespace
	if _, err := testRepo.Create(ctx, newUser); mongo.IsDuplicateKeyError(err) == false {
		t.Fatal(errors.New("creating a user with the same optional id and same namespace is not allowed"))
	}
}

func TestCreateDuplicateDifferentNamespace(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	user := user_mock.GetRandomUser(nil)
	createdUser, err := testRepo.Create(ctx, user)
	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	// act & validate
	user.Namespace = uuid.NewV4().String()
	user.Id = uuid.NewV4().String()
	if _, err := testRepo.Create(ctx, user); err != nil {
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
	_, err := testRepo.Create(ctx, user)
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
	_, err := testRepo.Create(ctx, user)
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
	createdUser, err := testRepo.Create(ctx, user)
	// validate
	assert.Nil(t, err)
	assert.NotEmpty(t, createdUser.Id)
	assert.Equal(t, id, createdUser.Id)
}
