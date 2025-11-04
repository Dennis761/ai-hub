package projectapp

import (
	"context"
	"time"

	projectsecret "ai_hub.com/app/core/app/project/shared/secret"
	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/core/ports/projectports"
)

type CreateProjectHandler struct {
	readRepo  projectports.ProjectReadRepository
	writeRepo projectports.ProjectWriteRepository
	uow       projectports.UnitOfWorkPort
	idGen     projectports.IDGenerator
	hasher    projectports.Hasher
}

func NewCreateProjectHandler(
	readRepo projectports.ProjectReadRepository,
	writeRepo projectports.ProjectWriteRepository,
	uow projectports.UnitOfWorkPort,
	idGen projectports.IDGenerator,
	hasher projectports.Hasher,
) *CreateProjectHandler {
	return &CreateProjectHandler{
		readRepo:  readRepo,
		writeRepo: writeRepo,
		uow:       uow,
		idGen:     idGen,
		hasher:    hasher,
	}
}

func (h *CreateProjectHandler) Create(
	ctx context.Context,
	cmd CreateProjectCommand,
) (*projectdomain.Project, error) {
	res, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// normalize name
		name, err := projectdomain.NewProjectName(cmd.Name)
		if err != nil {
			return nil, err
		}

		// ensure unique name
		existing, err := h.readRepo.FindByName(txCtx, name)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, projectdomain.ProjectNameAlreadyExists()
		}

		// generate id
		rawID := h.idGen.NewID()
		projectID, err := projectdomain.NewProjectID(rawID)
		if err != nil {
			return nil, err
		}

		// normalize owner
		ownerID, err := projectdomain.NewOwnerID(cmd.OwnerID)
		if err != nil {
			return nil, err
		}

		// normalize secret
		plainSecret, err := projectdomain.NewPlainProjectAPIKey(cmd.APIKey)
		if err != nil {
			return nil, err
		}

		// check secret strength
		policy := projectsecret.NewSimpleProjectPasswordPolicy(true)
		if !policy.Validate(plainSecret.Value()).OK {
			return nil, projectdomain.ProjectPasswordTooWeak()
		}

		// hash secret
		hashed, err := h.hasher.Hash(plainSecret.Value())
		if err != nil {
			return nil, err
		}
		hashedKey, err := projectdomain.NewHashedAPIKey(hashed)
		if err != nil {
			return nil, err
		}

		// map admin access
		var admins []projectdomain.AccessAdminID
		if len(cmd.AdminAccess) > 0 {
			admins = make([]projectdomain.AccessAdminID, 0, len(cmd.AdminAccess))
			for _, raw := range cmd.AdminAccess {
				adminID, err := projectdomain.NewAccessAdminID(raw)
				if err != nil {
					return nil, err
				}
				admins = append(admins, adminID)
			}
		}

		now := time.Now().UTC()
		props := projectdomain.CreateProps{
			ID:          projectID,
			Name:        name,
			APIKey:      hashedKey,
			OwnerID:     ownerID,
			AdminAccess: admins,
			Now:         &now,
		}

		// optional status
		if cmd.Status != nil && *cmd.Status != "" {
			props.Status = cmd.Status
		}

		// build aggregate
		project, err := projectdomain.Create(props)
		if err != nil {
			return nil, err
		}

		// persist
		created, err := h.writeRepo.Create(txCtx, &project)
		if err != nil {
			return nil, err
		}

		return created, nil
	})
	if err != nil {
		return nil, err
	}

	return res.(*projectdomain.Project), nil
}
