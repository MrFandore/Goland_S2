package jobs

const (
	MainQueue   = "task_jobs"
	DLQQueue    = "task_jobs_dlq"
	MaxAttempts = 3
)

type TaskJob struct {
	Job       string `json:"job"`
	TaskID    string `json:"task_id"`
	Attempt   int    `json:"attempt"`
	MessageID string `json:"message_id"`
}

type ProcessTaskRequest struct {
	TaskID string `json:"task_id"`
}

type AcceptedResponse struct {
	Status    string `json:"status"`
	TaskID    string `json:"task_id"`
	Attempt   int    `json:"attempt"`
	MessageID string `json:"message_id"`
	Queue     string `json:"queue"`
}
