package user_repository

import (
	"context"
)

func (r *mongoRepository) DeleteAll(ctx context.Context) error {
	if err := r.collection.Drop(ctx); err != nil {
		return err
	}
	return nil
}
