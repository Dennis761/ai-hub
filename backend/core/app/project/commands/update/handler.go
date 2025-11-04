package projectapp

import (
	"context"

	projectsecret "ai_hub.com/app/core/app/project/shared/secret"
	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/core/ports/projectports"
)

type UpdateProjectHandler struct {
	readRepo  projectports.ProjectReadRepository
	writeRepo projectports.ProjectWriteRepository
	uow       projectports.UnitOfWorkPort
	cache     projectports.ProjectCachePort
	hasher    projectports.Hasher
}

func NewUpdateProjectHandler(
	readRepo projectports.ProjectReadRepository,
	writeRepo projectports.ProjectWriteRepository,
	uow projectports.UnitOfWorkPort,
	cache projectports.ProjectCachePort,
	hasher projectports.Hasher,
) *UpdateProjectHandler {
	return &UpdateProjectHandler{
		readRepo:  readRepo,
		writeRepo: writeRepo,
		uow:       uow,
		cache:     cache,
		hasher:    hasher,
	}
}

func (h *UpdateProjectHandler) Update(
	ctx context.Context,
	cmd UpdateProjectCommand,
) (*projectdomain.Project, error) {
	res, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// normalize id
		projectID, err := projectdomain.NewProjectID(cmd.ID)
		if err != nil {
			return nil, err
		}

		// load project
		project, err := h.readRepo.FindByID(txCtx, projectID)
		if err != nil {
			return nil, err
		}
		if project == nil {
			return nil, projectdomain.ProjectNotFound()
		}

		// check owner
		ownerID, err := projectdomain.NewOwnerID(cmd.OwnerID)
		if err != nil {
			return nil, err
		}
		if project.OwnerID().Value() != ownerID.Value() {
			return nil, projectdomain.Forbidden()
		}

		current := project.ToPrimitives()

		nameChanged := false
		apiKeyChanged := false
		statusChanged := false

		// update name
		if cmd.Name != nil {
			newName, err := projectdomain.NewProjectName(*cmd.Name)
			if err != nil {
				return nil, err
			}
			if current.Name != newName.Value() {
				// ensure unique name
				existing, err := h.readRepo.FindByName(txCtx, newName)
				if err != nil {
					return nil, err
				}
				if existing != nil && existing.ToPrimitives().ID != current.ID {
					return nil, projectdomain.ProjectNameAlreadyExists()
				}

				project.Rename(newName)
				nameChanged = true
			}
		}

		// update secret
		if cmd.APIKey != nil {
			plainSecretVO, err := projectdomain.NewPlainProjectAPIKey(*cmd.APIKey)
			if err != nil {
				return nil, err
			}
			plain := plainSecretVO.Value()

			pol := projectsecret.NewSimpleProjectPasswordPolicy(true).Validate(plain)
			if !pol.OK {
				return nil, projectdomain.ProjectPasswordTooWeak()
			}

			hashed, err := h.hasher.Hash(plain)
			if err != nil {
				return nil, err
			}
			newHashedVO, err := projectdomain.NewHashedAPIKey(hashed)
			if err != nil {
				return nil, err
			}

			if current.APIKey != newHashedVO.Value() {
				project.RebindAPIKey(newHashedVO)
				apiKeyChanged = true
			}
		}

		// update status
		if cmd.Status != nil {
			if err := project.SetStatus(*cmd.Status); err != nil {
				return nil, err
			}
			if current.Status != *cmd.Status {
				statusChanged = true
			}
		}

		// nothing to persist
		if !nameChanged && !apiKeyChanged && !statusChanged {
			return project, nil
		}

		// persist
		updated, err := h.writeRepo.Update(txCtx, project)
		if err != nil {
			return nil, err
		}

		return updated, nil
	})
	if err != nil {
		return nil, err
	}

	project := res.(*projectdomain.Project)
	p := project.ToPrimitives()

	// best-effort cache update
	if h.cache != nil {
		_ = h.cache.IncrementEditCount(ctx, p.ID)

		inTop, cacheErr := h.cache.IsInTop100(ctx, p.ID)
		if cacheErr == nil && inTop {
			_ = h.cache.CacheProject(ctx, p.ID, p)
		}
	}

	return project, nil
}
