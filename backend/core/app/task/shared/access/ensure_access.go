package access

// EnsureAccess checks if the admin has permission to access the project.
func EnsureAccess(ownerID string, adminAccess []string, adminID string) bool {
	if ownerID == adminID {
		return true
	}

	for _, id := range adminAccess {
		if id == adminID {
			return true
		}
	}

	return false
}
