package apikeyports

type IDGenerator interface {
	NewID() string
}
