package user_repository

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/softcorp-io/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func (r *mongodbRepository) UpdatePassword(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error) {
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
	updateUser := ProtoUserToUser(update)
	mongoUpdate := bson.M{
		"$set": bson.M{
			"password":   updateUser.Password,
			"updated_at": updateUser.UpdatedAt,
		},
	}
	filter := bson.M{}
	if get.Id != "" {
		filter = bson.M{"_id": get.Id}
	} else if get.Email != "" {
		filter = bson.M{"email_hash": fmt.Sprintf("%x", md5.Sum([]byte(get.Email)))}
	} else if get.OptionalId != "" {
		filter = bson.M{"optional_id": get.OptionalId}
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
