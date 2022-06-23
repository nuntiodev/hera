package user_repository

import (
	"context"
	"github.com/nuntiodev/hera-proto/go_hera"
)

/*
	Delete - this method deletes a specific user with id, email or username.
*/
func (r *mongodbRepository) Delete(ctx context.Context, user *go_hera.User) error {
	if user == nil {
		return UserIsNilErr
	}
	prepare(actionGet, user)
	filter, err := getUserFilter(user)
	if err != nil {
		return err
	}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return NoUsersDeletedErr
	}
	return nil
}
