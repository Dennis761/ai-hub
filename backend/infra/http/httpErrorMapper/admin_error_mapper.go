package httperrormapper

import "net/http"

// maps admin domain errors to HTTP responses
func MapAdminDomainErrorToHttp(err error) HttpError {
	if err == nil {
		return HttpError{
			Status: http.StatusInternalServerError,
			Body:   map[string]string{"error": "unknown error", "key": "InternalServerError"},
		}
	}

	var key string
	if c, ok := err.(interface{ Code() string }); ok {
		key = c.Code()
	}

	switch key {
	case "AdminNotFound":
		return newHttpError(http.StatusNotFound, err, key)
	case "Forbidden":
		return newHttpError(http.StatusForbidden, err, key)

	case "InvalidCredentials":
		return newHttpError(http.StatusUnauthorized, err, key)
	case "EmailNotVerified":
		return newHttpError(http.StatusForbidden, err, key)

	case "EmailAlreadyUsed",
		"InvalidRole",
		"VerificationCodeExpired",
		"InvalidVerificationCode",
		"ResetCodeNotConfirmed",
		"InvalidAdminID",
		"InvalidEmailFormat",
		"EmailTooLong",
		"NameTooShort",
		"NameTooLong",
		"InvalidPasswordHash",
		"PasswordTooShort",
		"PasswordTooLong",
		"ProjectPasswordTooWeak",
		"RoleRequired",
		"RoleTooLong",
		"InvalidVerificationExpiry":
		return newHttpError(http.StatusBadRequest, err, key)

	default:
		return HttpError{
			Status: http.StatusInternalServerError,
			Body:   map[string]string{"error": err.Error(), "key": "InternalServerError"},
		}
	}
}
