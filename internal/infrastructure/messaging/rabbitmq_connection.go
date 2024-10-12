package messaging

import (
	"github.com/streadway/amqp"
)

func NewRabbitMQConnection(rabbitMQURL string) (*amqp.Connection, error) {
	return amqp.Dial(rabbitMQURL)
}
