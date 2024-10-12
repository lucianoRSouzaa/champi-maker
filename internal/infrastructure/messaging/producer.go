package messaging

type Producer interface {
	Publish(message interface{}) error
}
