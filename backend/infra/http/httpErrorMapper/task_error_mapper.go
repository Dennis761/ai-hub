package httperrormapper

import "net/http"

// MapTaskDomainErrorToHttp maps task domain errors to HTTP responses.
func MapTaskDomainErrorToHttp(err error) HttpError {
	if err == nil {
		return HttpError{
			Status: http.StatusInternalServerError,
			Body: map[string]string{
				"error": "unknown error",
				"key":   "InternalServerError",
			},
		}
	}

	var key string
	if c, ok := err.(interface{ Code() string }); ok {
		key = c.Code()
	}

	switch key {
	case "TaskNotFound":
		return newHttpError(http.StatusNotFound, err, key)

	case "InvalidTaskName",
		"TaskNameTooShort",
		"TaskNameTooLong",
		"TaskDescriptionTooShort",
		"TaskDescriptionTooLong",
		"InvalidAPIEndpoint",
		"APIMethodRequired",
		"APIMethodTooLong",
		"InvalidStatus",
		"TaskStatusRequired",
		"TaskIDRequired",
		"CreatorIDRequired":
		return newHttpError(http.StatusBadRequest, err, key)

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
