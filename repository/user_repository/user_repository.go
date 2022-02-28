package user_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/badoux/checkmail"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

const (
	actionCreate = iota
	actionUpdatePassword
	actionUpdateEmail
	actionUpdateProfile
	actionUpdateNamespace
	actionUpdateSecurity
	actionGetByID
	actionGetByEmail
	actionGetAll
	actionDelete
)

const (
	maximumLimit = 50
)

var (
	NoUsersDeletedErr = errors.New("no users deleted")
)

type UserRepository interface {
	Create(ctx context.Context, user *block_user.User) (*block_user.User, error)
	UpdatePassword(ctx context.Context, user *block_user.User) (*block_user.User, error)
	UpdateEmail(ctx context.Context, user *block_user.User) (*block_user.User, error)
	UpdateProfile(ctx context.Context, user *block_user.User) (*block_user.User, error)
	UpdateNamespace(ctx context.Context, user *block_user.User) (*block_user.User, error)
	UpdateSecurity(ctx context.Context, user *block_user.User) (*block_user.User, error)
	GetById(ctx context.Context, user *block_user.User) (*block_user.User, error)
	GetByEmail(ctx context.Context, user *block_user.User) (*block_user.User, error)
	GetAll(ctx context.Context, userFilter *block_user.UserFilter) ([]*block_user.User, error)
	Search(ctx context.Context, search string, userFilter *block_user.UserFilter) ([]*block_user.User, error)
	Delete(ctx context.Context, user *block_user.User) error
	DeleteNamespace(ctx context.Context, namespace string) error
}

type UserMongoRepository struct {
	collection *mongo.Collection
	zapLog     *zap.Logger
}

func NewUserRepository(ctx context.Context, collection *mongo.Collection, zapLog *zap.Logger) (*UserMongoRepository, error) {
	zapLog.Info("creating user repository...")
	idNamespaceIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "id", Value: 1},
			{Key: "namespace", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	if _, err := collection.Indexes().CreateOne(ctx, idNamespaceIndexModel); err != nil {
		return nil, err
	}
	emailNamespaceIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
			{Key: "namespace", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	if _, err := collection.Indexes().CreateOne(ctx, emailNamespaceIndexModel); err != nil {
		return nil, err
	}
	return &UserMongoRepository{
		collection: collection,
		zapLog:     zapLog,
	}, nil
}

func prepare(action int, user *block_user.User) {
	if user == nil {
		return
	}
	switch action {
	case actionCreate:
		user.Blocked = false
		user.CreatedAt = ts.Now()
		user.UpdatedAt = ts.Now()
		if user.Id == "" {
			user.Id = uuid.NewV4().String()
		}
	case actionUpdatePassword, actionUpdateEmail, actionUpdateProfile, actionUpdateNamespace, actionUpdateSecurity:
		user.UpdatedAt = ts.Now()
	}
	user.Id = strings.TrimSpace(user.Id)
	user.Namespace = strings.TrimSpace(user.Namespace)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.Name = strings.TrimSpace(user.Name)
	user.Country = strings.TrimSpace(user.Country)
	user.Image = strings.TrimSpace(user.Image)
}

func validate(action int, user *block_user.User) error {
	if user == nil {
		return errors.New("user is nil")
	} else if len(user.Name) > 100 {
		return errors.New("invalid name")
	} else if len(user.Country) > 100 {
		return errors.New("invalid country")
	}
	switch action {
	case actionCreate:
		if user.Id == "" {
			return errors.New("invalid user id")
		} else if err := checkmail.ValidateFormat(user.Email); err != nil {
			return err
		} else if !user.CreatedAt.IsValid() {
			return errors.New("invalid created at date")
		} else if !user.UpdatedAt.IsValid() {
			return errors.New("invalid updated at date")
		} else if err := validatePassword(user.Password); err != nil {
			return err
		}
	case actionUpdatePassword:
		if user.Id == "" {
			return errors.New("invalid user id")
		} else if err := validatePassword(user.Password); err != nil {
			return err
		}
	case actionUpdateProfile, actionUpdateSecurity, actionGetByID, actionDelete:
		if user.Id == "" {
			return errors.New("invalid user id")
		}
	case actionUpdateEmail:
		if user.Id == "" {
			return errors.New("invalid user id")
		} else if err := checkmail.ValidateFormat(user.Email); err != nil {
			return err
		}
	case actionGetByEmail:
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return err
		}
	case actionGetAll:
		return nil
	case actionUpdateNamespace:
		if user.Id == "" {
			return errors.New("invalid user id")
		} else if user.Namespace == "" {
			return errors.New("invalid namespace id")
		}
	}
	return nil
}

