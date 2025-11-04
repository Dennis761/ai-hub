package taskdomain

func (t Task) ID() TaskID               { return t.id }
func (t Task) APIMethod() APIMethod     { return t.apiMethod }
func (t Task) ProjectID() TaskProjectID { return t.projectID }
