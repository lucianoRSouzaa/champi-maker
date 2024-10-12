package messaging

import (
	"champi-maker/internal/application"
	"champi-maker/internal/application/port"
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type rabbitMQPublisher struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRabbitMQPublisher(conn *amqp.Connection, queueName string) (port.MessagePublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &rabbitMQPublisher{channel: ch, queue: q}, nil
}

func (p *rabbitMQPublisher) PublishChampionshipCreated(ctx context.Context, championshipID uuid.UUID, teamIDs []uuid.UUID) error {
	message := application.ChampionshipCreatedMessage{
		ChampionshipID: championshipID,
		TeamIDs:        teamIDs,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
