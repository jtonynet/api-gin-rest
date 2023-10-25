package interfaces

type Broker interface {
	IsConnected() bool
	Publish(body string) error
	RunConsumer(customHandler func(string) (string, string, error)) error
	Shutdown() error
}
