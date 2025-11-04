package projectports

type IDGenerator interface {
	NewID() string
}
