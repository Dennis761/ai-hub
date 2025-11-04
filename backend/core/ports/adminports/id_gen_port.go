package adminports

type IDGenerator interface {
	NewID() string
}
