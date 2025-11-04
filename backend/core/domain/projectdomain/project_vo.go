package projectdomain

import (
	"strings"

	"ai_hub.com/app/core/domain/admindomain"
)

const (
	ProjectNameMinLen     = 6
	ProjectNameMaxLen     = 100
	ProjectPasswordMinLen = 8
	ProjectPasswordMaxLen = 100
)

// ----------------- ProjectID -----------------

type ProjectID struct {
	value string
}

func NewProjectID(raw string) (ProjectID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ProjectID{}, ProjectIDRequired()
	}
	return ProjectID{value: trimmed}, nil
}

func (id ProjectID) Value() string { return id.value }

// ----------------- ProjectName -----------------

type ProjectName struct {
	value string
}

func NewProjectName(raw string) (ProjectName, error) {
	trimmed := strings.TrimSpace(raw)

	l := len(trimmed)
	if l < ProjectNameMinLen {
		return ProjectName{}, ProjectNameTooShort(ProjectNameMinLen)
	}
	if l > ProjectNameMaxLen {
		return ProjectName{}, ProjectNameTooLong(ProjectNameMaxLen)
	}

	return ProjectName{value: trimmed}, nil
}

func (n ProjectName) Value() string { return n.value }

// ----------------- PlainProjectAPIKey (NEW) -----------------
type PlainProjectAPIKey struct {
	value string
}

func NewPlainProjectAPIKey(raw string) (PlainProjectAPIKey, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return PlainProjectAPIKey{}, ProjectPasswordRequired()
	}

	l := len(trimmed)
	if l < ProjectPasswordMinLen {
		return PlainProjectAPIKey{}, ProjectPasswordTooShort(ProjectPasswordMinLen)
	}
	if l > ProjectPasswordMaxLen {
		return PlainProjectAPIKey{}, ProjectPasswordTooLong(ProjectPasswordMaxLen)
	}
	return PlainProjectAPIKey{value: trimmed}, nil
}

func (k PlainProjectAPIKey) Value() string {
	return k.value
}

// ----------------- ProjectAPIKey -----------------

type ProjectAPIKey struct {
	value string
}

func NewHashedAPIKey(raw string) (ProjectAPIKey, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ProjectAPIKey{}, APIKeyRequired()
	}
	return ProjectAPIKey{value: trimmed}, nil
}

func (k ProjectAPIKey) Value() string { return k.value }

// ----------------- OwnerID -----------------

type OwnerID struct {
	value string
}

func NewOwnerID(raw string) (OwnerID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return OwnerID{}, OwnerIDRequired()
	}
	return OwnerID{value: trimmed}, nil
}

func (id OwnerID) Value() string { return id.value }

// ----------------- AccessAdminID -----------------

type AccessAdminID struct {
	value string
}

func NewAccessAdminID(raw string) (AccessAdminID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return AccessAdminID{}, admindomain.AdminIDRequired()
	}
	return AccessAdminID{value: trimmed}, nil
}

func (id AccessAdminID) Value() string { return id.value }
