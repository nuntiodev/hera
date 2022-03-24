package user_repository

import (
	"context"
	"errors"
	"github.com/softcorp-io/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (r *mongoRepository) UpdateSecurity(ctx context.Context, get *go_block.User, encryptionKey string) (*go_block.User, error) {
	prepare(actionGet, get)
	if err := r.validate(actionGet, get); err != nil {
		return nil, err
	}
	get, err := r.Get(ctx, get, "") // check if user encryption is turned on
	if err != nil {
		return nil, err
	}
	if get.Encrypted == false && encryptionKey != "" {
		if err := r.crypto.EncryptUser(encryptionKey, get); err != nil {
			return nil, err
		}
		get.Encrypted = true
		get.EncryptedAt = ts.Now()
	} else if get.Encrypted == true && encryptionKey != "" {
		if err := r.crypto.DecryptUser(encryptionKey, get); err != nil {
			return nil, err
		}
		get.Encrypted = false
		get.EncryptedAt = &ts.Timestamp{}
	}
	get.UpdatedAt = ts.Now()
	updateUser := ProtoUserToUser(get)
	mongoUpdate := bson.M{
		"$set": bson.M{
			"email":        updateUser.Email,
			"image":        updateUser.Image,
			"encrypted":    updateUser.Encrypted,
			"metadata":     updateUser.Metadata,
			"updated_at":   updateUser.UpdatedAt,
			"encrypted_at": updateUser.EncryptedAt,
		},
	}
	updateResult, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": get.Id},
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
