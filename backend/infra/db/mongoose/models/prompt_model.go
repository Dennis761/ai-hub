package models

import "time"

type HistoryEntry struct {
	Prompt    string    `bson:"prompt"`
	Response  *string   `bson:"response,omitempty"`
	Version   int       `bson:"version"`
	CreatedAt time.Time `bson:"createdAt"`
}

type PromptDoc struct {
	ID             string         `bson:"_id"`
	TaskId         string         `bson:"taskId"`
	Name           string         `bson:"name"`
	ModelId        string         `bson:"modelId"`
	PromptText     string         `bson:"promptText"`
	ResponseText   *string        `bson:"responseText,omitempty"`
	History        []HistoryEntry `bson:"history"`
	ExecutionOrder int            `bson:"executionOrder"`
	Version        int            `bson:"version"`
	CreatedAt      time.Time      `bson:"createdAt"`
	UpdatedAt      time.Time      `bson:"updatedAt"`
}
