package adminapp

import (
	"context"

	"ai_hub.com/app/core/domain/admindomain"
	"ai_hub.com/app/core/ports/adminports"
)

type LoginAdminCommand struct {
	Email    string
	Password string
}

type LoginAdminHandler struct {
	repo   adminports.AdminReadRepository
	hasher adminports.Hasher
	issuer adminports.TokenIssuer
}

func NewLoginAdminHandler(
	repo adminports.AdminReadRepository,
	hasher adminports.Hasher,
	issuer adminports.TokenIssuer,
) *LoginAdminHandler {
	return &LoginAdminHandler{
		repo:   repo,
		hasher: hasher,
		issuer: issuer,
	}
}

func (h *LoginAdminHandler) Login(
	ctx context.Context,
	cmd LoginAdminCommand,
) (*LoginAdminResponse, error) {
	// validate input
	if cmd.Email == "" || cmd.Password == "" {
		return nil, admindomain.InvalidCredentials()
	}

	// normalize email
	emailVal, err := admindomain.NewAdminEmail(cmd.Email)
	if err != nil {
		return nil, admindomain.InvalidCredentials()
	}
	email := emailVal.Value()

	// normalize password
	plainPwd, err := admindomain.NewPlainPassword(cmd.Password)
	if err != nil {
		return nil, admindomain.InvalidCredentials()
	}
	pass := plainPwd.Value()

	// load admin aggregate
	admin, err := h.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, admindomain.InvalidCredentials()
	}

	// check password
	storedHash := admin.ExposePasswordHash()
	if storedHash == "" {
		return nil, admindomain.InvalidCredentials()
	}
	ok, err := h.hasher.Compare(pass, storedHash)
	if err != nil || !ok {
		return nil, admindomain.InvalidCredentials()
	}

	// check verification
	if !admin.IsVerified().Value() {
		return nil, admindomain.EmailNotVerified()
	}

	// issue token
	data := admin.ToPrimitives()
	payload := adminports.TokenPayload{
		UserID: data.ID,
	}
	if data.Email != "" {
		e := data.Email
		payload.Email = &e
	}
	if data.Role != "" {
		payload.Roles = []string{data.Role}
	}

	token, err := h.issuer.Issue(payload, nil)
	if err != nil {
		return nil, err
	}

	// build response
	var emailPtr *string
	if data.Email != "" {
		e := data.Email
		emailPtr = &e
	}
	var rolePtr *string
	if data.Role != "" {
		r := data.Role
		rolePtr = &r
	}

	return &LoginAdminResponse{
		Token: token,
		Admin: LoginAdminResult{
			ID:         data.ID,
			Email:      emailPtr,
			Name:       data.Name,
			Role:       rolePtr,
			IsVerified: data.IsVerified,
			CreatedAt:  data.CreatedAt,
			UpdatedAt:  data.UpdatedAt,
		},
	}, nil
}
