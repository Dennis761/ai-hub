package adminapp

import (
	"context"

	"ai_hub.com/app/core/domain/admindomain"
	"ai_hub.com/app/core/ports/adminports"
)

type RenameAdminHandler struct {
	readRepo  adminports.AdminReadRepository
	writeRepo adminports.AdminWriteRepository
	uow       adminports.UnitOfWorkPort
}

func NewRenameAdminHandler(
	readRepo adminports.AdminReadRepository,
	writeRepo adminports.AdminWriteRepository,
	uow adminports.UnitOfWorkPort,
) *RenameAdminHandler {
	return &RenameAdminHandler{
		readRepo:  readRepo,
		writeRepo: writeRepo,
		uow:       uow,
	}
}

func (h *RenameAdminHandler) RenameAdmin(
	ctx context.Context,
	cmd RenameAdminCommand,
) (*admindomain.Admin, error) {
	res, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// validate id
		idVal, err := admindomain.NewAdminID(cmd.ID)
		if err != nil {
			return nil, err
		}

		// load aggregate
		adminAgg, err := h.readRepo.FindByID(txCtx, idVal.Value())
		if err != nil {
			return nil, err
		}
		if adminAgg == nil {
			return nil, admindomain.AdminNotFound()
		}

		// normalize name
		nameRaw := ""
		if cmd.Name != nil {
			nameRaw = *cmd.Name
		}
		nameVal, err := admindomain.NewAdminName(nameRaw)
		if err != nil {
			return nil, err
		}

		// mutate
		adminAgg.Rename(nameVal)

		// persist
		updatedAgg, err := h.writeRepo.Update(txCtx, adminAgg)
		if err != nil {
			return nil, err
		}
		return updatedAgg, nil
	})
	if err != nil {
		return nil, err
	}

	return res.(*admindomain.Admin), nil
}
