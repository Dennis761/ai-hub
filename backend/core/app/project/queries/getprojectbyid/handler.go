package projectapp

import (
	"context"

	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/core/ports/projectports"
)

type GetProjectByIDHandler struct {
	projectReadRepo projectports.ProjectReadRepository
	projectCache    projectports.ProjectCachePort
}

func NewGetProjectByIDHandler(
	projectReadRepo projectports.ProjectReadRepository,
	projectCache projectports.ProjectCachePort,
) *GetProjectByIDHandler {
	return &GetProjectByIDHandler{
		projectReadRepo: projectReadRepo,
		projectCache:    projectCache,
	}
}

func (h *GetProjectByIDHandler) GetProjectByID(
	ctx context.Context,
	query GetProjectByIDQuery,
) (GetProjectByIDResult, error) {
	// fast path: try cache
	if cached, ok := h.tryCache(ctx, query); ok {
		return cached, nil
	}

	// load aggregate
	projectID, err := projectdomain.NewProjectID(query.ProjectID)
	if err != nil {
		return GetProjectByIDResult{}, err
	}

	project, err := h.projectReadRepo.FindByID(ctx, projectID)
	if err != nil {
		return GetProjectByIDResult{}, err
	}
	if project == nil {
		return GetProjectByIDResult{}, projectdomain.ProjectNotFound()
	}

	// check access
	p := project.ToPrimitives()
	if err := ensureProjectAccess(p.OwnerID, p.AdminAccess, query.AdminID); err != nil {
		return GetProjectByIDResult{}, err
	}

	// return safe dto
	return buildProjectResult(p), nil
}

func (h *GetProjectByIDHandler) tryCache(
	ctx context.Context,
	query GetProjectByIDQuery,
) (GetProjectByIDResult, bool) {
	if h.projectCache == nil {
		return GetProjectByIDResult{}, false
	}

	cached, err := h.projectCache.GetCachedProject(ctx, query.ProjectID)
	if err != nil || cached == nil {
		return GetProjectByIDResult{}, false
	}

	var prim projectdomain.ProjectPrimitives
	switch v := cached.(type) {
	case projectdomain.ProjectPrimitives:
		prim = v
	case *projectdomain.ProjectPrimitives:
		if v == nil {
			return GetProjectByIDResult{}, false
		}
		prim = *v
	default:
		return GetProjectByIDResult{}, false
	}

	// still check ACL even for cache
	if err := ensureProjectAccess(prim.OwnerID, prim.AdminAccess, query.AdminID); err != nil {
		return GetProjectByIDResult{}, false
	}

	return buildProjectResult(prim), true
}

// require owner or admin
func ensureProjectAccess(ownerID string, adminAccess []string, adminID string) error {
	if ownerID == adminID {
		return nil
	}
	for _, id := range adminAccess {
		if id == adminID {
			return nil
		}
	}
	return projectdomain.Forbidden()
}

// project → safe dto
func buildProjectResult(p projectdomain.ProjectPrimitives) GetProjectByIDResult {
	return GetProjectByIDResult{
		ID:        p.ID,
		Name:      p.Name,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
