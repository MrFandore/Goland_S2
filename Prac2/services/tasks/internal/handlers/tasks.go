package handlers

import (
	"encoding/json"
	"github.com/MrFandore/Go_S2/Prac2/services/tasks/middleware"
	"net/http"
)

type Task struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Subject string `json:"subject"`
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	subject, ok := r.Context().Value(middleware.UserSubjectKey).(string)
	if !ok || subject == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Здесь можно сохранить задачу в БД, для демо возвращаем заглушку
	task := Task{
		ID:      "1",
		Title:   req.Title,
		Subject: subject,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}
