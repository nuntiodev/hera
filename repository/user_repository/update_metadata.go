package user_repository

import (
	"context"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/x/cryptox"

	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongodbRepository) UpdateMetadata(ctx context.Context, get *go_block.User, update *go_block.User) (*models.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	prepare(actionUpdateMetadata, update)
	if err := r.validate(actionUpdateMetadata, update); err != nil {
		return nil, err
	}
	updateUser := models.ProtoUserToUser(&go_block.User{
		Metadata:  update.Metadata,
		UpdatedAt: update.UpdatedAt,
	})
	if err := r.crypto.Encrypt(updateUser); err != nil {
		return nil, err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"metadata":     updateUser.Metadata,
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
	resp.Metadata = cryptox.Stringx{
		Body:                    update.Metadata,
		InternalEncryptionLevel: resp.Metadata.InternalEncryptionLevel,
		ExternalEncryptionLevel: resp.Metadata.ExternalEncryptionLevel,
	}
	resp.UpdatedAt = update.UpdatedAt.AsTime()
	return &resp, nil
}
