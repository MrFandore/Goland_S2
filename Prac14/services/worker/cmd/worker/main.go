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

	"Prac14/internal/amqpclient"
	"Prac14/internal/publisher"
	"Prac14/internal/rabbitsetup"
	"Prac14/services/worker/internal/consumer"
	"Prac14/services/worker/internal/store"
)

func main() {
	rabbitURL := getEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672/")
	prefetch := getIntEnv("PREFETCH_COUNT", 1)

	logger := log.New(os.Stdout, "[worker] ", log.LstdFlags|log.Lmicroseconds)

	conn := amqpclient.MustConnect(rabbitURL)
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatalf("rabbit channel error: %v", err)
	}
	defer ch.Close()

	if err := rabbitsetup.DeclareQueues(ch); err != nil {
		logger.Fatalf("queue declare error: %v", err)
	}

	pub := publisher.New(ch)
	processed := store.NewProcessedStore()
	worker, err := consumer.New(ch, pub, processed, prefetch, logger)
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
