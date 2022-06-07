package user_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func (r *mongodbRepository) UpdatePassword(ctx context.Context, get *go_block.User, update *go_block.User) (*models.User, error) {
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
	updateUser := models.ProtoUserToUser(update)
	mongoUpdate := bson.M{
		"$set": bson.M{
			"password":   updateUser.Password,
			"updated_at": updateUser.UpdatedAt,
		},
	}
	filter, err := getUserFilter(get)
	if err != nil {
		return nil, err
	}
	result := r.collection.FindOneAndUpdate(
		ctx,
		filter,
		mongoUpdate,
	)
	if err := result.Err(); err != nil {
		return nil, err
	}
	var resp models.User
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	if err := r.crypto.Decrypt(&resp); err != nil {
		return nil, err
	}
	// set updated fields
	resp.Password = string(hashedPassword)
	return &resp, nil
}
