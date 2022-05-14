package text_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (t *defaultTextRepository) Get(ctx context.Context, id go_block.LanguageCode) (*go_block.Text, error) {
	resp := Text{}
	result := t.collection.FindOne(ctx, bson.M{"_id": id})
	if err := result.Err(); err != nil {
		return nil, err
	}
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	if resp.InternalEncryptionLevel > 0 && len(t.internalEncryptionKeys) > 0 {
		if resp.InternalEncryptionLevel > int32(len(t.internalEncryptionKeys)) {
			return nil, errors.New("internal encryption level is illegally higher than amount of internal encryption keys")
		}
		if err := t.DecryptText(&resp); err != nil {
			return nil, err
		}
		if resp.InternalEncryptionLevel > int32(len(t.internalEncryptionKeys)) {
			// upgrade user to new internal encryption level
			if err := t.upgradeEncryptionLevel(ctx, resp); err != nil {
				return nil, err
			}
		}
	}
	return TextToProtoText(&resp), nil
}
