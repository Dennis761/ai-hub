package taskports

type IDGenerator interface {
	NewID() string
}
