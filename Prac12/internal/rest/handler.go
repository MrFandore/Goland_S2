package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"Prac12/internal/task"
)

type Handler struct {
	service *task.Service
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(service *task.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Tasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, h.service.List())
	case http.MethodPost:
		h.createTask(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

func (h *Handler) TaskByID(w http.ResponseWriter, r *http.Request) {
	id, ok := taskIDFromPath(r.URL.Path)
	if !ok {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "not found"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		item, err := h.service.Get(id)
		if err != nil {
			writeTaskError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	case http.MethodPatch:
		h.updateTask(w, r, id)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var input task.CreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}

	item, err := h.service.Create(input)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request, id string) {
	var input task.UpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}

	item, err := h.service.Update(id, input)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "task not found"})
			return
		}
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func taskIDFromPath(path string) (string, bool) {
	const prefix = "/v1/tasks/"
	if !strings.HasPrefix(path, prefix) {
		return "", false
	}

	id := strings.TrimPrefix(path, prefix)
	if id == "" || strings.Contains(id, "/") {
		return "", false
	}

	return id, true
}

func writeTaskError(w http.ResponseWriter, err error) {
	if errors.Is(err, task.ErrTaskNotFound) {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "task not found"})
		return
	}

	writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
