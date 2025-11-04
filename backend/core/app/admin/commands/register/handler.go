package adminapp

import (
	"context"
	"time"

	security "ai_hub.com/app/core/app/admin/shared/secret"
	"ai_hub.com/app/core/domain/admindomain"
	"ai_hub.com/app/core/ports/adminports"
)

type CreateAdminHandler struct {
	readRepo  adminports.AdminReadRepository
	writeRepo adminports.AdminWriteRepository
	hasher    adminports.Hasher
	codeGen   adminports.CodeGenerator
	mailer    adminports.Mailer
	uow       adminports.UnitOfWorkPort
	idGen     adminports.IDGenerator
}

func NewCreateAdminHandler(
	readRepo adminports.AdminReadRepository,
	writeRepo adminports.AdminWriteRepository,
	hasher adminports.Hasher,
	codeGen adminports.CodeGenerator,
	mailer adminports.Mailer,
	uow adminports.UnitOfWorkPort,
	idGen adminports.IDGenerator,
) *CreateAdminHandler {
	return &CreateAdminHandler{
		readRepo:  readRepo,
		writeRepo: writeRepo,
		hasher:    hasher,
		codeGen:   codeGen,
		mailer:    mailer,
		uow:       uow,
		idGen:     idGen,
	}
}

func (h *CreateAdminHandler) CreateAdmin(
	ctx context.Context,
	cmd CreateAdminCommand,
) (*admindomain.Admin, error) {
	// prepare verification
	code := h.codeGen.GenerateVerificationCode()
	ttl := h.codeGen.VerificationTTL()
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	now := time.Now().UTC()
	expiresAt := now.Add(ttl)

	res, err := h.uow.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		// validate password
		plainPwd, err := admindomain.NewPlainPassword(cmd.Password)
		if err != nil {
			return nil, err
		}
		pass := plainPwd.Value()

		// check password strength
		if !security.NewSimpleAdminPasswordPolicy(true).Validate(pass).OK {
			return nil, admindomain.PasswordTooWeak()
		}

		// validate email
		emailVal, err := admindomain.NewAdminEmail(cmd.Email)
		if err != nil {
			return nil, err
		}
		email := emailVal.Value()

		// ensure email unique
		existingAgg, err := h.readRepo.FindByEmail(txCtx, email)
		if err != nil {
			return nil, err
		}
		if existingAgg != nil {
			return nil, admindomain.EmailAlreadyUsed()
		}

		// generate id
		rawID := h.idGen.NewID()
		idVal, err := admindomain.NewAdminID(rawID)
		if err != nil {
			return nil, err
		}

		// normalize name
		nameVal, err := admindomain.NewAdminName(cmd.Name)
		if err != nil {
			return nil, err
		}

		// normalize role
		roleVal, err := admindomain.NewAdminRole(cmd.Role)
		if err != nil {
			return nil, err
		}

		// hash password
		hashed, err := h.hasher.Hash(pass)
		if err != nil {
			return nil, err
		}
		passHash, err := admindomain.NewPasswordHashFromString(hashed)
		if err != nil {
			return nil, err
		}

		// build aggregate
		adminAgg, err := admindomain.NewAdmin(admindomain.NewAdminProps{
			ID:           idVal,
			Email:        emailVal,
			Name:         nameVal,
			PasswordHash: passHash,
			Role:         roleVal,
			CurrentTime:  now,
		})
		if err != nil {
			return nil, err
		}

		// attach verification
		codeVal, err := admindomain.NewVerificationCode(&code)
		if err != nil {
			return nil, err
		}
		expVal, err := admindomain.NewVerificationExpiry(&expiresAt)
		if err != nil {
			return nil, err
		}
		adminAgg.SetVerificationCode(codeVal, expVal)

		// persist
		createdAgg, err := h.writeRepo.Create(txCtx, adminAgg)
		if err != nil {
			return nil, err
		}

		return createdAgg, nil
	})
	if err != nil {
		return nil, err
	}

	adminAgg := res.(*admindomain.Admin)

	// send verification mail
	pr := adminAgg.ToPrimitives()
	if err := h.mailer.SendVerificationCode(ctx, pr.Email, code, expiresAt); err != nil {
		return nil, err
	}

	return adminAgg, nil
}
