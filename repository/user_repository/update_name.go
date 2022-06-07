package user_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) UpdateName(ctx context.Context, get *go_block.User, update *go_block.User) (*models.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateName, update)
	if err := r.validate(actionUpdateName, update); err != nil {
		return nil, err
	}
	updateUser := models.ProtoUserToUser(&go_block.User{
		FirstName: update.FirstName,
		LastName:  update.LastName,
		UpdatedAt: update.UpdatedAt,
	})
	// encrypt user if user has previously been encrypted
	if err := r.crypto.Encrypt(&updateUser); err != nil {
		return nil, err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"first_name": updateUser.FirstName,
			"last_name":  updateUser.LastName,
			"updated_at": updateUser.UpdatedAt,
		},
	}
	result := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": get.Id},
		mongoUpdate,
	)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var resp models.User
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	if err := r.crypto.Decrypt(&resp); err != nil {
		return nil, err
	}
	// set updated fields
	resp.FirstName = cryptox.Stringx{
		Body:                    update.FirstName,
		InternalEncryptionLevel: resp.FirstName.InternalEncryptionLevel,
		ExternalEncryptionLevel: resp.FirstName.ExternalEncryptionLevel,
	}
	resp.LastName = cryptox.Stringx{
		Body:                    update.LastName,
		InternalEncryptionLevel: resp.LastName.InternalEncryptionLevel,
		ExternalEncryptionLevel: resp.LastName.ExternalEncryptionLevel,
	}
	resp.UpdatedAt = updateUser.UpdatedAt
	return &resp, nil
}
