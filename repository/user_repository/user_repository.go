package user_repository

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/badoux/checkmail"
	uuid "github.com/satori/go.uuid"
	"github.com/softcorp-io/block-proto/go_block/block_user"
	"github.com/softcorp-io/block-user-service/crypto"
	"go.mongodb.org/mongo-driver/bson"
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
	actionUpdateOptionalId
	actionUpdateImage
	actionUpdateMetadata
	actionUpdateNamespace
	actionUpdateSecurity
	actionGet
	actionGetAll
)

const (
	maximumGetLimit = 75
	maxFieldLength  = 150
)

var (
	NoUsersDeletedErr = errors.New("no users deleted")
)

type User struct {
	Id         string    `bson:"id" json:"id"`
	OptionalId string    `bson:"optional_id" json:"optional_id"`
	Namespace  string    `bson:"namespace" json:"namespace"`
	Email      string    `bson:"email" json:"email"`
	Role       string    `bson:"role" json:"role"`
	Password   string    `bson:"password" json:"password"`
	Image      string    `bson:"image" json:"image"`
	Encrypted  bool      `bson:"encrypted" json:"encrypted"`
	EmailHash  string    `bson:"email_hash" json:"email_hash"`
	Metadata   string    `bson:"metadata" json:"metadata"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error)
	UpdatePassword(ctx context.Context, get *block_user.User, update *block_user.User) (*block_user.User, error)
	UpdateEmail(ctx context.Context, get *block_user.User, update *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error)
	UpdateOptionalId(ctx context.Context, get *block_user.User, update *block_user.User) (*block_user.User, error)
	UpdateImage(ctx context.Context, get *block_user.User, update *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error)
	UpdateMetadata(ctx context.Context, get *block_user.User, update *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error)
	UpdateSecurity(ctx context.Context, get *block_user.User, update *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error)
	Get(ctx context.Context, user *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error)
	GetAll(ctx context.Context, userFilter *block_user.UserFilter, namespace string, encryptionOptions *EncryptionOptions) ([]*block_user.User, error)
	Count(ctx context.Context, namespace string) (int64, error)
	Delete(ctx context.Context, user *block_user.User) error
	DeleteNamespace(ctx context.Context, namespace string) error
}

type mongoRepository struct {
	collection   *mongo.Collection
	crypto       crypto.Crypto
	metadataType block_user.MetadataType
	zapLog       *zap.Logger
}

func NewUserRepository(ctx context.Context, collection *mongo.Collection, crypto crypto.Crypto, metadataType block_user.MetadataType, zapLog *zap.Logger) (UserRepository, error) {
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
	return &mongoRepository{
		collection:   collection,
		zapLog:       zapLog,
		crypto:       crypto,
		metadataType: metadataType,
	}, nil
}

func prepare(action int, user *block_user.User) {
	if user == nil {
		return
	}
	switch action {
	case actionCreate:
		user.CreatedAt = ts.Now()
		user.UpdatedAt = ts.Now()
		if user.Id == "" {
			user.Id = uuid.NewV4().String()
		}
	case actionUpdatePassword, actionUpdateImage, actionUpdateMetadata,
		actionUpdateNamespace, actionUpdateSecurity, actionUpdateEmail,
		actionUpdateOptionalId:
		user.UpdatedAt = ts.Now()
	}
	user.Id = strings.TrimSpace(user.Id)
	user.Namespace = strings.TrimSpace(user.Namespace)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.Image = strings.TrimSpace(user.Image)
	user.OptionalId = strings.TrimSpace(user.OptionalId)
	user.Metadata = strings.TrimSpace(user.Metadata)
}

func (r *mongoRepository) validate(action int, user *block_user.User) error {
	if user == nil {
		return errors.New("user is nil")
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
		} else if err := validatePassword(user.Password); err != nil && user.Password != "" {
			return err
		} else if !user.CreatedAt.IsValid() {
			return errors.New("invalid created at date")
		} else if !user.UpdatedAt.IsValid() {
			return errors.New("invalid updated at date")
		} else if r.metadataType == block_user.MetadataType_METADATA_TYPE_JSON && !json.Valid([]byte(user.Metadata)) && user.Metadata != "" {
			return errors.New("invalid json type")
		}
	case actionUpdatePassword:
		if err := validatePassword(user.Password); user.Password != "" && err != nil {
			return err
		} else if !user.UpdatedAt.IsValid() {
			return errors.New("invalid updated at")
		}
	case actionUpdateEmail:
		if err := checkmail.ValidateFormat(user.Email); user.Email != "" && err != nil {
			return err
		} else if !user.UpdatedAt.IsValid() {
			return errors.New("invalid updated at")
		}
	case actionUpdateMetadata:
		if !user.UpdatedAt.IsValid() {
			return errors.New("invalid updated at")
		} else if r.metadataType == block_user.MetadataType_METADATA_TYPE_JSON && !json.Valid([]byte(user.Metadata)) && user.Metadata != "" {
			return errors.New("invalid json type")
		}
	case actionUpdateSecurity:
		if !user.UpdatedAt.IsValid() {
			return errors.New("invalid updated at")
		}
	case actionGetAll, actionUpdateOptionalId:
		return nil
	}
	if len(user.Email) > maxFieldLength {

	} else if len(user.Role) > maxFieldLength {

	} else if len(user.OptionalId) > maxFieldLength {

	} else if len(user.Metadata) > 10*maxFieldLength {

	}
	return nil
}

func (r *mongoRepository) Create(ctx context.Context, user *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error) {
	prepare(actionCreate, user)
	if err := r.validate(actionCreate, user); err != nil {
		return nil, err
	}
	createUser := protoUserToUser(user)
	if user.Email != "" {
		createUser.EmailHash = fmt.Sprintf("%x", md5.Sum([]byte(user.Email)))
	}
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		createUser.Password = string(hashedPassword)
		user.Password = string(hashedPassword)
	}
	if encryptionOptions != nil && encryptionOptions.Key != "" {
		if err := r.encryptUser(encryptionOptions.Key, createUser); err != nil {
			return nil, err
		}
		createUser.Encrypted = true
		user.Encrypted = true
	} else {
		createUser.Encrypted = false
		user.Encrypted = false
	}
	if _, err := r.collection.InsertOne(ctx, createUser); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *mongoRepository) UpdatePassword(ctx context.Context, get *block_user.User, update *block_user.User) (*block_user.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdatePassword, update)
	if err := r.validate(actionUpdatePassword, update); err != nil {
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
	updateResult, err := r.collection.UpdateOne(
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

func (r *mongoRepository) UpdateEmail(ctx context.Context, get *block_user.User, update *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateEmail, update)
	if err := r.validate(actionUpdateEmail, update); err != nil {
		return nil, err
	}
	updateUser := protoUserToUser(&block_user.User{
		Email:     update.Email,
		UpdatedAt: update.UpdatedAt,
	})
	if updateUser.Email != "" {
		updateUser.EmailHash = fmt.Sprintf("%x", md5.Sum([]byte(updateUser.Email)))
	}
	// check if user encryption is turned on
	getUser, err := r.Get(ctx, get, encryptionOptions)
	if err != nil {
		return nil, err
	}
	if getUser.Encrypted == false && (encryptionOptions != nil && encryptionOptions.Key != "") {
		return nil, errors.New("you need to update the users security profile (UpdateSecurity) and set encrypted=true if you want to encrypt users data")
	} else if getUser.Encrypted == true && (encryptionOptions == nil || encryptionOptions.Key == "") {
		return nil, errors.New("in order to update an encrypted user, you need to pass the encryption key. If you want to store the user in plaintext, update the users security profile (UpdateSecurity) and turn set encrypted=false")
	} else if getUser.Encrypted && encryptionOptions != nil && encryptionOptions.Key != "" {
		if err := r.encryptUser(encryptionOptions.Key, updateUser); err != nil {
			return nil, err
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"email":      updateUser.Email,
			"email_hash": fmt.Sprintf("%x", md5.Sum([]byte(update.Email))),
			"updated_at": updateUser.UpdatedAt,
		},
	}
	updateResult, err := r.collection.UpdateOne(
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
	return update, nil
}

func (r *mongoRepository) UpdateOptionalId(ctx context.Context, get *block_user.User, update *block_user.User) (*block_user.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateOptionalId, update)
	if err := r.validate(actionUpdateOptionalId, update); err != nil {
		return nil, err
	}
	updateUser := protoUserToUser(&block_user.User{
		OptionalId: update.OptionalId,
		UpdatedAt:  update.UpdatedAt,
	})
	mongoUpdate := bson.M{
		"$set": bson.M{
			"optional_id": updateUser.OptionalId,
			"updated_at":  updateUser.UpdatedAt,
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
	updateResult, err := r.collection.UpdateOne(
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

func (r *mongoRepository) UpdateImage(ctx context.Context, get *block_user.User, update *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateImage, update)
	if err := r.validate(actionUpdateImage, update); err != nil {
		return nil, err
	}
	updateUser := protoUserToUser(&block_user.User{
		Image:     update.Image,
		UpdatedAt: update.UpdatedAt,
	})
	// check if user encryption is turned on
	getUser, err := r.Get(ctx, get, encryptionOptions)
	if err != nil {
		return nil, err
	}
	if getUser.Encrypted == false && (encryptionOptions != nil && encryptionOptions.Key != "") {
		return nil, errors.New("you need to update the users security profile (UpdateSecurity) and set encrypted=true if you want to encrypt users data")
	} else if getUser.Encrypted == true && (encryptionOptions == nil || encryptionOptions.Key == "") {
		return nil, errors.New("in order to update an encrypted user, you need to pass the encryption key. If you want to store the user in plaintext, update the users security profile (UpdateSecurity) and turn set encrypted=false")
	} else if getUser.Encrypted && encryptionOptions != nil && encryptionOptions.Key != "" {
		if err := r.encryptUser(encryptionOptions.Key, updateUser); err != nil {
			return nil, err
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"image":      updateUser.Image,
			"updated_at": updateUser.UpdatedAt,
		},
	}
	updateResult, err := r.collection.UpdateOne(
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
	return update, nil
}

func (r *mongoRepository) UpdateMetadata(ctx context.Context, get *block_user.User, update *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateMetadata, update)
	if err := r.validate(actionUpdateMetadata, update); err != nil {
		return nil, err
	}
	updateUser := protoUserToUser(&block_user.User{
		Metadata:  update.Metadata,
		UpdatedAt: update.UpdatedAt,
	}) // check if user encryption is turned on
	getUser, err := r.Get(ctx, get, encryptionOptions)
	if err != nil {
		return nil, err
	}
	if getUser.Encrypted == false && (encryptionOptions != nil && encryptionOptions.Key != "") {
		return nil, errors.New("you need to update the users security profile (UpdateSecurity) and set encrypted=true if you want to encrypt users data")
	} else if getUser.Encrypted == true && (encryptionOptions == nil || encryptionOptions.Key == "") {
		return nil, errors.New("in order to update an encrypted user, you need to pass the encryption key. If you want to store the user in plaintext, update the users security profile (UpdateSecurity) and turn set encrypted=false")
	} else if getUser.Encrypted && encryptionOptions != nil && encryptionOptions.Key != "" {
		if err := r.encryptUser(encryptionOptions.Key, updateUser); err != nil {
			return nil, err
		}
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"metadata":   updateUser.Metadata,
			"updated_at": updateUser.UpdatedAt,
		},
	}
	updateResult, err := r.collection.UpdateOne(
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
	return update, nil
}

func (r *mongoRepository) UpdateSecurity(ctx context.Context, get *block_user.User, update *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateSecurity, update)
	if err := r.validate(actionUpdateSecurity, update); err != nil {
		return nil, err
	}
	updateUser := protoUserToUser(update)
	// check if user encryption is turned on
	get, err := r.Get(ctx, get, encryptionOptions)
	if err != nil {
		return nil, err
	}
	getUser := protoUserToUser(get)
	getUser.Role = updateUser.Role
	if getUser.Encrypted == false && encryptionOptions != nil && encryptionOptions.Key != "" {
		if err := r.encryptUser(encryptionOptions.Key, getUser); err != nil {
			return nil, err
		}
		getUser.Encrypted = true
	} else if getUser.Encrypted == true && encryptionOptions != nil && encryptionOptions.Key == "" {
		if err := r.decryptUser(encryptionOptions.Key, getUser); err != nil {
			return nil, err
		}
		getUser.Encrypted = false
	}
	getUser.UpdatedAt = updateUser.UpdatedAt
	mongoUpdate := bson.M{
		"$set": bson.M{
			"email":      getUser.Email,
			"image":      getUser.Image,
			"role":       getUser.Role,
			"encrypted":  getUser.Encrypted,
			"metadata":   getUser.Metadata,
			"updated_at": getUser.UpdatedAt,
		},
	}
	updateResult, err := r.collection.UpdateOne(
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

func (r *mongoRepository) Get(ctx context.Context, user *block_user.User, encryptionOptions *EncryptionOptions) (*block_user.User, error) {
	prepare(actionGet, user)
	if err := r.validate(actionGet, user); err != nil {
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
	if err := r.collection.FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, err
	}
	if resp.Encrypted == true && encryptionOptions != nil && encryptionOptions.Key != "" {
		if err := r.decryptUser(encryptionOptions.Key, &resp); err != nil {
			return nil, err
		}
	}
	return userToProtoUser(&resp), nil
}

func (r *mongoRepository) GetAll(ctx context.Context, userFilter *block_user.UserFilter, namespace string, encryptionOptions *EncryptionOptions) ([]*block_user.User, error) {
	var resp []*block_user.User
	sortOptions := options.FindOptions{}
	limitOptions := options.Find()
	limitOptions.SetLimit(maximumGetLimit)
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
		default:
			return nil, errors.New("invalid sorting")
		}
		if userFilter.From >= 0 && userFilter.To > 0 {
			if userFilter.To-userFilter.From > maximumGetLimit {
				return nil, errors.New(fmt.Sprintf("exceeding maximum range of %d", maximumGetLimit))
			}
			limitOptions.SetLimit(int64(userFilter.To - userFilter.From))
			limitOptions.SetSkip(int64(userFilter.From))
		}
	}
	cursor, err := r.collection.Find(ctx, filter, &sortOptions, limitOptions)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		if user.Encrypted == true && encryptionOptions != nil && encryptionOptions.Key != "" {
			if err := r.decryptUser(encryptionOptions.Key, &user); err != nil {
				return nil, err
			}
		}
		resp = append(resp, userToProtoUser(&user))
	}

	return resp, nil
}

func (r *mongoRepository) Count(ctx context.Context, namespace string) (int64, error) {
	filter := bson.M{"namespace": namespace}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *mongoRepository) Delete(ctx context.Context, user *block_user.User) error {
	prepare(actionGet, user)
	if err := r.validate(actionGet, user); err != nil {
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
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return NoUsersDeletedErr
	}
	return nil
}

func (r *mongoRepository) DeleteNamespace(ctx context.Context, namespace string) error {
	filter := bson.M{"namespace": namespace}
	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return NoUsersDeletedErr
	}
	return nil
}

/*
func (r *mongoRepository) Search(ctx context.Context, search string, namespace string, encryptionOptions *EncryptionOptions) ([]*block_user.User, error) {
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
	cursor, err := r.collection.Find(ctx, filter, &sortOptions, limitOptions)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		if user.Encrypted == true && encryptionOptions != nil && encryptionOptions.Key != "" {
			if err := r.decryptUser(encryptionOptions.Key, &user); err != nil {
				return nil, err
			}
		}
		resp = append(resp, userToProtoUser(&user))
	}

	return resp, nil
}
*/
