package email_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (e *defaultEmailRepository) Delete(ctx context.Context, email *go_block.Email) error {
	if email == nil {
		return errors.New("email is nil")
	} else if email.Id == "" {
		return errors.New("missing required id")
	}
	prepare(actionDelete, email)
	filter := bson.M{"_id": email.Id}
	if _, err := e.collection.DeleteOne(ctx, filter); err != nil {
		return err
	}
	return nil
}
