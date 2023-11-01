package interfaces

type Consumer interface {
	RunConsumer(consumerHandler func(string) (string, error)) error
}

type Broker interface {
	IsConnected() bool
	Publish(body string) error
	Shutdown() error
	//ConsumerHandler() Consumer // <-- DA INTERFACE
	RunConsumer(consumerHandler func(string) (string, error)) error // <-- Implementacao anterior
}
