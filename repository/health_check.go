package repository

import "context"

func (r *defaultRepository) performRepositoryHealthCheck(ctx context.Context) error {
	return r.mongodbClient.Ping(ctx, nil)

}
