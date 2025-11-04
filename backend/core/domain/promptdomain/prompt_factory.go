package promptdomain

import "time"

type CreateProps struct {
	ID             PromptID
	TaskID         TaskRefID
	Name           PromptName
	ModelID        ModelRefID
	PromptText     PromptText
	ResponseText   *ResponseText
	History        []HistoryEntry
	ExecutionOrder *ExecutionOrder
	Version        *VersionNumber
	Now            *time.Time
}

// Create builds a new prompt aggregate.
func Create(props CreateProps) (Prompt, error) {
	now := time.Now().UTC()
	if props.Now != nil {
		now = props.Now.UTC()
	}

	var resp ResponseText
	if props.ResponseText != nil {
		resp = *props.ResponseText
	} else {
		v, err := NewResponseText(nil)
		if err != nil {
			return Prompt{}, err
		}
		resp = v
	}

	hist := props.History
	if hist == nil {
		hist = []HistoryEntry{}
	}

	var order ExecutionOrder
	if props.ExecutionOrder != nil {
		order = *props.ExecutionOrder
	} else {
		v, err := NewExecutionOrder(0)
		if err != nil {
			return Prompt{}, err
		}
		order = v
	}

	var ver VersionNumber
	if props.Version != nil {
		ver = *props.Version
	} else {
		v, err := NewVersionNumber(1)
		if err != nil {
			return Prompt{}, err
		}
		ver = v
	}

	return Prompt{
		id:             props.ID,
		taskID:         props.TaskID,
		name:           props.Name,
		modelID:        props.ModelID,
		promptText:     props.PromptText,
		responseText:   resp,
		history:        hist,
		executionOrder: order,
		version:        ver,
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

type RestoreProps struct {
	ID             PromptID
	TaskID         TaskRefID
	Name           PromptName
	ModelID        ModelRefID
	PromptText     PromptText
	ResponseText   ResponseText
	History        []HistoryEntry
	ExecutionOrder ExecutionOrder
	Version        VersionNumber
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Restore rebuilds an Prompt aggregate from persisted data.
func Restore(props RestoreProps) Prompt {
	return Prompt{
		id:             props.ID,
		taskID:         props.TaskID,
		name:           props.Name,
		modelID:        props.ModelID,
		promptText:     props.PromptText,
		responseText:   props.ResponseText,
		history:        append([]HistoryEntry(nil), props.History...),
		executionOrder: props.ExecutionOrder,
		version:        props.Version,
		createdAt:      props.CreatedAt.UTC(),
		updatedAt:      props.UpdatedAt.UTC(),
	}
}
