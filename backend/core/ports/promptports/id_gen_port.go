package promptports

type IDGenerator interface {
	NewID() string
}
