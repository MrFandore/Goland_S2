package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"Prac13/internal/amqpclient"
	"Prac13/services/worker/internal/consumer"
)

func main() {
	rabbitURL := getEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672/")
	queueName := getEnv("QUEUE_NAME", "task_events")
	prefetch := getIntEnv("PREFETCH_COUNT", 1)

	logger := log.New(os.Stdout, "[worker] ", log.LstdFlags|log.Lmicroseconds)

	conn := amqpclient.MustConnect(rabbitURL)
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatalf("rabbit channel error: %v", err)
	}
	defer ch.Close()

	worker, err := consumer.New(ch, queueName, prefetch, logger)
	if err != nil {
		logger.Fatalf("consumer init error: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := worker.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		logger.Fatalf("worker stopped with error: %v", err)
	}

	logger.Println("worker stopped")
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func getIntEnv(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}
