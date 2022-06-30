package user_repository

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mongodbRepository) Search(ctx context.Context, search string) (*go_hera.User, error) {
	if search == "" {
		return nil, errors.New("missing required search parameter")
	}
	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"_id", primitive.Regex{Pattern: search, Options: ""}}},
			bson.D{{"email_hash", primitive.Regex{Pattern: fmt.Sprintf("%x", md5.Sum([]byte(search))), Options: ""}}},
			bson.D{{"username_hash", primitive.Regex{Pattern: fmt.Sprintf("%x", md5.Sum([]byte(search))), Options: ""}}},
			bson.D{{"phone_hash", primitive.Regex{Pattern: fmt.Sprintf("%x", md5.Sum([]byte(search))), Options: ""}}},
		},
		},
	}
	resp := models.User{}
	if err := r.collection.FindOne(ctx, filter).Decode(&resp); err != nil {
		return nil, fmt.Errorf("could not find user with id: %v and err: %v", filter, err)
	}
	if err := r.crypto.Decrypt(&resp); err != nil {
		return nil, err
	}
	// check if we should upgrade the encryption level
	if upgradable, _ := r.crypto.Upgradeble(&resp); upgradable {
		if err := r.upgradeEncryptionLevel(ctx, &resp); err != nil {
			return nil, err
		}
	}
	return models.UserToProtoUser(&resp), nil
}
