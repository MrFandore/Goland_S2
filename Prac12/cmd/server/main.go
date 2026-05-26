package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"tech-ip-sem2-gql-vs-rest/graph"
	"tech-ip-sem2-gql-vs-rest/internal/rest"
	"tech-ip-sem2-gql-vs-rest/internal/task"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	service := task.NewService()
	restHandler := rest.NewHandler(service)
	graphHandler := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{Service: service},
	}))

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("Task GraphQL playground", "/query"))
	mux.Handle("/query", graphHandler)
	mux.HandleFunc("/v1/tasks", restHandler.Tasks)
	mux.HandleFunc("/v1/tasks/", restHandler.TaskByID)

	log.Printf("server started on http://localhost:%s", port)
	log.Printf("rest: GET/POST http://localhost:%s/v1/tasks", port)
	log.Printf("graphql playground: http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
