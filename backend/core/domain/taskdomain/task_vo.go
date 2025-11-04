package taskdomain

import (
	"regexp"
	"strings"

	"ai_hub.com/app/core/domain/projectdomain"
)

const (
	MinTaskNameLen        = 6
	MaxTaskNameLen        = 200
	MaxTaskDescriptionLen = 5000
	MaxAPIMethodLen       = 300
)

// ----------------- TaskID -----------------

type TaskID struct {
	value string
}

func NewTaskID(raw string) (TaskID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return TaskID{}, TaskIDRequired()
	}
	return TaskID{value: trimmed}, nil
}

func (id TaskID) Value() string { return id.value }

// ----------------- TaskName -----------------

type TaskName struct {
	value string
}

func NewTaskName(raw string) (TaskName, error) {
	trimmed := strings.TrimSpace(raw)
	length := len(trimmed)

	switch {
	case length < MinTaskNameLen:
		return TaskName{}, TaskNameTooShort(MinTaskNameLen)
	case length > MaxTaskNameLen:
		return TaskName{}, TaskNameTooLong(MaxTaskNameLen)
	}

	return TaskName{value: trimmed}, nil
}

func (n TaskName) Value() string { return n.value }

// ----------------- TaskDescription -------------------

type TaskDescription struct {
	value *string
}

func NewTaskDescription(raw *string) (TaskDescription, error) {
	if raw == nil {
		return TaskDescription{value: nil}, nil
	}

	trimmed := strings.TrimSpace(*raw)

	if len(trimmed) > MaxTaskDescriptionLen {
		return TaskDescription{}, TaskDescriptionTooLong(MaxTaskDescriptionLen)
	}

	copied := trimmed
	return TaskDescription{value: &copied}, nil
}

func (d TaskDescription) Value() *string { return d.value }

// ----------------- TaskProjectID -----------------

type TaskProjectID struct {
	value string
}

func NewTaskProjectID(raw string) (TaskProjectID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return TaskProjectID{}, projectdomain.ProjectIDRequired()
	}
	return TaskProjectID{value: trimmed}, nil
}

func (id TaskProjectID) Value() string { return id.value }

// ----------------- TaskCreatorID -----------------

type TaskCreatorID struct {
	value string
}

func NewTaskCreatorID(raw string) (TaskCreatorID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return TaskCreatorID{}, CreatorIDRequired()
	}
	return TaskCreatorID{value: trimmed}, nil
}

func (id TaskCreatorID) Value() string { return id.value }

// ----------------- TaskStatus -----------------

type TaskStatus struct {
	value string
}

func NewTaskStatus(raw string) (TaskStatus, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	if normalized == "" {
		return TaskStatus{}, TaskStatusRequired()
	}
	if _, ok := allowedStatuses[normalized]; !ok {
		return TaskStatus{}, InvalidStatus()
	}
	return TaskStatus{value: normalized}, nil
}

func (s TaskStatus) Value() string { return s.value }

// ----------------- APIMethod -----------------

type APIMethod struct {
	value string
}

var apiMethodRe = regexp.MustCompile(`^(?:https?:\/\/[A-Za-z0-9\.\-]+)?(?:\/[A-Za-z0-9_\-\.]*)*(?:\?[A-Za-z0-9_\-=&%]*)?$`)

func NewAPIMethod(raw string) (APIMethod, error) {
	trimmed := strings.TrimSpace(raw)
	normalized := strings.ToLower(trimmed)

	switch {
	case normalized == "":
		return APIMethod{}, APIMethodRequired()
	case len(normalized) > MaxAPIMethodLen:
		return APIMethod{}, APIMethodTooLong(MaxAPIMethodLen)
	case !apiMethodRe.MatchString(trimmed):
		return APIMethod{}, InvalidAPIEndpoint()
	}

	return APIMethod{value: normalized}, nil
}

func (m APIMethod) Value() string { return m.value }
