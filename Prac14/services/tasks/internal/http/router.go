package httpapi

import "net/http"

func NewRouter(handler *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handler.Health)
	mux.HandleFunc("/v1/jobs/process-task", handler.ProcessTask)
	return mux
}
