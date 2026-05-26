package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"Prac13/services/tasks/internal/publisher"
	"Prac13/services/tasks/internal/service"
)

const (
	PublishModeBestEffort = "best-effort"
	PublishModeStrict     = "strict"
)

type Handler struct {
	service     *service.Service
	publisher   *publisher.Publisher
	logger      *log.Logger
	publishMode string
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(service *service.Service, publisher *publisher.Publisher, logger *log.Logger, publishMode string) *Handler {
	return &Handler{
		service:     service,
		publisher:   publisher,
		logger:      logger,
		publishMode: publishMode,
	}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	var input service.CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}

	task, err := h.service.Create(input)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	requestID := strings.TrimSpace(r.Header.Get("X-Request-ID"))
	publishErr := h.publisher.PublishTaskCreated(r.Context(), task.ID, requestID)
	if publishErr != nil {
		h.logger.Printf("request_id=%s publish_error task_id=%s err=%v", requestID, task.ID, publishErr)
		if h.publishMode == PublishModeStrict {
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "task created but event publish failed"})
			return
		}
	}

	h.logger.Printf(
		"request_id=%s task_created id=%s publish_mode=%s publish_ok=%t",
		requestID,
		task.ID,
		h.publishMode,
		publishErr == nil,
	)

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":           task.ID,
		"title":        task.Title,
		"description":  task.Description,
		"done":         task.Done,
		"published":    publishErr == nil,
		"publish_mode": h.publishMode,
		"ts":           time.Now().UTC().Format(time.RFC3339),
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil && !errors.Is(err, http.ErrHandlerTimeout) {
		return
	}
}
