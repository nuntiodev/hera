package email_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (e *defaultEmailRepository) GetAll(ctx context.Context, email *go_block.Email) ([]*go_block.Email, error) {
	if email == nil {
		return nil, errors.New("email is nil")
	} else if email.Id == "" {
		return nil, errors.New("missing required id")
	}
	prepare(actionGet, email)
	resp := []*go_block.Email{}
	cursor, err := e.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		temp := Email{}
		if err := cursor.Decode(&temp); err != nil {
			return nil, err
		}
		if temp.InternalEncryptionLevel > 0 && len(e.internalEncryptionKeys) > 0 {
			if temp.InternalEncryptionLevel > int32(len(e.internalEncryptionKeys)) {
				return nil, errors.New("internal encryption level is illegally higher than amount of internal encryption keys")
			}
			if err := e.DecryptEmail(&temp); err != nil {
				return nil, err
			}
		}
		resp = append(resp, EmailToProtoEmail(&temp))
	}
	return resp, nil
}
