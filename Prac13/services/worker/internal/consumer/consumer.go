package consumer

import (
	"context"
	"encoding/json"
	"log"

	"Prac13/internal/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	channel   *amqp.Channel
	queueName string
	prefetch  int
	logger    *log.Logger
}

func New(channel *amqp.Channel, queueName string, prefetch int, logger *log.Logger) (*Consumer, error) {
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

	if err := channel.Qos(prefetch, 0, false); err != nil {
		return nil, err
	}

	return &Consumer{
		channel:   channel,
		queueName: queueName,
		prefetch:  prefetch,
		logger:    logger,
	}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	c.logger.Printf("worker started queue=%s prefetch=%d waiting_for_messages=true", c.queueName, c.prefetch)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d, ok := <-msgs:
			if !ok {
				return nil
			}

			if err := c.handleDelivery(d); err != nil {
				c.logger.Printf("delivery_handle_error err=%v", err)
			}
		}
	}
}

func (c *Consumer) handleDelivery(d amqp.Delivery) error {
	var ev events.TaskEvent
	if err := json.Unmarshal(d.Body, &ev); err != nil {
		c.logger.Printf("bad_message err=%v body=%s", err, string(d.Body))
		return d.Nack(false, false)
	}

	c.logger.Printf(
		"received event=%s task_id=%s ts=%s request_id=%s producer=%s version=%s",
		ev.Event,
		ev.TaskID,
		ev.TS,
		ev.RequestID,
		ev.Producer,
		ev.Version,
	)

	return d.Ack(false)
}
