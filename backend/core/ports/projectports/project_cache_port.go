package projectports

import "context"

type ProjectCachePort interface {
	IncrementEditCount(ctx context.Context, projectID string) error

	IsInTop100(ctx context.Context, projectID string) (bool, error)

	CacheProject(ctx context.Context, projectID string, projectData any) error

	GetCachedProject(ctx context.Context, projectID string) (any, error)

	DeleteFromCache(ctx context.Context, projectID string) error

	DeleteFromTop100(ctx context.Context, projectID string) error
}
