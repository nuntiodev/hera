package text_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (t *defaultTextRepository) GetAll(ctx context.Context) ([]*go_block.Text, error) {
	var resp []*go_block.Text
	cursor, err := t.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var temp Text
		if err := cursor.Decode(&temp); err != nil {
			return nil, err
		}
		if temp.InternalEncryptionLevel > 0 && len(t.internalEncryptionKeys) > 0 {
			if temp.InternalEncryptionLevel > int32(len(t.internalEncryptionKeys)) {
				return nil, errors.New("internal encryption level is illegally higher than amount of internal encryption keys")
			}
			if err := t.DecryptText(&temp); err != nil {
				return nil, err
			}
			if temp.InternalEncryptionLevel > int32(len(t.internalEncryptionKeys)) {
				// upgrade user to new internal encryption level
				if err := t.upgradeEncryptionLevel(ctx, temp); err != nil {
					return nil, err
				}
			}
		}
		resp = append(resp, TextToProtoText(&temp))
	}
	return resp, nil
}
