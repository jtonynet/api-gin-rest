package interfaces

type Insert interface {
	InsertMethod(msg string) (string, error)
}
