package user_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/x/cryptox"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) UpdateImage(ctx context.Context, get *go_block.User, update *go_block.User) (*models.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateImage, update)
	if err := r.validate(actionUpdateImage, update); err != nil {
		return nil, err
	}
	updateUser := models.ProtoUserToUser(&go_block.User{
		Image:     update.Image,
		UpdatedAt: update.UpdatedAt,
	})
	if err := r.crypto.Encrypt(updateUser); err != nil {
		return nil, err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"image":        updateUser.Image,
			"updated_at":   updateUser.UpdatedAt,
			"encrypted_at": updateUser.EncryptedAt,
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
	resp.Image = cryptox.Stringx{
		Body:                    update.Image,
		InternalEncryptionLevel: resp.Image.InternalEncryptionLevel,
		ExternalEncryptionLevel: resp.Image.ExternalEncryptionLevel,
	}
	resp.UpdatedAt = updateUser.UpdatedAt
	return &resp, nil
}
