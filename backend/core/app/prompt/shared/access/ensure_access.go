package access

// check if admin has access (owner or listed in adminAccess)
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