func (umr *UserMongoRepository) Create(ctx context.Context, user *block_user.User) (*block_user.User, error) {
	prepare(actionCreate, user)
	if err := validate(actionCreate, user); err != nil {
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	if err != nil {
		return nil, err
	}
	if _, err := umr.collection.InsertOne(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (umr *UserMongoRepository) UpdatePassword(ctx context.Context, user *block_user.User) (*block_user.User, error) {
	prepare(actionUpdatePassword, user)
	if err := validate(actionUpdatePassword, user); err != nil {
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	updateUser := bson.M{
		"$set": bson.M{
			"password":  user.Password,
			"updatedAt": user.UpdatedAt,
		},
	}
	filter := bson.M{"id": user.Id, "namespace": user.Namespace}
	updateResult, err := umr.collection.UpdateOne(
		ctx,
		filter,
		updateUser,
	)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("could not find user")
	}
	return user, nil
}

func (umr *UserMongoRepository) UpdateEmail(ctx context.Context, user *block_user.User) (*block_user.User, error) {
	prepare(actionUpdateEmail, user)
	if err := validate(actionUpdateEmail, user); err != nil {
		return nil, err
	}
	updateUser := bson.M{
		"$set": bson.M{
			"email":     user.Email,
			"updatedAt": user.UpdatedAt,
		},
	}
	filter := bson.M{"id": user.Id, "namespace": user.Namespace}
	updateResult, err := umr.collection.UpdateOne(
		ctx,
		filter,
		updateUser,
	)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("could not find user")
	}
	return user, nil
}

func (umr *UserMongoRepository) UpdateProfile(ctx context.Context, user *block_user.User) (*block_user.User, error) {
	prepare(actionUpdateProfile, user)
	if err := validate(actionUpdateProfile, user); err != nil {
		return nil, err
	}
	updateUser := bson.M{
		"$set": bson.M{
			"name":      user.Name,
			"gender":    user.Gender,
			"image":     user.Image,
			"country":   user.Country,
			"birthdate": user.Birthdate,
			"updatedAt": user.UpdatedAt,
		},
	}
	filter := bson.M{"id": user.Id, "namespace": user.Namespace}
	updateResult, err := umr.collection.UpdateOne(
		ctx,
		filter,
		updateUser,
	)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("could not find user")
	}
	return user, nil
}

func (umr *UserMongoRepository) UpdateNamespace(ctx context.Context, user *block_user.User) (*block_user.User, error) {
	prepare(actionUpdateNamespace, user)
	if err := validate(actionUpdateNamespace, user); err != nil {
		return nil, err
	}
	updateUser := bson.M{
		"$set": bson.M{
			"namespace": user.Namespace,
			"updatedAt": user.UpdatedAt,
		},
	}
	filter := bson.M{"id": user.Id}
	updateResult, err := umr.collection.UpdateOne(
		ctx,
		filter,
		updateUser,
	)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("could not find user")
	}
	return user, nil
}

func (umr *UserMongoRepository) UpdateSecurity(ctx context.Context, user *block_user.User) (*block_user.User, error) {
	prepare(actionUpdateSecurity, user)
	if err := validate(actionUpdateSecurity, user); err != nil {
		return nil, err
	}
	updateUser := bson.M{
		"$set": bson.M{
			"role":      user.Role,
			"blocked":   user.Blocked,
			"updatedAt": user.UpdatedAt,
		},
	}
	filter := bson.M{"id": user.Id}
	updateResult, err := umr.collection.UpdateOne(
		ctx,
		filter,
		updateUser,
	)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("could not find user")
	}
	return user, nil
}

