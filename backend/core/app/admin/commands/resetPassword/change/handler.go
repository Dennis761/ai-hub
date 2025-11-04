package adminapp

import (
	"context"

	security "ai_hub.com/app/core/app/admin/shared/secret"
	"ai_hub.com/app/core/domain/admindomain"
	"ai_hub.com/app/core/ports/adminports"
)

type ChangePasswordWithCodeHandler struct {
	readRepo  adminports.AdminReadRepository
	writeRepo adminports.AdminWriteRepository
	hasher    adminports.Hasher
	uow       adminports.UnitOfWorkPort
}

func NewChangePasswordWithCodeHandler(
	readRepo adminports.AdminReadRepository,
	writeRepo adminports.AdminWriteRepository,
	hasher adminports.Hasher,
	uow adminports.UnitOfWorkPort,
) *ChangePasswordWithCodeHandler {
	return &ChangePasswordWithCodeHandler{
		readRepo:  readRepo,
		writeRepo: writeRepo,
		hasher:    hasher,
		uow:       uow,
	}
}

func (h *ChangePasswordWithCodeHandler) ChangePasswordWithCode(
	ctx context.Context,
	cmd ChangePasswordWithCodeCommand,
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

		// check reset confirmation
		if !adminAgg.IsResetCodeConfirmed().Value() {
			return nil, admindomain.ResetCodeNotConfirmed()
		}

		// normalize password
		plainPwd, err := admindomain.NewPlainPassword(cmd.NewPassword)
		if err != nil {
			return nil, err
		}
		pass := plainPwd.Value()

		// check strength
		if !security.NewSimpleAdminPasswordPolicy(true).Validate(pass).OK {
			return nil, admindomain.PasswordTooWeak()
		}

		// hash password
		hashStr, err := h.hasher.Hash(pass)
		if err != nil {
			return nil, err
		}
		hashVal, err := admindomain.NewPasswordHashFromString(hashStr)
		if err != nil {
			return nil, err
		}

		// mutate aggregate
		adminAgg.ChangePassword(hashVal)
		adminAgg.ResetConfirmation()
		adminAgg.ClearVerificationCode()

		// persist
		if _, err := h.writeRepo.Update(txCtx, adminAgg); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}
