package httpErrorMapper

import "net/http"

// MapPromptDomainErrorToHttp maps prompt domain errors to HTTP responses.
func MapPromptDomainErrorToHttp(err error) HttpError {
	if err == nil {
		return HttpError{
			Status: http.StatusInternalServerError,
			Body: map[string]string{
				"error": "unknown error",
				"key":   "InternalServerError",
			},
		}
	}

	key := extractKey(err)

	switch key {
	case "PromptNotFound", "TaskNotFound":
		return newHttpError(http.StatusNotFound, err, key)

	case "MissingParameters",
		"NoPlaceholdersProvided",
		"InvalidExecutionOrder",
		"InvalidVersionNumber",
		"HistoryIndexOutOfRange":
		return newHttpError(http.StatusBadRequest, err, key)

	case "SamePromptConsecutive":
		return newHttpError(http.StatusConflict, err, key)

	default:
		return HttpError{
			Status: http.StatusInternalServerError,
			Body: map[string]string{
				"error": err.Error(),
				"key":   "InternalServerError",
			},
		}
	}
}
