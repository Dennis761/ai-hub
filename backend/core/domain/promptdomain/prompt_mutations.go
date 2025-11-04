package promptdomain

func (a *Prompt) Rename(newName PromptName) {
	a.name = newName
	a.touch()
}

func (a *Prompt) RebindModel(newModelID ModelRefID) {
	a.modelID = newModelID
	a.touch()
}

func (a *Prompt) SetPromptText(newText PromptText, bumpVersion bool) error {
	a.promptText = newText
	if bumpVersion {
		next, err := NewVersionNumber(a.version.Value() + 1)
		if err != nil {
			return err
		}
		a.version = next
	}
	a.touch()
	return nil
}

func (a *Prompt) SetResponseText(newResponse ResponseText) {
	a.responseText = newResponse
	a.touch()
}

func (a *Prompt) SetExecutionOrder(order ExecutionOrder) {
	a.executionOrder = order
	a.touch()
}

func (a *Prompt) AddHistoryEntry(entry HistoryEntry) {
	a.history = append(a.history, entry)
	a.touch()
}

func (a *Prompt) ForceSetVersion(v VersionNumber) {
	a.version = v
	a.touch()
}

func (a *Prompt) RollbackToIndex(index int) error {
	if index < 0 || index >= len(a.history) {
		return HistoryIndexOutOfRange()
	}

	snap := a.history[index]

	pt, err := NewPromptText(snap.Prompt())
	if err != nil {
		return err
	}
	rt, err := NewResponseText(snap.Response())
	if err != nil {
		return err
	}

	maxVer := 0
	for _, h := range a.history {
		if h.Version() > maxVer {
			maxVer = h.Version()
		}
	}
	next, err := NewVersionNumber(maxVer + 1)
	if err != nil {
		return err
	}

	a.promptText = pt
	a.responseText = rt
	a.version = next
	a.touch()
	return nil
}
