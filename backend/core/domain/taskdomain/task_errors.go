package taskdomain

import "fmt"

// DomainError is a common contract for Task domain errors.
type DomainError interface {
	error
	Code() string
}

// basic types

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

// not found

func TaskNotFound() DomainError {
	return simpleError{
		code: "TaskNotFound",
		msg:  "Task not found",
	}
}

// validation

func InvalidTaskName() DomainError {
	return simpleError{
		code: "InvalidTaskName",
		msg:  "Invalid task name",
	}
}

func TaskNameTooShort(min int) DomainError {
	return limitError{
		code:   "TaskNameTooShort",
		format: "Task name too short: min %d characters",
		limit:  min,
		def:    6,
	}
}

func TaskNameTooLong(max int) DomainError {
	return limitError{
		code:   "TaskNameTooLong",
		format: "Task name too long: max %d characters",
		limit:  max,
		def:    200,
	}
}

func TaskDescriptionTooLong(max int) DomainError {
	return limitError{
		code:   "TaskDescriptionTooLong",
		format: "Task description too long: max %d characters",
		limit:  max,
		def:    5000,
	}
}

func TaskDescriptionTooShort(min int) DomainError {
	return limitError{
		code:   "TaskDescriptionTooShort",
		format: "Task description too short: min %d characters",
		limit:  min,
		def:    1,
	}
}

// api method

func InvalidAPIEndpoint() DomainError {
	return simpleError{
		code: "InvalidAPIEndpoint",
		msg:  "Invalid API endpoint",
	}
}

func APIMethodRequired() DomainError {
	return simpleError{
		code: "APIMethodRequired",
		msg:  "API method is required",
	}
}

func APIMethodTooLong(max int) DomainError {
	return limitError{
		code:   "APIMethodTooLong",
		format: "API method too long: max %d characters",
		limit:  max,
		def:    300,
	}
}

// status

func InvalidStatus() DomainError {
	return simpleError{
		code: "InvalidStatus",
		msg:  "Invalid status. Allowed: archived, active, inactive",
	}
}

func TaskStatusRequired() DomainError {
	return simpleError{
		code: "TaskStatusRequired",
		msg:  "Task status is required",
	}
}

// required ids

func TaskIDRequired() DomainError {
	return simpleError{
		code: "TaskIDRequired",
		msg:  "Task _id is required",
	}
}

func CreatorIDRequired() DomainError {
	return simpleError{
		code: "CreatorIDRequired",
		msg:  "Creator _id is required",
	}
}
