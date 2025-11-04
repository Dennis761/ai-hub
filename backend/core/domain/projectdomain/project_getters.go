package projectdomain

func (p Project) APIKey() ProjectAPIKey {
	return p.apiKey
}

func (p Project) OwnerID() OwnerID {
	return p.ownerID
}

func (p Project) AdminAccess() []string {
	out := make([]string, 0, len(p.adminAccess))
	for adminID := range p.adminAccess {
		out = append(out, adminID)
	}
	return out
}

func (p Project) HasAdminAccess(adminID AccessAdminID) bool {
	id := adminID.Value()
	if p.ownerID.Value() == id {
		return true
	}
	_, ok := p.adminAccess[id]
	return ok
}
