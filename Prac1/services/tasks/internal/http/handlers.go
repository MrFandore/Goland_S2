package http

import (
	"encoding/json"
	"net/http"

	"Prac1/services/tasks/client/authclient"
	"Prac1/services/tasks/internal/service"
	"Prac1/shared/middleware"
	"github.com/gorilla/mux"
)

var authClient *authclient.Client

func SetAuthClient(client *authclient.Client) {
	authClient = client
}

// authMiddleware проверяет токен через Auth service перед выполнением хендлера.
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		requestID := middleware.GetRequestID(r.Context())
		valid, err := authClient.Verify(r.Context(), token, requestID)
		if err != nil {
			http.Error(w, "auth service unavailable", http.StatusServiceUnavailable)
			return
		}
		if !valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("authorization")
	if len(authHeader) > 7 && authHeader[:7] == "bearer " {
		return authHeader[7:]
	}
	return ""
}

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/v1/tasks", authMiddleware(createTaskHandler)).Methods("POST")
	r.HandleFunc("/v1/tasks", authMiddleware(listTasksHandler)).Methods("GET")
	r.HandleFunc("/v1/tasks/{id}", authMiddleware(getTaskHandler)).Methods("GET")
	r.HandleFunc("/v1/tasks/{id}", authMiddleware(updateTaskHandler)).Methods("PATCH")
	r.HandleFunc("/v1/tasks/{id}", authMiddleware(deleteTaskHandler)).Methods("DELETE")
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req service.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	task, err := service.CreateTask(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func listTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := service.ListTasks()
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	task, err := service.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var req service.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	task, err := service.UpdateTask(id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := service.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
