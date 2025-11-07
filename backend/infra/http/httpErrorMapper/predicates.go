package httperrormapper

import "net/http"

// domain predicates
func IsAdminError(err error) bool {
	if err == nil {
		return false
	}

	var key string
	if c, ok := err.(interface{ Code() string }); ok {
		key = c.Code()
	}

	if key == "" {
		return false
	}

	switch key {
	case "AdminNotFound",
		"Forbidden",
		"InvalidCredentials",
		"EmailNotVerified",
		"EmailAlreadyUsed",
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
		"PasswordTooWeak",
		"RoleRequired",
		"RoleTooLong",
		"InvalidVerificationExpiry":
		return true
	default:
		return false
	}
}

func IsApiKeyError(err error) bool {
	if err == nil {
		return false
	}

	var key string
	if c, ok := err.(interface{ Code() string }); ok {
		key = c.Code()
	}

	if key == "" {
		return false
	}

	switch key {
	case "APIKeyNotFound",
		"Forbidden",
		"APIKeyIDRequired",
		"OwnerIDRequired",
		"InvalidAPIKeyName",
		"APIKeyNameTooLong",
		"KeyNameAlreadyUsedByAnotherUser",
		"MaxThreeKeysPerNameExceeded",
		"EnvAlreadyExistsForThisName",
		"InvalidUsageEnv",
		"InvalidStatus",
		"InvalidEncryptedKeyFormat",
		"InvalidPlainAPIKeyValue",
		"ProviderRequired",
		"ProviderTooLong",
		"ModelNameRequired",
		"ModelNameTooLong",
		"ModelNotSupported",
		"ProviderAuthFailed":
		return true
	default:
		return false
	}
}

func IsProjectError(err error) bool {
	if err == nil {
		return false
	}

	var key string
	if c, ok := err.(interface{ Code() string }); ok {
		key = c.Code()
	}

	if key == "" {
		return false
	}

	switch key {
	case "ProjectNotFound",
		"Forbidden",
		"InvalidStatus",
		"ProjectIDRequired",
		"OwnerIDRequired",
		"APIKeyRequired",
		"ProjectPasswordRequired",
		"ProjectNameTooShort",
		"ProjectNameTooLong",
		"ProjectNameAlreadyUsed",
		"ProjectNameAlreadyExists",
		"InvalidProjectName",
		"ProjectPasswordTooShort",
		"ProjectPasswordTooWeak",
		"ProjectPasswordTooLong",
		"ProjectPasswordMismatch",
		"AlreadyProjectMember",
		"InvalidProjectCredentials":
		return true
	default:
		return false
	}
}

func IsPromptError(err error) bool {
	if err == nil {
		return false
	}

	var key string
	if c, ok := err.(interface{ Code() string }); ok {
		key = c.Code()
	}

	if key == "" {
		return false
	}

	switch key {
	case "PromptNotFound",
		"TaskNotFound",
		"MissingParameters",
		"NoPlaceholdersProvided",
		"InvalidExecutionOrder",
		"InvalidVersionNumber",
		"HistoryIndexOutOfRange",
		"SamePromptConsecutive":
		return true
	default:
		return false
	}
}

func IsTaskError(err error) bool {
	if err == nil {
		return false
	}

	var key string
	if c, ok := err.(interface{ Code() string }); ok {
		key = c.Code()
	}

	if key == "" {
		return false
	}

	switch key {
	case "TaskNotFound",
		"InvalidTaskName",
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
		return true
	default:
		return false
	}
}

// wrappers returning HTTP status and body

func MapAdminErrorToHttp(err error) (int, map[string]string) {
	he := MapAdminDomainErrorToHttp(err)
	return he.Status, he.Body
}

func MapApiKeyErrorToHttp(err error) (int, map[string]string) {
	he := MapAPIKeyDomainErrorToHttp(err)
	return he.Status, he.Body
}

func MapProjectErrorToHttp(err error) (int, map[string]string) {
	he := MapProjectDomainErrorToHttp(err)
	return he.Status, he.Body
}

func MapPromptErrorToHttp(err error) (int, map[string]string) {
	he := MapPromptDomainErrorToHttp(err)
	return he.Status, he.Body
}

func MapTaskErrorToHttp(err error) (int, map[string]string) {
	he := MapTaskDomainErrorToHttp(err)
	return he.Status, he.Body
}

// generic fallback for unknown errors
func FallbackInternal(err error) (int, map[string]string) {
	if err == nil {
		return http.StatusInternalServerError, map[string]string{
			"error": "unknown error",
			"key":   "InternalServerError",
		}
	}

	return http.StatusInternalServerError, map[string]string{
		"error": err.Error(),
		"key":   "InternalServerError",
	}
}
