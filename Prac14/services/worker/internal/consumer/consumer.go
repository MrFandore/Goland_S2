package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"Prac14/internal/jobs"
	"Prac14/internal/publisher"
	"Prac14/services/worker/internal/store"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	channel   *amqp.Channel
	publisher *publisher.Publisher
	store     *store.ProcessedStore
	prefetch  int
	logger    *log.Logger
}

func New(
	channel *amqp.Channel,
	publisher *publisher.Publisher,
	store *store.ProcessedStore,
	prefetch int,
	logger *log.Logger,
) (*Consumer, error) {
	if err := channel.Qos(prefetch, 0, false); err != nil {
		return nil, err
	}

	return &Consumer{
		channel:   channel,
		publisher: publisher,
		store:     store,
		prefetch:  prefetch,
		logger:    logger,
	}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		jobs.MainQueue,
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

	c.logger.Printf("worker started queue=%s dlq=%s prefetch=%d", jobs.MainQueue, jobs.DLQQueue, c.prefetch)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d, ok := <-msgs:
			if !ok {
				return nil
			}

			if err := c.handleDelivery(ctx, d); err != nil {
				c.logger.Printf("delivery_handle_error err=%v", err)
			}
		}
	}
}

func (c *Consumer) handleDelivery(ctx context.Context, d amqp.Delivery) error {
	var job jobs.TaskJob
	if err := json.Unmarshal(d.Body, &job); err != nil {
		c.logger.Printf("bad_message err=%v body=%s", err, string(d.Body))
		return d.Ack(false)
	}

	c.logger.Printf("received job=%s task_id=%s attempt=%d message_id=%s", job.Job, job.TaskID, job.Attempt, job.MessageID)

	if c.store.Exists(job.MessageID) {
		c.logger.Printf("duplicate_message message_id=%s task_id=%s action=skip", job.MessageID, job.TaskID)
		return d.Ack(false)
	}

	if err := processTask(job); err != nil {
		c.logger.Printf("process_error task_id=%s attempt=%d message_id=%s err=%v", job.TaskID, job.Attempt, job.MessageID, err)
		return c.handleRetryOrDLQ(ctx, d, job)
	}

	c.store.MarkDone(job.MessageID)
	c.logger.Printf("process_success task_id=%s attempt=%d message_id=%s action=ack", job.TaskID, job.Attempt, job.MessageID)
	return d.Ack(false)
}

func (c *Consumer) handleRetryOrDLQ(ctx context.Context, d amqp.Delivery, job jobs.TaskJob) error {
	job.Attempt++

	if job.Attempt <= jobs.MaxAttempts {
		if err := c.publisher.PublishJob(ctx, jobs.MainQueue, job); err != nil {
			return fmt.Errorf("retry publish: %w", err)
		}

		c.logger.Printf("retry_published task_id=%s attempt=%d message_id=%s queue=%s", job.TaskID, job.Attempt, job.MessageID, jobs.MainQueue)
		return d.Ack(false)
	}

	if err := c.publisher.PublishJob(ctx, jobs.DLQQueue, job); err != nil {
		return fmt.Errorf("dlq publish: %w", err)
	}

	c.logger.Printf("dlq_published task_id=%s attempt=%d message_id=%s queue=%s", job.TaskID, job.Attempt, job.MessageID, jobs.DLQQueue)
	return d.Ack(false)
}

func processTask(job jobs.TaskJob) error {
	time.Sleep(2 * time.Second)

	if job.TaskID == "t_fail" {
		return fmt.Errorf("simulated processing error")
	}

	return nil
}
