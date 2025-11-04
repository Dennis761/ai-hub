package projectdomain

import "fmt"

// DomainError is a common contract for project domain errors.
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

func ProjectNotFound() DomainError {
	return simpleError{"ProjectNotFound", "Project not found"}
}

func Forbidden() DomainError {
	return simpleError{"Forbidden", "Forbidden"}
}

// required fields

func ProjectIDRequired() DomainError {
	return simpleError{"ProjectIDRequired", "Project _id is required"}
}

func OwnerIDRequired() DomainError {
	return simpleError{"OwnerIDRequired", "Owner _id is required"}
}

func APIKeyRequired() DomainError {
	return simpleError{"APIKeyRequired", "Project API key is required"}
}

// status

func InvalidStatus() DomainError {
	return simpleError{"InvalidStatus", "Invalid status. Allowed: active, inactive, archived"}
}

// project name

func ProjectNameAlreadyUsed() DomainError {
	return simpleError{"ProjectNameAlreadyUsed", "Project name already used"}
}

func ProjectNameAlreadyExists() DomainError {
	return simpleError{"ProjectNameAlreadyExists", "Project name already exists"}
}

func ProjectNameTooShort(min int) DomainError {
	return limitError{
		code:   "ProjectNameTooShort",
		format: "Project name is too short (min %d characters)",
		limit:  min,
		def:    6,
	}
}

func ProjectNameTooLong(max int) DomainError {
	return limitError{
		code:   "ProjectNameTooLong",
		format: "Project name is too long (max %d characters)",
		limit:  max,
		def:    100,
	}
}

func InvalidProjectName(msg string) DomainError {
	if msg == "" {
		msg = "Invalid project name"
	}
	return dynamicError{"InvalidProjectName", msg}
}

// project password / secret

func ProjectPasswordRequired() DomainError {
	return simpleError{"ProjectPasswordRequired", "Project password is required"}
}

func ProjectPasswordTooShort(min int) DomainError {
	return limitError{
		code:   "ProjectPasswordTooShort",
		format: "Project password is too short (min %d characters)",
		limit:  min,
		def:    8,
	}
}

func ProjectPasswordTooWeak() DomainError {
	return simpleError{"ProjectPasswordTooWeak", "Project password is too weak (must include upper/lowercase, number, special char)"}
}

func ProjectPasswordTooLong(max int) DomainError {
	return limitError{
		code:   "ProjectPasswordTooLong",
		format: "Project password is too long (max %d characters)",
		limit:  max,
		def:    100,
	}
}

func ProjectPasswordMismatch() DomainError {
	return simpleError{"ProjectPasswordMismatch", "Project password confirmation does not match"}
}

// project auth

func InvalidProjectCredentials() DomainError {
	return simpleError{"InvalidProjectCredentials", "Invalid project name or password"}
}

// access / membership

func AlreadyProjectMember() DomainError {
	return simpleError{"AlreadyProjectMember", "Admin already has access to this project"}
}
