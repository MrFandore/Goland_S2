package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	DueDate     string    `json:"due_date,omitempty"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	DueDate     *string `json:"due_date"`
	Done        *bool   `json:"done"`
}

var (
	tasks   = make(map[string]Task)
	tasksMu sync.RWMutex
)

func CreateTask(req CreateTaskRequest) (Task, error) {
	if req.Title == "" {
		return Task{}, fmt.Errorf("title is required")
	}
	tasksMu.Lock()
	defer tasksMu.Unlock()
	id := uuid.New().String()[:8] // короткий идентификатор
	now := time.Now()
	task := Task{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Done:        false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	tasks[id] = task
	return task, nil
}

func ListTasks() []Task {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	result := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		result = append(result, t)
	}
	return result
}

func GetTask(id string) (Task, error) {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	task, ok := tasks[id]
	if !ok {
		return Task{}, fmt.Errorf("task not found")
	}
	return task, nil
}

func UpdateTask(id string, req UpdateTaskRequest) (Task, error) {
	tasksMu.Lock()
	defer tasksMu.Unlock()
	task, ok := tasks[id]
	if !ok {
		return Task{}, fmt.Errorf("task not found")
	}
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.DueDate != nil {
		task.DueDate = *req.DueDate
	}
	if req.Done != nil {
		task.Done = *req.Done
	}
	task.UpdatedAt = time.Now()
	tasks[id] = task
	return task, nil
}

func DeleteTask(id string) error {
	tasksMu.Lock()
	defer tasksMu.Unlock()
	if _, ok := tasks[id]; !ok {
		return fmt.Errorf("task not found")
	}
	delete(tasks, id)
	return nil
}
