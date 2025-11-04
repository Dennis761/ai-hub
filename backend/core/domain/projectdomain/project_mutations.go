package projectdomain

func (p *Project) Rename(newName ProjectName) {
	p.name = newName
	p.touch()
}

func (p *Project) SetStatus(newStatus string) error {
	if _, ok := allowedProjectStatuses[newStatus]; !ok {
		return InvalidStatus()
	}
	if p.status != newStatus {
		p.status = newStatus
		p.touch()
	}
	return nil
}

func (p *Project) RebindAPIKey(newRef ProjectAPIKey) {
	p.apiKey = newRef
	p.touch()
}

func (p *Project) GrantAdminAccess(adminID AccessAdminID) {
	before := len(p.adminAccess)
	p.adminAccess[adminID.Value()] = struct{}{}
	if len(p.adminAccess) != before {
		p.touch()
	}
}
