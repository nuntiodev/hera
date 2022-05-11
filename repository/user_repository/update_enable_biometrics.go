package user_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (r *mongodbRepository) UpdateEnableBiometrics(ctx context.Context, get *go_block.User, update *go_block.User) (*go_block.User, error) {
	prepare(actionGet, get)
	prepare(actionUpdateEnableBiometrics, update)
	if err := r.validate(actionUpdateEnableBiometrics, update); err != nil {
		return nil, err
	}
	filter, err := getUserFilter(get)
	if err != nil {
		return nil, err
	}
	updateUser := ProtoUserToUser(&go_block.User{
		EnableBiometrics: update.EnableBiometrics,
		UpdatedAt:        update.UpdatedAt,
	})
	mongoUpdate := bson.M{
		"$set": bson.M{
			"enable_biometrics": updateUser.EnableBiometrics,
			"updated_at":        updateUser.UpdatedAt,
		},
	}
	if _, err := r.collection.UpdateOne(
		ctx,
		filter,
		mongoUpdate,
	); err != nil {
		return nil, err
	}
	// set updated fields
	get.EnableBiometrics = updateUser.EnableBiometrics
	get.UpdatedAt = ts.New(updateUser.UpdatedAt)
	return get, nil
}
