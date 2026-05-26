package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"Prac13/internal/amqpclient"
	httpapi "Prac13/services/tasks/internal/http"
	"Prac13/services/tasks/internal/publisher"
	"Prac13/services/tasks/internal/service"
)

func main() {
	port := getEnv("TASKS_PORT", "8082")
	rabbitURL := getEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672/")
	queueName := getEnv("QUEUE_NAME", "task_events")
	publishMode := normalizePublishMode(getEnv("PUBLISH_MODE", httpapi.PublishModeBestEffort))

	logger := log.New(os.Stdout, "[tasks] ", log.LstdFlags|log.Lmicroseconds)

	conn := amqpclient.MustConnect(rabbitURL)
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatalf("rabbit channel error: %v", err)
	}
	defer ch.Close()

	taskPublisher, err := publisher.New(ch, queueName)
	if err != nil {
		logger.Fatalf("publisher init error: %v", err)
	}

	taskService := service.New()
	handler := httpapi.NewHandler(taskService, taskPublisher, logger, publishMode)
	router := httpapi.NewRouter(handler)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 3 * time.Second,
	}

	logger.Printf("tasks service started on :%s queue=%s publish_mode=%s", port, queueName, publishMode)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("tasks service failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func normalizePublishMode(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case httpapi.PublishModeStrict:
		return httpapi.PublishModeStrict
	default:
		return httpapi.PublishModeBestEffort
	}
}
