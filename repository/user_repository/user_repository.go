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
	"time"
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
	actionGetByOptionalId
	actionGetAll
	actionDelete
)

const (
	maximumLimit = 50
)

var (
	NoUsersDeletedErr = errors.New("no users deleted")
)

type User struct {
	Id         string            `bson:"id" json:"id"`
	OptionalId string            `bson:"optional_id" json:"optional_id"`
	Namespace  string            `bson:"namespace" json:"namespace"`
	Role       string            `bson:"role" json:"role"`
	Name       string            `bson:"name" json:"name"`
	Email      string            `bson:"email" json:"email"`
	Password   string            `bson:"password" json:"password"`
	Gender     block_user.Gender `bson:"gender" json:"gender"`
	Country    string            `bson:"country" json:"country"`
	Image      string            `bson:"image" json:"image"`
	Blocked    bool              `bson:"blocked" json:"blocked"`
	Birthdate  time.Time         `bson:"birthdate" json:"birthdate"`
	CreatedAt  time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time         `bson:"updated_at" json:"updated_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *block_user.User) (*block_user.User, error)
	UpdatePassword(ctx context.Context, user *block_user.User) (*block_user.User, error)
	UpdateEmail(ctx context.Context, user *block_user.User) (*block_user.User, error)
	UpdateProfile(ctx context.Context, user *block_user.User) (*block_user.User, error)
	UpdateNamespace(ctx context.Context, user *block_user.User) (*block_user.User, error)
	UpdateSecurity(ctx context.Context, user *block_user.User) (*block_user.User, error)
	GetById(ctx context.Context, user *block_user.User) (*block_user.User, error)
	GetByEmail(ctx context.Context, user *block_user.User) (*block_user.User, error)
	GetByOptionalId(ctx context.Context, user *block_user.User) (*block_user.User, error)
	GetAll(ctx context.Context, userFilter *block_user.UserFilter, namespace string) ([]*block_user.User, error)
	Search(ctx context.Context, search string, namespace string) ([]*block_user.User, error)
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
	optionalIdIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "optional_id", Value: 1},
			{Key: "namespace", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(
			bson.D{
				{
					"optional_id", bson.D{
						{
							"$gt", "",
						},
					},
				},
			},
		),
	}
	if _, err := collection.Indexes().CreateOne(ctx, optionalIdIndexModel); err != nil {
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
	user.OptionalId = strings.TrimSpace(user.OptionalId)
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
	case actionGetByOptionalId:
		if user.OptionalId == "" {
			return errors.New("missing required optional id")
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
	if _, err := umr.collection.InsertOne(ctx, protoUserToUser(user)); err != nil {
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
	updateUser := protoUserToUser(user)
	update := bson.M{
		"$set": bson.M{
			"password":   updateUser.Password,
			"updated_at": updateUser.UpdatedAt,
		},
	}
	filter := bson.M{"id": user.Id, "namespace": user.Namespace}
	updateResult, err := umr.collection.UpdateOne(
		ctx,
		filter,
		update,
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
	updateUser := protoUserToUser(user)
	update := bson.M{
		"$set": bson.M{
			"email":      updateUser.Email,
			"updated_at": updateUser.UpdatedAt,
		},
	}
	filter := bson.M{"id": user.Id, "namespace": user.Namespace}
	updateResult, err := umr.collection.UpdateOne(
		ctx,
		filter,
		update,
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
	updateUser := protoUserToUser(user)
	update := bson.M{
		"$set": bson.M{
			"name":       updateUser.Name,
			"gender":     updateUser.Gender,
			"image":      updateUser.Image,
			"country":    updateUser.Country,
			"birthdate":  updateUser.Birthdate,
			"updated_at": updateUser.UpdatedAt,
		},
	}
	filter := bson.M{"id": user.Id, "namespace": user.Namespace}
	updateResult, err := umr.collection.UpdateOne(
		ctx,
		filter,
		update,
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
	updateUser := protoUserToUser(user)
	update := bson.M{
		"$set": bson.M{
			"namespace":  updateUser.Namespace,
			"updated_at": updateUser.UpdatedAt,
		},
	}
	filter := bson.M{"id": user.Id}
	updateResult, err := umr.collection.UpdateOne(
		ctx,
		filter,
		update,
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
	updateUser := protoUserToUser(user)
	update := bson.M{
		"$set": bson.M{
			"role":       updateUser.Role,
			"blocked":    updateUser.Blocked,
			"updated_at": updateUser.UpdatedAt,
		},
	}
	filter := bson.M{"id": user.Id}
	updateResult, err := umr.collection.UpdateOne(
		ctx,
		filter,
		update,
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
	resp := User{}
	if err := umr.collection.FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, err
	}
	return userToProtoUser(&resp), nil
}

func (umr *UserMongoRepository) GetByEmail(ctx context.Context, user *block_user.User) (*block_user.User, error) {
	prepare(actionGetByEmail, user)
	if err := validate(actionGetByEmail, user); err != nil {
		return nil, err
	}
	filter := bson.M{"id": user.Id, "namespace": user.Namespace}
	resp := User{}
	if err := umr.collection.FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, err
	}
	return userToProtoUser(&resp), nil
}

func (umr *UserMongoRepository) GetByOptionalId(ctx context.Context, user *block_user.User) (*block_user.User, error) {
	prepare(actionGetByOptionalId, user)
	if err := validate(actionGetByOptionalId, user); err != nil {
		return nil, err
	}
	filter := bson.M{"optional_id": user.OptionalId, "namespace": user.Namespace}
	resp := User{}
	if err := umr.collection.FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, err
	}
	return userToProtoUser(&resp), nil
}

func (umr *UserMongoRepository) GetAll(ctx context.Context, userFilter *block_user.UserFilter, namespace string) ([]*block_user.User, error) {
	var resp []*block_user.User
	sortOptions := options.FindOptions{}
	limitOptions := options.Find()
	limitOptions.SetLimit(maximumLimit)
	filter := bson.M{"namespace": namespace}
	if userFilter != nil {
		order := -1
		if userFilter.Order == block_user.UserFilter_INC {
			order = 1
		}
		switch userFilter.Sort {
		case block_user.UserFilter_CREATED_AT:
			sortOptions.SetSort(bson.D{{"created_at", order}, {"_id", order}})
		case block_user.UserFilter_UPDATE_AT:
			sortOptions.SetSort(bson.D{{"updated_at", order}, {"_id", order}})
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
	}
	cursor, err := umr.collection.Find(ctx, filter, &sortOptions, limitOptions)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		resp = append(resp, userToProtoUser(&user))
	}

	return resp, nil
}

func (umr *UserMongoRepository) Search(ctx context.Context, search string, namespace string) ([]*block_user.User, error) {
	if search == "" {
		return nil, errors.New("empty search string")
	}
	var resp []*block_user.User
	limitOptions := options.Find()
	limitOptions.SetLimit(50)
	filter := bson.D{
		{"namespace", namespace},
		{"$or", bson.A{
			bson.D{{"id", primitive.Regex{Pattern: search, Options: ""}}},
			bson.D{{"optional_id", primitive.Regex{Pattern: search, Options: ""}}},
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
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		resp = append(resp, userToProtoUser(&user))
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
