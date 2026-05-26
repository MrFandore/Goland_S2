package rabbitsetup

import (
	"Prac14/internal/jobs"

	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareQueues(ch *amqp.Channel) error {
	_, err := ch.QueueDeclare(
		jobs.DLQQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	args := amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": jobs.DLQQueue,
	}

	_, err = ch.QueueDeclare(
		jobs.MainQueue,
		true,
		false,
		false,
		false,
		args,
	)

	return err
}
