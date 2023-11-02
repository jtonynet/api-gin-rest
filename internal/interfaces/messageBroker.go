package interfaces

type MessageBroker interface {
	IsConnected() bool
	Publish(body string) error
	Shutdown() error
	RunConsumer(consumerHandler func(string) (string, error)) error
}
