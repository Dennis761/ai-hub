package taskdomain

import "time"

// Task — Aggregate Root.
type Task struct {
	id          TaskID
	name        TaskName
	description TaskDescription
	projectID   TaskProjectID
	apiMethod   APIMethod
	status      string
	createdBy   TaskCreatorID
	createdAt   time.Time
	updatedAt   time.Time
}

func (t *Task) touch() {
	t.updatedAt = time.Now().UTC()
}
