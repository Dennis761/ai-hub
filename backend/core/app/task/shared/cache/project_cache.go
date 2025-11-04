package cache

import (
	"context"

	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/core/ports/projectports"
)

type ProjectCacheDeps struct {
	ProjectReadRepo projectports.ProjectReadRepository
	ProjectCache    projectports.ProjectCachePort
}

// UpdateProjectCache updates project cache and increments edit counter (best-effort).
func UpdateProjectCache(ctx context.Context, deps ProjectCacheDeps, projectID string) error {
	if deps.ProjectCache == nil {
		return nil
	}

	// increment project edit counter
	_ = deps.ProjectCache.IncrementEditCount(ctx, projectID)

	// refresh project snapshot in cache
	projectIDVO, err := projectdomain.NewProjectID(projectID)
	if err != nil {
		return err
	}

	project, err := deps.ProjectReadRepo.FindByID(ctx, projectIDVO)
	if err != nil || project == nil {
		return err
	}

	return deps.ProjectCache.CacheProject(ctx, projectID, project.ToPrimitives())
}
