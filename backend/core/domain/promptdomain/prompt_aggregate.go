package promptdomain

import "time"

// Prompt — Aggregate Root.
type Prompt struct {
	id             PromptID
	taskID         TaskRefID
	name           PromptName
	modelID        ModelRefID
	promptText     PromptText
	responseText   ResponseText
	history        []HistoryEntry
	executionOrder ExecutionOrder
	version        VersionNumber
	createdAt      time.Time
	updatedAt      time.Time
}

func (a *Prompt) touch() {
	a.updatedAt = time.Now().UTC()
}
