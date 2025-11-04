package apikeydomain

import "fmt"

// DomainError is a common contract for API key domain errors.
type DomainError interface {
	error
	Code() string
}

// basic error types

type simpleError struct {
	code string
	msg  string
}

func (e simpleError) Error() string { return e.msg }
func (e simpleError) Code() string  { return e.code }

type limitError struct {
	code   string
	format string
	limit  int
	def    int
}

func (e limitError) Error() string {
	lim := e.limit
	if lim <= 0 {
		lim = e.def
	}
	return fmt.Sprintf(e.format, lim)
}

func (e limitError) Code() string { return e.code }

type dynamicError struct {
	code string
	msg  string
}

func (e dynamicError) Error() string { return e.msg }
func (e dynamicError) Code() string  { return e.code }

// resource / access

func APIKeyNotFound() DomainError {
	return simpleError{
		code: "APIKeyNotFound",
		msg:  "API key not found",
	}
}

func Forbidden() DomainError {
	return simpleError{
		code: "Forbidden",
		msg:  "Forbidden",
	}
}

// name / quota / collisions

func InvalidAPIKeyName(msgOverride string) DomainError {
	if msgOverride != "" {
		return simpleError{
			code: "InvalidAPIKeyName",
			msg:  msgOverride,
		}
	}
	return simpleError{
		code: "InvalidAPIKeyName",
		msg:  "Invalid API key name (min length: 6)",
	}
}

func APIKeyNameTooLong(max int) DomainError {
	return limitError{
		code:   "APIKeyNameTooLong",
		format: "API key name is too long (max %d)",
		limit:  max,
		def:    100,
	}
}

func KeyNameAlreadyUsedByAnotherUser(name string) DomainError {
	return dynamicError{
		code: "KeyNameAlreadyUsedByAnotherUser",
		msg:  fmt.Sprintf("The key name %q is already used by another user", name),
	}
}

func MaxThreeKeysPerNameExceeded(name string) DomainError {
	return dynamicError{
		code: "MaxThreeKeysPerNameExceeded",
		msg:  fmt.Sprintf("Maximum of 3 keys with name %q allowed", name),
	}
}

func EnvAlreadyExistsForThisName(name, env string) DomainError {
	return dynamicError{
		code: "EnvAlreadyExistsForThisName",
		msg:  fmt.Sprintf("A key %q already exists in %q environment", name, env),
	}
}

// status / env validation

func InvalidStatus() DomainError {
	return simpleError{
		code: "InvalidStatus",
		msg:  "Invalid status. Allowed: active, inactive",
	}
}

func InvalidUsageEnv() DomainError {
	return simpleError{
		code: "InvalidUsageEnv",
		msg:  "Invalid usageEnv. Allowed: prod, dev, test",
	}
}

// provider / model

func ModelNotSupported(provider, model string) DomainError {
	return dynamicError{
		code: "ModelNotSupported",
		msg:  fmt.Sprintf("Model %q is not supported by provider %q", model, provider),
	}
}

func ProviderAuthFailed(optionalMsg string) DomainError {
	if optionalMsg == "" {
		return simpleError{
			code: "ProviderAuthFailed",
			msg:  "Invalid api key",
		}
	}
	return dynamicError{
		code: "ProviderAuthFailed",
		msg:  optionalMsg,
	}
}

// VO validation

func APIKeyIDRequired() DomainError {
	return simpleError{
		code: "APIKeyIDRequired",
		msg:  "API key \"_id\" is required",
	}
}

func OwnerIDRequired() DomainError {
	return simpleError{
		code: "OwnerIDRequired",
		msg:  "Owner \"_id\" is required",
	}
}

func InvalidEncryptedKeyFormat() DomainError {
	return simpleError{
		code: "InvalidEncryptedKeyFormat",
		msg:  "Invalid encrypted API key format",
	}
}

func ProviderRequired() DomainError {
	return simpleError{
		code: "ProviderRequired",
		msg:  "Provider is required",
	}
}

func ProviderTooLong(max int) DomainError {
	return limitError{
		code:   "ProviderTooLong",
		format: "Provider is too long (max %d)",
		limit:  max,
		def:    100,
	}
}

func ModelNameRequired() DomainError {
	return simpleError{
		code: "ModelNameRequired",
		msg:  "Model name is required",
	}
}

func ModelNameTooLong(max int) DomainError {
	return limitError{
		code:   "ModelNameTooLong",
		format: "Model name is too long (max %d)",
		limit:  max,
		def:    120,
	}
}

// plain key validation

func InvalidPlainAPIKeyValue() DomainError {
	return simpleError{
		code: "InvalidPlainAPIKeyValue",
		msg:  "Plain API key value is required",
	}
}
