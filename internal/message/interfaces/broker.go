package interfaces

type Broker interface {
	Publish(body string) error
	RunConsumer(customHandler func(string) error) error
	Shutdown() error
	CheckReadiness() error
}
