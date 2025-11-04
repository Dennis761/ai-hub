package promptdomain

func (a Prompt) ID() PromptID           { return a.id }
func (a Prompt) TaskID() TaskRefID      { return a.taskID }
func (a Prompt) Version() VersionNumber { return a.version }

func (a Prompt) History() []HistoryEntry {
	if len(a.history) == 0 {
		return nil
	}
	cp := make([]HistoryEntry, len(a.history))
	copy(cp, a.history)
	return cp
}
