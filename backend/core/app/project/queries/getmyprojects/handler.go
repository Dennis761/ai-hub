package projectapp

import (
	"context"

	"ai_hub.com/app/core/ports/projectports"
)

type GetMyProjectsHandler struct {
	projectReadRepo projectports.ProjectReadRepository
}

func NewGetMyProjectsHandler(
	projectReadRepo projectports.ProjectReadRepository,
) *GetMyProjectsHandler {
	return &GetMyProjectsHandler{
		projectReadRepo: projectReadRepo,
	}
}

func (h *GetMyProjectsHandler) GetMyProjects(
	ctx context.Context,
	query GetMyProjectsQuery,
) (GetMyProjectsResult, error) {
	// load all projects where admin has access
	projects, err := h.projectReadRepo.FindAllAccessibleByAdmin(ctx, query.AdminID)
	if err != nil {
		return GetMyProjectsResult{}, err
	}

	// build safe dto
	items := make([]ProjectListItem, 0, len(projects))
	for _, projectAgg := range projects {
		p := projectAgg.ToPrimitives()
		items = append(items, ProjectListItem{
			ID:     p.ID,
			Name:   p.Name,
			Status: p.Status,
		})
	}

	return GetMyProjectsResult{
		Projects: items,
	}, nil
}
