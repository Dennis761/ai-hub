package promptdomain

import "time"

type HistoryPrimitive struct {
	Prompt    string  `json:"prompt"`
	Response  *string `json:"response"`
	Version   int     `json:"version"`
	CreatedAt string  `json:"createdAt"`
}

type PromptPrimitives struct {
	ID             string             `json:"_id"`
	TaskID         string             `json:"taskID"`
	Name           string             `json:"name"`
	ModelID        string             `json:"modelID"`
	PromptText     string             `json:"promptText"`
	ResponseText   *string            `json:"responseText"`
	History        []HistoryPrimitive `json:"history"`
	ExecutionOrder int                `json:"executionOrder"`
	Version        int                `json:"version"`
	CreatedAt      string             `json:"createdAt"`
	UpdatedAt      string             `json:"updatedAt"`
}

func (a Prompt) ToPrimitives() PromptPrimitives {
	hp := make([]HistoryPrimitive, 0, len(a.history))
	for _, e := range a.history {
		hp = append(hp, HistoryPrimitive{
			Prompt:    e.Prompt(),
			Response:  e.Response(),
			Version:   e.Version(),
			CreatedAt: e.CreatedAt().UTC().Format(time.RFC3339),
		})
	}
	return PromptPrimitives{
		ID:             a.id.Value(),
		TaskID:         a.taskID.Value(),
		Name:           a.name.Value(),
		ModelID:        a.modelID.Value(),
		PromptText:     a.promptText.Value(),
		ResponseText:   a.responseText.Value(),
		History:        hp,
		ExecutionOrder: a.executionOrder.Value(),
		Version:        a.version.Value(),
		CreatedAt:      a.createdAt.UTC().Format(time.RFC3339),
		UpdatedAt:      a.updatedAt.UTC().Format(time.RFC3339),
	}
}
