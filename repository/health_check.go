package repository

import "context"

func (r *Repository) performRepositoryHealthCheck(ctx context.Context) error {
	return r.mongoClient.Ping(ctx, nil)

}
