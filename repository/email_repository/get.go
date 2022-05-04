package email_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (e *defaultEmailRepository) Get(ctx context.Context, email *go_block.Email) (*go_block.Email, error) {
	if email == nil {
		return nil, errors.New("email is nil")
	} else if email.Id == "" {
		return nil, errors.New("missing required id")
	}
	prepare(actionGet, email)
	resp := Email{}
	result := e.collection.FindOne(ctx, bson.M{"_id": email.Id})
	if err := result.Err(); err != nil {
		return nil, err
	}
	if err := result.Decode(&resp); err != nil {
		return nil, err
	}
	if resp.InternalEncryptionLevel > 0 && len(e.internalEncryptionKeys) > 0 {
		if resp.InternalEncryptionLevel > int32(len(e.internalEncryptionKeys)) {
			return nil, errors.New("internal encryption level is illegally higher than amount of internal encryption keys")
		}
		if err := e.DecryptEmail(&resp); err != nil {
			return nil, err
		}
		if resp.InternalEncryptionLevel > int32(len(e.internalEncryptionKeys)) {
			// upgrade user to new internal encryption level
			if err := e.upgradeEncryptionLevel(ctx, resp); err != nil {
				return nil, err
			}
		}
	}
	return EmailToProtoEmail(&resp), nil
}
