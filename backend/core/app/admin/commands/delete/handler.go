package adminapp

import (
	"context"

	"ai_hub.com/app/core/domain/admindomain"
	"ai_hub.com/app/core/ports/adminports"
)

// handles admin deletion logic
type DeleteAdminHandler struct {
	adminReader adminports.AdminReadRepository
	adminWriter adminports.AdminWriteRepository
	uow         adminports.UnitOfWorkPort
}

func NewDeleteAdminHandler(
	adminReader adminports.AdminReadRepository,
	adminWriter adminports.AdminWriteRepository,
	uow adminports.UnitOfWorkPort,
) *DeleteAdminHandler {
	return &DeleteAdminHandler{
		adminReader: adminReader,
		adminWriter: adminWriter,
		uow:         uow,
	}
}

func (h *DeleteAdminHandler) DeleteAdmin(ctx context.Context, cmd DeleteAdminCommand) error {
	// validate and normalize admin ID
	adminIDVO, err := admindomain.NewAdminID(cmd.ID)
	if err != nil {
		return err
	}
	adminID := adminIDVO.Value()

	// execute deletion in transaction
	_, err = h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// check if aggregate exists
		agg, err := h.adminReader.FindByID(ctx, adminID)
		if err != nil {
			return nil, err
		}
		if agg == nil {
			return nil, admindomain.AdminNotFound()
		}

		// remove aggregate
		if err := h.adminWriter.Delete(ctx, adminID); err != nil {
			return nil, err
		}
		return nil, nil
	})

	return err
}
