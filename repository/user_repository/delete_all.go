package user_repository

import (
	"context"
)

/*
	DeleteAll - this method deletes all users in the collection.
*/
func (r *mongodbRepository) DeleteAll(ctx context.Context) error {
	if err := r.collection.Drop(ctx); err != nil {
		return err
	}
	return nil
}
