package user_repository

import (
	"context"
	"crypto/md5"
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
	actionUpdateProfile
	actionUpdateNamespace
	actionUpdateSecurity
	actionGet
	actionGetAll
)

const (
	maximumLimit = 50
)

var (
	NoUsersDeletedErr = errors.New("no users deleted")
)

type User struct {
	Id                        string    `bson:"id" json:"id"`
	OptionalId                string    `bson:"optional_id" json:"optional_id"`
	Namespace                 string    `bson:"namespace" json:"namespace"`
	Role                      string    `bson:"role" json:"role"`
	Name                      string    `bson:"name" json:"name"`
	Email                     string    `bson:"email" json:"email"`
	Password                  string    `bson:"password" json:"password"`
	Gender                    string    `bson:"gender" json:"gender"`
	Country                   string    `bson:"country" json:"country"`
	Image                     string    `bson:"image" json:"image"`
	Blocked                   bool      `bson:"blocked" json:"blocked"`
	Verified                  bool      `bson:"verified" json:"verified"`
	DisablePasswordValidation bool      `bson:"disable_password_validation" json:"disable_password_validation"`
	Encrypted                 bool      `bson:"encrypted" json:"encrypted"`
	Birthdate                 string    `bson:"birthdate" json:"birthdate"`
	EmailHash                 string    `bson:"email_hash" json:"email_hash"`
	CreatedAt                 time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt                 time.Time `bson:"updated_at" json:"updated_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error)
	UpdatePassword(ctx context.Context, get *block_user.User, update *block_user.User) (*block_user.User, error)
	UpdateProfile(ctx context.Context, get *block_user.User, update *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error)
	UpdateSecurity(ctx context.Context, get *block_user.User, update *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error)
	Get(ctx context.Context, user *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error)
	GetAll(ctx context.Context, userFilter *block_user.UserFilter, namespace string, encryptionOptions *EncryptionOptions) ([]*block_user.User, error)
	Search(ctx context.Context, search string, namespace string, encryptionOptions *EncryptionOptions) ([]*block_user.User, error)
	Delete(ctx context.Context, user *block_user.User) error
	DeleteNamespace(ctx context.Context, namespace string) error
}

type UserMongoRepository struct {
	collection *mongo.Collection
	zapLog     *zap.Logger
}

func NewUserRepository(ctx context.Context, collection *mongo.Collection, zapLog *zap.Logger) (UserRepository, error) {
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
			{Key: "email_hash", Value: 1},
			{Key: "namespace", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(
			bson.D{
				{
					"email_hash", bson.D{
						{
							"$gt", "",
						},
					},
				},
			},
		),
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
	case actionUpdatePassword, actionUpdateProfile, actionUpdateNamespace, actionUpdateSecurity:
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
	case actionGet:
		if user.Id == "" && user.Email == "" && user.OptionalId == "" {
			return errors.New("missing required search parameter")
		}
	case actionCreate:
		if user.Id == "" {
			return errors.New("invalid user id")
		} else if err := checkmail.ValidateFormat(user.Email); user.Email != "" && err != nil {
			return err
		} else if !user.CreatedAt.IsValid() {
			return errors.New("invalid created at date")
		} else if !user.UpdatedAt.IsValid() {
			return errors.New("invalid updated at date")
		} else if err := validatePassword(user.Password); err != nil {
			if user.DisablePasswordValidation == false || (user.Password != "" && user.DisablePasswordValidation == true) {
				return err
			}
		}
	case actionUpdatePassword:
		if err := validatePassword(user.Password); err != nil {
			return err
		} else if !user.UpdatedAt.IsValid() {
			return errors.New("invalid updated at")
		}
	case actionUpdateProfile:
		if err := checkmail.ValidateFormat(user.Email); user.Email != "" && err != nil {
			return err
		} else if !user.UpdatedAt.IsValid() {
			return errors.New("invalid updated at")
		}
	case actionUpdateSecurity:
		if !user.UpdatedAt.IsValid() {
			return errors.New("invalid updated at")
		}
	case actionGetAll:
		return nil
	}
	return nil
}

func (umr *UserMongoRepository) Create(ctx context.Context, user *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error) {
	prepare(actionCreate, user)
	if err := validate(actionCreate, user); err != nil {
		return nil, err
	}
	createUser := protoUserToUser(user)
	if user.Email != "" {
		createUser.EmailHash = fmt.Sprintf("%x", md5.Sum([]byte(user.Email)))
	}
	if encryptionOptions != nil && encryptionOptions.Key != "" {
		if err := createUser.encryptUser(encryptionOptions.Key); err != nil {
			return nil, err
		}
		createUser.Encrypted = true
	}
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		createUser.Password = string(hashedPassword)
		user.Password = string(hashedPassword)
	}
	if _, err := umr.collection.InsertOne(ctx, createUser); err != nil {
		return nil, err
	}
	return user, nil
}

func (umr *UserMongoRepository) UpdatePassword(ctx context.Context, get *block_user.User, update *block_user.User) (*block_user.User, error) {
	prepare(actionGet, get)
	if err := validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdatePassword, update)
	if err := validate(actionUpdatePassword, update); err != nil {
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(update.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	update.Password = string(hashedPassword)
	updateUser := protoUserToUser(update)
	mongoUpdate := bson.M{
		"$set": bson.M{
			"password":   updateUser.Password,
			"updated_at": updateUser.UpdatedAt,
		},
	}
	filter := bson.M{}
	if get.Id != "" {
		filter = bson.M{"id": get.Id, "namespace": get.Namespace}
	} else if get.Email != "" {
		filter = bson.M{"email_hash": fmt.Sprintf("%x", md5.Sum([]byte(get.Email))), "namespace": get.Namespace}
	} else if get.OptionalId != "" {
		filter = bson.M{"optional_id": get.OptionalId, "namespace": get.Namespace}
	}
	updateResult, err := umr.collection.UpdateOne(
		ctx,
		filter,
		mongoUpdate,
	)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("could not find get")
	}
	return update, nil
}

func (umr *UserMongoRepository) UpdateProfile(ctx context.Context, get *block_user.User, update *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error) {
	prepare(actionGet, get)
	if err := validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateProfile, update)
	if err := validate(actionUpdateProfile, update); err != nil {
		return nil, err
	}
	updateUser := protoUserToUser(update)
	if updateUser.Email != "" {
		updateUser.EmailHash = fmt.Sprintf("%x", md5.Sum([]byte(updateUser.Email)))
	}
	// check if user encryption is turned on
	getUser, err := umr.Get(ctx, get, encryptionOptions)
	if err != nil {
		return nil, err
	}
	if getUser.Encrypted == false && (encryptionOptions != nil && encryptionOptions.Key != "") {
		fmt.Println(getUser.Encrypted)
		fmt.Println(encryptionOptions)
		fmt.Println(encryptionOptions.Key)
		return nil, errors.New("you need to update the users security profile (UpdateSecurity) and set encrypted=true if you want to encrypt users data")
	} else if getUser.Encrypted == true && (encryptionOptions == nil || encryptionOptions.Key == "") {
		return nil, errors.New("in order to update an encrypted user, you need to pass the encryption key. If you want to store the user in plaintext, update the users security profile (UpdateSecurity) and turn set encrypted=false")
	} else if getUser.Encrypted && encryptionOptions != nil && encryptionOptions.Key != "" {
		if err := updateUser.encryptUser(encryptionOptions.Key); err != nil {
			return nil, err
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"name":       updateUser.Name,
			"gender":     updateUser.Gender,
			"image":      updateUser.Image,
			"country":    updateUser.Country,
			"email":      updateUser.Email,
			"email_hash": updateUser.EmailHash,
			"birthdate":  updateUser.Birthdate,
			"updated_at": updateUser.UpdatedAt,
		},
	}
	updateResult, err := umr.collection.UpdateOne(
		ctx,
		bson.M{"id": getUser.Id, "namespace": getUser.Namespace},
		mongoUpdate,
	)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("could not find get")
	}
	return get, nil
}

func (umr *UserMongoRepository) UpdateSecurity(ctx context.Context, get *block_user.User, update *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error) {
	prepare(actionGet, get)
	if err := validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateSecurity, update)
	if err := validate(actionUpdateSecurity, update); err != nil {
		return nil, err
	}
	updateUser := protoUserToUser(update)
	// check if user encryption is turned on
	get, err := umr.Get(ctx, get, encryptionOptions)
	if err != nil {
		return nil, err
	}
	getUser := protoUserToUser(get)
	getUser.Role = updateUser.Role
	if getUser.Encrypted == false && encryptionOptions != nil && encryptionOptions.Key != "" {
		if err := getUser.encryptUser(encryptionOptions.Key); err != nil {
			return nil, err
		}
		getUser.Encrypted = true
	} else if getUser.Encrypted == true && encryptionOptions != nil && encryptionOptions.Key == "" {
		if err := getUser.decryptUser(encryptionOptions.Key); err != nil {
			return nil, err
		}
		getUser.Encrypted = false
	}
	getUser.Blocked = updateUser.Blocked
	getUser.Verified = updateUser.Verified
	getUser.DisablePasswordValidation = updateUser.DisablePasswordValidation
	getUser.UpdatedAt = updateUser.UpdatedAt
	mongoUpdate := bson.M{
		"$set": bson.M{
			"role":                        getUser.Role,
			"name":                        getUser.Name,
			"email":                       getUser.Email,
			"gender":                      getUser.Gender,
			"country":                     getUser.Country,
			"image":                       getUser.Image,
			"blocked":                     getUser.Blocked,
			"verified":                    getUser.Verified,
			"disable_password_validation": getUser.DisablePasswordValidation,
			"encrypted":                   getUser.Encrypted,
			"birthdate":                   getUser.Birthdate,
			"updated_at":                  getUser.UpdatedAt,
		},
	}
	updateResult, err := umr.collection.UpdateOne(
		ctx,
		bson.M{"id": getUser.Id, "namespace": getUser.Namespace},
		mongoUpdate,
	)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("could not find get")
	}
	return userToProtoUser(getUser), nil
}

func (umr *UserMongoRepository) Get(ctx context.Context, user *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error) {
	prepare(actionGet, user)
	if err := validate(actionGet, user); err != nil {
		return nil, err
	}
	filter := bson.M{}
	if user.Id != "" {
		filter = bson.M{"id": user.Id, "namespace": user.Namespace}
	} else if user.Email != "" {
		filter = bson.M{"email_hash": fmt.Sprintf("%x", md5.Sum([]byte(user.Email))), "namespace": user.Namespace}
	} else if user.OptionalId != "" {
		filter = bson.M{"optional_id": user.OptionalId, "namespace": user.Namespace}
	}
	resp := User{}
	if err := umr.collection.FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, err
	}
	if resp.Encrypted == true && encryptionOptions != nil && encryptionOptions.Key != "" {
		if err := resp.decryptUser(encryptionOptions.Key); err != nil {
			return nil, err
		}
	}
	return userToProtoUser(&resp), nil
}

func (umr *UserMongoRepository) GetAll(ctx context.Context, userFilter *block_user.UserFilter, namespace string, encryptionOptions *EncryptionOptions) ([]*block_user.User, error) {
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
		if user.Encrypted == true && encryptionOptions != nil && encryptionOptions.Key != "" {
			if err := user.decryptUser(encryptionOptions.Key); err != nil {
				return nil, err
			}
		}
		resp = append(resp, userToProtoUser(&user))
	}

	return resp, nil
}

func (umr *UserMongoRepository) Search(ctx context.Context, search string, namespace string, encryptionOptions *EncryptionOptions) ([]*block_user.User, error) {
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
		if user.Encrypted == true && encryptionOptions != nil && encryptionOptions.Key != "" {
			if err := user.decryptUser(encryptionOptions.Key); err != nil {
				return nil, err
			}
		}
		resp = append(resp, userToProtoUser(&user))
	}

	return resp, nil
}

func (umr *UserMongoRepository) Delete(ctx context.Context, user *block_user.User) error {
	prepare(actionGet, user)
	if err := validate(actionGet, user); err != nil {
		return err
	}
	filter := bson.M{}
	if user.Id != "" {
		filter = bson.M{"id": user.Id, "namespace": user.Namespace}
	} else if user.Email != "" {
		filter = bson.M{"email_hash": fmt.Sprintf("%x", md5.Sum([]byte(user.Email))), "namespace": user.Namespace}
	} else if user.OptionalId != "" {
		filter = bson.M{"optional_id": user.OptionalId, "namespace": user.Namespace}
	}
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
