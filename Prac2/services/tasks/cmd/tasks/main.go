package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/MrFandore/Go_S2/Prac2/services/tasks/internal/client"
	"github.com/MrFandore/Go_S2/Prac2/services/tasks/internal/handlers"
	authMiddleware "github.com/MrFandore/Go_S2/Prac2/services/tasks/middleware"
)

func main() {
	authAddr := os.Getenv("AUTH_GRPC_ADDR")
	if authAddr == "" {
		authAddr = "localhost:50051"
	}

	authClient, err := client.NewAuthClient(authAddr)
	if err != nil {
		log.Fatalf("failed to create auth client: %v", err)
	}
	defer authClient.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.AuthMiddleware(authClient))
		r.Post("/tasks", handlers.CreateTask)
	})

	port := os.Getenv("TASKS_PORT")
	if port == "" {
		port = "8082"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("shutting down HTTP server...")
		if err := srv.Shutdown(nil); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	log.Printf("Tasks HTTP server listening on :%s", port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
