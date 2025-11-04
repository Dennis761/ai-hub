package projectapp

import (
	"context"

	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/core/ports/projectports"
)

type DeleteProjectHandler struct {
	readRepo  projectports.ProjectReadRepository
	writeRepo projectports.ProjectWriteRepository
	uow       projectports.UnitOfWorkPort
	cachePort projectports.ProjectCachePort
}

func NewDeleteProjectHandler(
	readRepo projectports.ProjectReadRepository,
	writeRepo projectports.ProjectWriteRepository,
	uow projectports.UnitOfWorkPort,
	cachePort projectports.ProjectCachePort,
) *DeleteProjectHandler {
	return &DeleteProjectHandler{
		readRepo:  readRepo,
		writeRepo: writeRepo,
		uow:       uow,
		cachePort: cachePort,
	}
}

func (h *DeleteProjectHandler) Delete(
	ctx context.Context,
	cmd DeleteProjectCommand,
) error {
	var deletedID string

	_, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// normalize project id
		projectID, err := projectdomain.NewProjectID(cmd.ID)
		if err != nil {
			return nil, err
		}

		// load aggregate
		project, err := h.readRepo.FindByID(txCtx, projectID)
		if err != nil {
			return nil, err
		}
		if project == nil {
			return nil, nil
		}

		// normalize admin id
		adminID, err := projectdomain.NewAccessAdminID(cmd.AdminID)
		if err != nil {
			return nil, err
		}

		// check access (owner or admin)
		if !project.HasAdminAccess(adminID) {
			return nil, projectdomain.Forbidden()
		}

		// delete aggregate
		if err := h.writeRepo.Delete(txCtx, projectID); err != nil {
			return nil, err
		}

		deletedID = projectID.Value()
		return nil, nil
	})
	if err != nil {
		return err
	}

	// best-effort cache cleanup (after commit)
	if deletedID != "" && h.cachePort != nil {
		_ = h.cachePort.DeleteFromCache(ctx, deletedID)
		_ = h.cachePort.DeleteFromTop100(ctx, deletedID)
	}

	return nil
}
