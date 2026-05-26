package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Done        bool   `json:"done"`
}

type CreateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Service struct {
	mu      sync.Mutex
	tasks   map[string]Task
	counter int
}

func New() *Service {
	return &Service{
		tasks: make(map[string]Task),
	}
}

func (s *Service) Create(input CreateTaskInput) (Task, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return Task{}, errors.New("title is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	id := fmt.Sprintf("t_%03d", s.counter)
	task := Task{
		ID:          id,
		Title:       title,
		Description: strings.TrimSpace(input.Description),
		Done:        false,
	}

	s.tasks[id] = task
	return task, nil
}