func (umr *UserMongoRepository) GetById(ctx context.Context, user *block_user.User) (*block_user.User, error) {
	prepare(actionGetByID, user)
	if err := validate(actionGetByID, user); err != nil {
		return nil, err
	}
	filter := bson.M{"id": user.Id, "namespace": user.Namespace}
	resp := block_user.User{}
	if err := umr.collection.FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (umr *UserMongoRepository) GetByEmail(ctx context.Context, user *block_user.User) (*block_user.User, error) {
	prepare(actionGetByEmail, user)
	if err := validate(actionGetByEmail, user); err != nil {
		return nil, err
	}
	filter := bson.M{"id": user.Id, "namespace": user.Namespace}
	resp := block_user.User{}
	if err := umr.collection.FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (umr *UserMongoRepository) GetAll(ctx context.Context, userFilter *block_user.UserFilter) ([]*block_user.User, error) {
	var resp []*block_user.User
	sortOptions := options.FindOptions{}
	limitOptions := options.Find()
	limitOptions.SetLimit(maximumLimit)
	filter := bson.M{"namespace": ""}
	if userFilter != nil {
		order := -1
		if userFilter.Order == block_user.UserFilter_INC {
			order = 1
		}
		switch userFilter.Sort {
		case block_user.UserFilter_CREATED_AT:
			sortOptions.SetSort(bson.D{{"createdAt", order}, {"_id", order}})
		case block_user.UserFilter_UPDATE_AT:
			sortOptions.SetSort(bson.D{{"updatedAt", order}, {"_id", order}})
		case block_user.UserFilter_BIRTHDATE:
			sortOptions.SetSort(bson.D{{"birthdate", order}, {"_id", order}})
		case block_user.UserFilter_NAME:
			sortOptions.SetSort(bson.D{{"name", order}, {"_id", order}})
		default:
			return nil, errors.New("invalid sorting")
		}
		if userFilter.From >= 0 && userFilter.To > 0 {
			if userFilter.To-userFilter.From > maximumLimit {
				return nil, errors.New(fmt.Sprintf("exceeding maximum range of %d", maximumLimit))
			}
			limitOptions.SetLimit(int64(userFilter.To - userFilter.From))
			limitOptions.SetSkip(int64(userFilter.From))
		}
		if userFilter.Namespace != "" {
			filter = bson.M{"namespace": userFilter.Namespace}
		}
	}
	cursor, err := umr.collection.Find(ctx, filter, &sortOptions, limitOptions)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var user block_user.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		resp = append(resp, &user)
	}

	return resp, nil
}

func (umr *UserMongoRepository) Search(ctx context.Context, search string, userFilter *block_user.UserFilter) ([]*block_user.User, error) {
	if search == "" {
		return nil, errors.New("empty search string")
	}
	var resp []*block_user.User
	namespace := ""
	limitOptions := options.Find()
	limitOptions.SetLimit(50)
	if userFilter != nil {
		namespace = userFilter.Namespace
	}
	filter := bson.D{
		{"namespace", namespace},
		{"$or", bson.A{
			bson.D{{"id", primitive.Regex{Pattern: search, Options: ""}}},
			bson.D{{"email", primitive.Regex{Pattern: search, Options: ""}}},
			bson.D{{"name", primitive.Regex{Pattern: search, Options: ""}}},
			bson.D{{"country", primitive.Regex{Pattern: search, Options: ""}}},
		},
		},
	}
	sortOptions := options.FindOptions{}
	sortOptions.SetSort(bson.D{{"name", 1}, {"_id", 1}})
	cursor, err := umr.collection.Find(ctx, filter, &sortOptions, limitOptions)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var user block_user.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		resp = append(resp, &user)
	}

	return resp, nil
}

func (umr *UserMongoRepository) Delete(ctx context.Context, user *block_user.User) error {
	prepare(actionDelete, user)
	if err := validate(actionDelete, user); err != nil {
		return err
	}
	filter := bson.M{"id": user.Id, "namespace": user.Namespace}
	result, err := umr.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return NoUsersDeletedErr
	}
	return nil
}

func (umr *UserMongoRepository) DeleteNamespace(ctx context.Context, namespace string) error {
	filter := bson.M{"namespace": namespace}
	result, err := umr.collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return NoUsersDeletedErr
	}
	return nil
}
