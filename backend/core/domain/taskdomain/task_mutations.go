package taskdomain

func (t *Task) Rename(newName TaskName) {
	t.name = newName
	t.touch()
}

func (t *Task) SetDescription(description TaskDescription) {
	t.description = description
	t.touch()
}

func (t *Task) SetAPIMethod(method APIMethod) {
	t.apiMethod = method
	t.touch()
}

func (t *Task) SetStatus(newStatus string) error {
	if _, ok := allowedStatuses[newStatus]; !ok {
		return InvalidStatus()
	}
	if t.status != newStatus {
		t.status = newStatus
		t.touch()
	}
	return nil
}
