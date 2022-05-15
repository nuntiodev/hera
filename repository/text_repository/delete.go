package text_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"go.mongodb.org/mongo-driver/bson"
)

func (t *defaultTextRepository) Delete(ctx context.Context, id go_block.LanguageCode) error {
	if id == go_block.LanguageCode_INVALID_LANGUAGE_CODE {
		return errors.New("invalid language code")
	}
	filter := bson.M{}
	filter = bson.M{"_id": id}
	result, err := t.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no document deleted")
	}
	return nil
}
