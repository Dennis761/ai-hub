package adminapp

import (
	"context"
	"time"

	"ai_hub.com/app/core/domain/admindomain"
	"ai_hub.com/app/core/ports/adminports"
)

type ConfirmResetCodeHandler struct {
	readRepo  adminports.AdminReadRepository
	writeRepo adminports.AdminWriteRepository
	uow       adminports.UnitOfWorkPort
}

func NewConfirmResetCodeHandler(
	readRepo adminports.AdminReadRepository,
	writeRepo adminports.AdminWriteRepository,
	uow adminports.UnitOfWorkPort,
) *ConfirmResetCodeHandler {
	return &ConfirmResetCodeHandler{
		readRepo:  readRepo,
		writeRepo: writeRepo,
		uow:       uow,
	}
}

func (h *ConfirmResetCodeHandler) ConfirmResetCode(
	ctx context.Context,
	cmd ConfirmResetCodeCommand,
) error {
	_, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// normalize email
		emailVal, err := admindomain.NewAdminEmail(cmd.Email)
		if err != nil {
			return nil, admindomain.AdminNotFound()
		}
		email := emailVal.Value()

		// load admin
		adminAgg, err := h.readRepo.FindByEmail(txCtx, email)
		if err != nil {
			return nil, err
		}
		if adminAgg == nil {
			return nil, admindomain.AdminNotFound()
		}

		// get stored code and expiry
		storedCode := adminAgg.VerificationCode().Value()
		expiry := adminAgg.VerificationCodeExpiry()
		expiryTime := expiry.Value()
		if storedCode == nil || expiryTime == nil {
			return nil, admindomain.InvalidVerificationCode()
		}

		// check expiry
		now := time.Now().UTC()
		if expiry.IsExpired(now) {
			return nil, admindomain.VerificationCodeExpired()
		}

		// normalize incoming code
		incomingCode, err := admindomain.NewVerificationCode(&cmd.Code)
		if err != nil {
			return nil, admindomain.InvalidVerificationCode()
		}
		inCode := incomingCode.Value()
		if inCode == nil || *inCode == "" {
			return nil, admindomain.InvalidVerificationCode()
		}

		// compare codes
		if *storedCode != *inCode {
			return nil, admindomain.InvalidVerificationCode()
		}

		// mark as confirmed
		adminAgg.ConfirmResetCode()

		// persist
		if _, err := h.writeRepo.Update(txCtx, adminAgg); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}
