package httpErrorMapper

import "net/http"

// MapProjectDomainErrorToHttp maps project domain errors to HTTP responses.
func MapProjectDomainErrorToHttp(err error) HttpError {
	if err == nil {
		return HttpError{
			Status: http.StatusInternalServerError,
			Body:   map[string]string{"error": "unknown error", "key": "InternalServerError"},
		}
	}

	key := extractKey(err)

	switch key {
	case "ProjectNotFound":
		return newHttpError(http.StatusNotFound, err, key)
	case "Forbidden":
		return newHttpError(http.StatusForbidden, err, key)

	case "InvalidProjectCredentials":
		return newHttpError(http.StatusUnauthorized, err, key)

	case "ProjectIDRequired",
		"OwnerIDRequired",
		"APIKeyRequired",
		"InvalidStatus",
		"ProjectNameTooShort",
		"ProjectNameTooLong",
		"InvalidProjectName",
		"ProjectPasswordRequired",
		"ProjectPasswordTooShort",
		"ProjectPasswordTooWeak",
		"ProjectPasswordTooLong",
		"ProjectPasswordMismatch":
		return newHttpError(http.StatusBadRequest, err, key)

	case "ProjectNameAlreadyUsed",
		"ProjectNameAlreadyExists",
		"AlreadyProjectMember":
		return newHttpError(http.StatusConflict, err, key)

	default:
		return HttpError{
			Status: http.StatusInternalServerError,
			Body:   map[string]string{"error": err.Error(), "key": "InternalServerError"},
		}
	}
}
