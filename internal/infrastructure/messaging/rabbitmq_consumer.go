package messaging

import (
	"champi-maker/internal/application"
	"champi-maker/internal/application/service"
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type rabbitMQConsumer struct {
	channel      *amqp.Channel
	queue        amqp.Queue
	matchService service.MatchService
}

func NewRabbitMQConsumer(conn *amqp.Connection, queueName string, matchService service.MatchService) (*rabbitMQConsumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return &rabbitMQConsumer{channel: ch, queue: q, matchService: matchService}, nil
}

func (c *rabbitMQConsumer) StartConsuming() error {
	msgs, err := c.channel.Consume(
		c.queue.Name,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var message application.ChampionshipCreatedMessage
			if err := json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				// Enviar Nack e não reencaminhar a mensagem
				d.Nack(false, false)
				continue
			}

			// Processar a mensagem e gerar as partidas
			if err := c.matchService.GenerateMatches(context.Background(), message); err != nil {
				log.Printf("Error generating matches: %v", err)
				// Enviar Nack e reencaminhar a mensagem para tentar novamente
				d.Nack(false, true)
				continue
			}

			// Se tudo deu certo, enviar Ack
			d.Ack(false)
		}
	}()

	return nil
}
