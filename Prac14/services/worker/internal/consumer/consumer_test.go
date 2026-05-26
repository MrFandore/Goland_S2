package consumer

import (
	"testing"

	"Prac14/internal/jobs"
)

func TestProcessTask(t *testing.T) {
	if err := processTask(jobs.TaskJob{TaskID: "t_001"}); err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if err := processTask(jobs.TaskJob{TaskID: "t_fail"}); err == nil {
		t.Fatal("expected simulated processing error")
	}
}
