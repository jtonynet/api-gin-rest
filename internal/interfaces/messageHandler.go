package interfaces

type MessageHandler interface {
	Execute(msg string) (string, error)
}
