package publisher

import (
	"context"
	"encoding/json"
	"time"

	"Prac13/internal/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channel   *amqp.Channel
	queueName string
}

func New(channel *amqp.Channel, queueName string) (*Publisher, error) {
	_, err := channel.QueueDeclare(
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

	return &Publisher{
		channel:   channel,
		queueName: queueName,
	}, nil
}

func (p *Publisher) PublishTaskCreated(ctx context.Context, taskID, requestID string) error {
	msg := events.TaskEvent{
		Event:     "task.created",
		TaskID:    taskID,
		TS:        time.Now().UTC().Format(time.RFC3339),
		RequestID: requestID,
		Producer:  "tasks",
		Version:   "v1",
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.channel.PublishWithContext(
		ctx,
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}
