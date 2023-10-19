package interfaces

type Broker interface {
	Publish(body string) error
	Consume(customHandler func(string) error) error
	Shutdown() error
	CheckReadiness() error
}
