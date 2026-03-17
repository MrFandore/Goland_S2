package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"Prac1/services/tasks/client/authclient"
	tasksHttp "Prac1/services/tasks/internal/http"
	"Prac1/shared/middleware"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("TASKS_PORT")
	if port == "" {
		port = "8082"
	}
	authBaseURL := os.Getenv("AUTH_BASE_URL")
	if authBaseURL == "" {
		authBaseURL = "http://localhost:8081"
	}

	authClient := authclient.NewClient(authBaseURL, 3*time.Second)

	r := mux.NewRouter()
	r.Use(middleware.RequestIDMiddleware)
	r.Use(middleware.LoggingMiddleware)

	tasksHttp.SetAuthClient(authClient)
	tasksHttp.RegisterRoutes(r)

	addr := ":" + port
	log.Printf("tasks service starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
