package adminapp

import (
	"context"
	"time"

	"ai_hub.com/app/core/domain/admindomain"
	"ai_hub.com/app/core/ports/adminports"
)

type StartPasswordResetHandler struct {
	readRepo  adminports.AdminReadRepository
	writeRepo adminports.AdminWriteRepository
	codeGen   adminports.CodeGenerator
	mailer    adminports.Mailer
	uow       adminports.UnitOfWorkPort
}

func NewStartPasswordResetHandler(
	readRepo adminports.AdminReadRepository,
	writeRepo adminports.AdminWriteRepository,
	codeGen adminports.CodeGenerator,
	mailer adminports.Mailer,
	uow adminports.UnitOfWorkPort,
) *StartPasswordResetHandler {
	return &StartPasswordResetHandler{
		readRepo:  readRepo,
		writeRepo: writeRepo,
		codeGen:   codeGen,
		mailer:    mailer,
		uow:       uow,
	}
}

// Start password reset flow (silent on unknown email).
func (h *StartPasswordResetHandler) StartPasswordReset(
	ctx context.Context,
	cmd StartPasswordResetCommand,
) error {
	// normalize email
	emailVal, err := admindomain.NewAdminEmail(cmd.Email)
	if err != nil {
		return err
	}
	email := emailVal.Value()

	// prepare code
	code := h.codeGen.GenerateVerificationCode()
	ttl := h.codeGen.VerificationTTL()
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	now := time.Now().UTC()
	expiresAt := now.Add(ttl)

	// prepare VOs
	codeVal, err := admindomain.NewVerificationCode(&code)
	if err != nil {
		return err
	}
	expVal, err := admindomain.NewVerificationExpiry(&expiresAt)
	if err != nil {
		return err
	}

	// persist reset code
	_, err = h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// load admin (silent)
		adminAgg, err := h.readRepo.FindByEmail(txCtx, email)
		if err != nil {
			return nil, err
		}
		if adminAgg == nil {
			return nil, admindomain.AdminNotFound()
		}

		// update aggregate
		adminAgg.SetVerificationCode(codeVal, expVal)
		adminAgg.ResetConfirmation()

		// save
		if _, err := h.writeRepo.Update(txCtx, adminAgg); err != nil {
			return nil, err
		}

		return adminAgg, nil
	})
	if err != nil {
		return err
	}

	// send mail (best-effort)
	if err := h.mailer.SendResetCode(ctx, email, code, expiresAt); err != nil {
		return err
	}

	return nil
}
