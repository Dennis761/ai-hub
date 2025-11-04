// src/core/domain/prompt/prompt.errors.go
package promptdomain

import (
	"fmt"
	"strings"
)

// DomainError is a common contract for prompt domain errors.
type DomainError interface {
	error
	Code() string
}

// simpleError represents a fixed-message error.
type simpleError struct {
	code string
	msg  string
}

func (e simpleError) Error() string { return e.msg }
func (e simpleError) Code() string  { return e.code }

// limitError represents errors with numeric limits (min/max).
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

// not found / relations

func PromptNotFound() DomainError {
	return simpleError{
		code: "PromptNotFound",
		msg:  "The requested prompt could not be found.",
	}
}

func TaskNotFound() DomainError {
	return simpleError{
		code: "TaskNotFound",
		msg:  "The task associated with this prompt could not be found.",
	}
}

// validation / workflow

func InvalidExecutionOrder() DomainError {
	return simpleError{
		code: "InvalidExecutionOrder",
		msg:  "The provided execution order is invalid.",
	}
}

func HistoryIndexOutOfRange() DomainError {
	return simpleError{
		code: "HistoryIndexOutOfRange",
		msg:  "The requested history entry does not exist.",
	}
}

func MissingParameters(missing []string) DomainError {
	text := "The following required parameters are missing: " + strings.Join(missing, ", ") + "."
	return simpleError{
		code: "MissingParameters",
		msg:  text,
	}
}

func InvalidVersionNumber() DomainError {
	return simpleError{
		code: "InvalidVersionNumber",
		msg:  "The specified version number is not valid.",
	}
}

// required fields

func PromptIDRequired() DomainError {
	return simpleError{
		code: "PromptIDRequired",
		msg:  `Prompt "_id" is required.`,
	}
}

func TaskIDRequired() DomainError {
	return simpleError{
		code: "TaskIDRequired",
		msg:  `Task "_id" is required.`,
	}
}

func ModelIDRequired() DomainError {
	return simpleError{
		code: "ModelIDRequired",
		msg:  `"modelID" is required.`,
	}
}

// prompt name / text

func SamePromptConsecutive() DomainError {
	return simpleError{
		code: "SamePromptConsecutive",
		msg:  "Consecutive run of the same prompt text is not allowed.",
	}
}

func NoPlaceholdersProvided() DomainError {
	return simpleError{
		code: "NoPlaceholdersProvided",
		msg:  "prompt must contain at least one placeholder ({{...}})",
	}
}

func PromptNameRequired() DomainError {
	return simpleError{
		code: "PromptNameRequired",
		msg:  "Prompt name is required",
	}
}

func PromptNameTooShort(min int) DomainError {
	return limitError{
		code:   "PromptNameTooShort",
		format: "Prompt name is too short (min %d)",
		limit:  min,
		def:    6,
	}
}

func PromptNameTooLong(max int) DomainError {
	return limitError{
		code:   "PromptNameTooLong",
		format: "Prompt name is too long (max %d)",
		limit:  max,
		def:    120,
	}
}

func PromptTextRequired() DomainError {
	return simpleError{
		code: "PromptTextRequired",
		msg:  "Prompt text is required.",
	}
}

func PromptTextTooLong(max int) DomainError {
	return limitError{
		code:   "PromptTextTooLong",
		format: "Prompt text is too long (max %d).",
		limit:  max,
		def:    10000,
	}
}

func ResponseTextTooLong(max int) DomainError {
	return limitError{
		code:   "ResponseTextTooLong",
		format: "Response text is too long (max %d).",
		limit:  max,
		def:    50000,
	}
}
