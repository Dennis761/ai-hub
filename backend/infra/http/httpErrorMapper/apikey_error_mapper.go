package httpErrorMapper

import "net/http"

// maps API key domain errors to HTTP responses
func MapAPIKeyDomainErrorToHttp(err error) HttpError {
	if err == nil {
		return HttpError{
			Status: http.StatusInternalServerError,
			Body:   map[string]string{"error": "unknown error", "key": "InternalServerError"},
		}
	}

	key := extractKey(err)

	switch key {
	case "APIKeyNotFound":
		return newHttpError(http.StatusNotFound, err, key)
	case "Forbidden":
		return newHttpError(http.StatusForbidden, err, key)

	case "InvalidAPIKeyName",
		"APIKeyNameTooLong",
		"InvalidEncryptedKeyFormat",
		"InvalidStatus",
		"InvalidUsageEnv",
		"ProviderRequired",
		"ProviderTooLong",
		"ModelNameRequired",
		"ModelNameTooLong",
		"APIKeyIDRequired",
		"OwnerIDRequired",
		"KeyNameAlreadyUsedByAnotherUser",
		"MaxThreeKeysPerNameExceeded",
		"EnvAlreadyExistsForThisName",
		"InvalidPlainAPIKeyValue":
		return newHttpError(http.StatusBadRequest, err, key)

	case "ModelNotSupported", "ProviderAuthFailed":
		return newHttpError(http.StatusBadRequest, err, key)

	default:
		return HttpError{
			Status: http.StatusInternalServerError,
			Body:   map[string]string{"error": err.Error(), "key": "InternalServerError"},
		}
	}
}
