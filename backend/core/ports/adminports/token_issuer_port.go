package adminports

type TokenPayload struct {
	UserID string
	Email  *string
	Roles  []string
}

type TokenOptions struct {
	ExpiresIn string
}

type TokenIssuer interface {
	Issue(payload TokenPayload, opts *TokenOptions) (string, error)
}
