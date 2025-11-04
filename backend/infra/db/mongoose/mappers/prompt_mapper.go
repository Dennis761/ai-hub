package mappers

import (
	"time"

	"ai_hub.com/app/core/domain/promptdomain"
)

type PromptDoc struct {
	ID             string          `bson:"_id"`
	TaskID         string          `bson:"taskId"`
	Name           string          `bson:"name"`
	ModelID        string          `bson:"modelId"`
	PromptText     string          `bson:"promptText"`
	ResponseText   *string         `bson:"responseText,omitempty"`
	History        []HistoryRecord `bson:"history"`
	ExecutionOrder int             `bson:"executionOrder"`
	Version        int             `bson:"version"`
	CreatedAt      time.Time       `bson:"createdAt"`
	UpdatedAt      time.Time       `bson:"updatedAt"`
}

type HistoryRecord struct {
	Prompt    string    `bson:"prompt"`
	Response  *string   `bson:"response,omitempty"`
	Version   int       `bson:"version"`
	CreatedAt time.Time `bson:"createdAt"`
}

func PromptFromDoc(doc *PromptDoc) (*promptdomain.Prompt, error) {
	if doc == nil {
		return nil, nil
	}
	// ID
	id, err := promptdomain.NewPromptID(doc.ID)
	if err != nil {
		return nil, err
	}
	// TaskID
	task, err := promptdomain.NewTaskRefID(doc.TaskID)
	if err != nil {
		return nil, err
	}
	// ModelID
	model, err := promptdomain.NewModelRefID(doc.ModelID)
	if err != nil {
		return nil, err
	}
	// Name
	name, err := promptdomain.NewPromptName(doc.Name)
	if err != nil {
		return nil, err
	}
	// PromptText
	prompt, err := promptdomain.NewPromptText(doc.PromptText)
	if err != nil {
		return nil, err
	}
	// ResponseText
	response, err := promptdomain.NewResponseText(doc.ResponseText)
	if err != nil {
		return nil, err
	}

	// History
	history := make([]promptdomain.HistoryEntry, 0, len(doc.History))
	for _, h := range doc.History {
		created := h.CreatedAt
		entry, err := promptdomain.NewHistoryEntry(h.Prompt, h.Response, h.Version, &created)
		if err != nil {
			return nil, err
		}
		history = append(history, entry)
	}

	// ExecutionOrder
	execOrder, err := promptdomain.NewExecutionOrder(doc.ExecutionOrder)
	if err != nil {
		return nil, err
	}
	// Version
	version, err := promptdomain.NewVersionNumber(doc.Version)
	if err != nil {
		return nil, err
	}

	// CreatedAt & UpdatedAt
	createdAt := doc.CreatedAt
	updatedAt := doc.UpdatedAt

	props := promptdomain.RestoreProps{
		ID:             id,
		TaskID:         task,
		Name:           name,
		ModelID:        model,
		PromptText:     prompt,
		ResponseText:   response,
		History:        history,
		ExecutionOrder: execOrder,
		Version:        version,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}

	p := promptdomain.Restore(props)
	return &p, nil
}

func PromptToPersistence(entity *promptdomain.Prompt) (*PromptDoc, error) {
	p := entity.ToPrimitives()

	createdAt, err := time.Parse(time.RFC3339, p.CreatedAt)
	if err != nil {
		return nil, err
	}
	updatedAt, err := time.Parse(time.RFC3339, p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	historyDocs := make([]HistoryRecord, 0, len(p.History))
	for _, h := range p.History {
		ht, err := time.Parse(time.RFC3339, h.CreatedAt)
		if err != nil {
			return nil, err
		}
		historyDocs = append(historyDocs, HistoryRecord{
			Prompt:    h.Prompt,
			Response:  h.Response,
			Version:   h.Version,
			CreatedAt: ht.UTC(),
		})
	}

	return &PromptDoc{
		ID:             p.ID,
		TaskID:         p.TaskID,
		Name:           p.Name,
		ModelID:        p.ModelID,
		PromptText:     p.PromptText,
		ResponseText:   p.ResponseText,
		History:        historyDocs,
		ExecutionOrder: p.ExecutionOrder,
		Version:        p.Version,
		CreatedAt:      createdAt.UTC(),
		UpdatedAt:      updatedAt.UTC(),
	}, nil
}
