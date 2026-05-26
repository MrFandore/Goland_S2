package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"Prac14/internal/jobs"
	"Prac14/internal/publisher"
)

type Handler struct {
	publisher *publisher.Publisher
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(pub *publisher.Publisher) *Handler {
	return &Handler{publisher: pub}
}

func (h *Handler) ProcessTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	var req jobs.ProcessTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}

	if req.TaskID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "task_id is required"})
		return
	}

	job := jobs.TaskJob{
		Job:       "process_task",
		TaskID:    req.TaskID,
		Attempt:   1,
		MessageID: uuid.NewString(),
	}

	if err := h.publisher.PublishJob(r.Context(), jobs.MainQueue, job); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to publish job"})
		return
	}

	writeJSON(w, http.StatusAccepted, jobs.AcceptedResponse{
		Status:    "accepted",
		TaskID:    job.TaskID,
		Attempt:   job.Attempt,
		MessageID: job.MessageID,
		Queue:     jobs.MainQueue,
	})
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
