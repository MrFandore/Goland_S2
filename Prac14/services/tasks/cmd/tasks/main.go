package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"Prac14/internal/amqpclient"
	"Prac14/internal/publisher"
	"Prac14/internal/rabbitsetup"
	httpapi "Prac14/services/tasks/internal/http"
)

func main() {
	port := getEnv("TASKS_PORT", "8082")
	rabbitURL := getEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672/")

	logger := log.New(os.Stdout, "[tasks] ", log.LstdFlags|log.Lmicroseconds)

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
	handler := httpapi.NewHandler(pub)
	router := httpapi.NewRouter(handler)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 3 * time.Second,
	}

	logger.Printf("tasks service started on :%s queue=%s dlq=%s", port, "task_jobs", "task_jobs_dlq")
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
