package interfaces

type Handler interface {
	Execute(msg string) (string, error)
}
