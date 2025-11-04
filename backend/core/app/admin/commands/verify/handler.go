package adminapp

import (
	"context"
	"time"

	"ai_hub.com/app/core/domain/admindomain"
	"ai_hub.com/app/core/ports/adminports"
)

type VerifyAdminHandler struct {
	readRepo  adminports.AdminReadRepository
	writeRepo adminports.AdminWriteRepository
	uow       adminports.UnitOfWorkPort
}

func NewVerifyAdminHandler(
	readRepo adminports.AdminReadRepository,
	writeRepo adminports.AdminWriteRepository,
	uow adminports.UnitOfWorkPort,
) *VerifyAdminHandler {
	return &VerifyAdminHandler{
		readRepo:  readRepo,
		writeRepo: writeRepo,
		uow:       uow,
	}
}

func (h *VerifyAdminHandler) Verify(
	ctx context.Context,
	cmd VerifyAdminCommand,
) error {
	// normalize id
	idVal, err := admindomain.NewAdminID(cmd.ID)
	if err != nil {
		return err
	}

	// normalize code
	codeVal, err := admindomain.NewVerificationCode(&cmd.Code)
	if err != nil {
		return admindomain.InvalidVerificationCode()
	}
	inCode := codeVal.Value()
	if inCode == nil || *inCode == "" {
		return admindomain.InvalidVerificationCode()
	}

	// load aggregate
	adminAgg, err := h.readRepo.FindByID(ctx, idVal.Value())
	if err != nil {
		return err
	}
	if adminAgg == nil {
		return admindomain.AdminNotFound()
	}

	// get stored code & expiry
	storedCode := adminAgg.VerificationCode().Value()
	storedExpiry := adminAgg.VerificationCodeExpiry().Value()
	if storedCode == nil || storedExpiry == nil {
		return admindomain.InvalidVerificationCode()
	}

	// check expiry
	now := time.Now().UTC()
	if now.After(*storedExpiry) {
		return admindomain.VerificationCodeExpired()
	}

	// compare codes
	if *storedCode != *inCode {
		return admindomain.InvalidVerificationCode()
	}

	// mutate
	adminAgg.Verify()
	adminAgg.ClearVerificationCode()

	// persist
	_, err = h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		updatedAgg, err := h.writeRepo.Update(txCtx, adminAgg)
		if err != nil {
			return nil, err
		}
		return updatedAgg, nil
	})

	return err
}
