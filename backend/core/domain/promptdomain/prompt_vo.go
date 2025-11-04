// src/core/domain/prompt/prompt.vo.go
package promptdomain

import (
	"strings"
	"time"
)

const (
	promptNameMinLen   = 6
	promptNameMaxLen   = 120
	promptTextMaxLen   = 10000
	responseTextMaxLen = 50000
)

// ----------------- PromptID -----------------

type PromptID struct {
	value string
}

func NewPromptID(raw string) (PromptID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return PromptID{}, PromptIDRequired()
	}
	return PromptID{value: trimmed}, nil
}

func (id PromptID) Value() string { return id.value }

// ----------------- TaskRefID -----------------

type TaskRefID struct {
	value string
}

func NewTaskRefID(raw string) (TaskRefID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return TaskRefID{}, TaskIDRequired()
	}
	return TaskRefID{value: trimmed}, nil
}

func (id TaskRefID) Value() string { return id.value }

// ----------------- ModelRefID -----------------

type ModelRefID struct {
	value string
}

func NewModelRefID(raw string) (ModelRefID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ModelRefID{}, ModelIDRequired()
	}
	return ModelRefID{value: trimmed}, nil
}

func (id ModelRefID) Value() string { return id.value }

// ----------------- PromptName -----------------

type PromptName struct {
	value string
}

func NewPromptName(raw string) (PromptName, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return PromptName{}, PromptNameRequired()
	}
	if len(trimmed) < promptNameMinLen {
		return PromptName{}, PromptNameTooShort(promptNameMinLen)
	}
	if len(trimmed) > promptNameMaxLen {
		return PromptName{}, PromptNameTooLong(promptNameMaxLen)
	}
	return PromptName{value: trimmed}, nil
}

func (n PromptName) Value() string { return n.value }

// ----------------- PromptText -----------------

type PromptText struct {
	value string
}

func NewPromptText(raw string) (PromptText, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return PromptText{}, PromptTextRequired()
	}
	if len(trimmed) > promptTextMaxLen {
		return PromptText{}, PromptTextTooLong(promptTextMaxLen)
	}
	return PromptText{value: trimmed}, nil
}

func (t PromptText) Value() string { return t.value }

// ----------------- ResponseText ------------------

type ResponseText struct {
	value *string
}

func NewResponseText(raw *string) (ResponseText, error) {
	if raw == nil {
		return ResponseText{value: nil}, nil
	}
	if len(*raw) > responseTextMaxLen {
		return ResponseText{}, ResponseTextTooLong(responseTextMaxLen)
	}
	v := *raw
	return ResponseText{value: &v}, nil
}

func (r ResponseText) Value() *string { return r.value }

// ----------------- ExecutionOrder (>=0) -----------------

type ExecutionOrder struct {
	value int
}

func NewExecutionOrder(raw int) (ExecutionOrder, error) {
	if raw < 0 {
		return ExecutionOrder{}, InvalidExecutionOrder()
	}
	return ExecutionOrder{value: raw}, nil
}

func (o ExecutionOrder) Value() int { return o.value }

// ----------------- VersionNumber (>=1) -----------------

type VersionNumber struct {
	value int
}

func NewVersionNumber(raw int) (VersionNumber, error) {
	if raw < 1 {
		return VersionNumber{}, InvalidVersionNumber()
	}
	return VersionNumber{value: raw}, nil
}

func (v VersionNumber) Value() int { return v.value }

// ----------------- HistoryEntry --------------------------

type HistoryEntry struct {
	prompt    string
	response  *string
	version   int
	createdAt time.Time
}

func NewHistoryEntry(prompt string, response *string, version int, createdAt *time.Time) (HistoryEntry, error) {
	pt, err := NewPromptText(prompt)
	if err != nil {
		return HistoryEntry{}, err
	}
	rt, err := NewResponseText(response)
	if err != nil {
		return HistoryEntry{}, err
	}
	vn, err := NewVersionNumber(version)
	if err != nil {
		return HistoryEntry{}, err
	}

	ts := time.Now().UTC()
	if createdAt != nil {
		ts = createdAt.UTC()
	}

	return HistoryEntry{
		prompt:    pt.Value(),
		response:  rt.Value(),
		version:   vn.Value(),
		createdAt: ts,
	}, nil
}

func (h HistoryEntry) Prompt() string       { return h.prompt }
func (h HistoryEntry) Response() *string    { return h.response }
func (h HistoryEntry) Version() int         { return h.version }
func (h HistoryEntry) CreatedAt() time.Time { return h.createdAt }
