package projectapp

import (
	"context"

	"ai_hub.com/app/core/domain/projectdomain"
	"ai_hub.com/app/core/ports/projectports"
)

type JoinProjectByNameHandler struct {
	readRepo  projectports.ProjectReadRepository
	writeRepo projectports.ProjectWriteRepository
	uow       projectports.UnitOfWorkPort
	hasher    projectports.Hasher
}

func NewJoinProjectByNameHandler(
	readRepo projectports.ProjectReadRepository,
	writeRepo projectports.ProjectWriteRepository,
	uow projectports.UnitOfWorkPort,
	hasher projectports.Hasher,
) *JoinProjectByNameHandler {
	return &JoinProjectByNameHandler{
		readRepo:  readRepo,
		writeRepo: writeRepo,
		uow:       uow,
		hasher:    hasher,
	}
}

func (h *JoinProjectByNameHandler) JoinProjectByName(
	ctx context.Context,
	cmd JoinProjectByNameCommand,
) (*projectdomain.Project, error) {
	res, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// normalize input
		adminID, err := projectdomain.NewAccessAdminID(cmd.AdminID)
		if err != nil {
			return nil, err
		}
		projectName, err := projectdomain.NewProjectName(cmd.Name)
		if err != nil {
			return nil, err
		}
		secretVO, err := projectdomain.NewPlainProjectAPIKey(cmd.APIKey)
		if err != nil {
			return nil, err
		}
		plainSecret := secretVO.Value()

		// load project by name
		project, err := h.readRepo.FindByName(txCtx, projectName)
		if err != nil {
			return nil, err
		}
		if project == nil {
			return nil, projectdomain.InvalidProjectCredentials()
		}

		// verify secret (stored as hash)
		storedHash := project.APIKey().Value()
		match, err := h.hasher.Compare(plainSecret, storedHash)
		if err != nil {
			return nil, err
		}
		if !match {
			return nil, projectdomain.InvalidProjectCredentials()
		}

		// already in access list
		if project.HasAdminAccess(adminID) {
			return nil, projectdomain.AlreadyProjectMember()
		}

		// grant access
		project.GrantAdminAccess(adminID)

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

	return res.(*projectdomain.Project), nil
}
