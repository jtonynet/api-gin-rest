package interfaces

type Broker interface {
	IsConnected() bool
	Publish(body string) error
	Shutdown() error
	RunConsumer(consumerHandler func(string) (string, string, error)) error
}
