package task

import (
	"errors"
	"fmt"
	"sync"
)

var ErrTaskNotFound = errors.New("task not found")

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type CreateInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateInput struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Done        *bool   `json:"done,omitempty"`
}

type Service struct {
	mu      sync.RWMutex
	tasks   map[string]Task
	counter int
}

func NewService() *Service {
	svc := &Service{
		tasks:   make(map[string]Task),
		counter: 2,
	}

	svc.tasks["t_001"] = Task{
		ID:          "t_001",
		Title:       "Первая задача",
		Description: "Учебный пример",
		Done:        false,
	}
	svc.tasks["t_002"] = Task{
		ID:          "t_002",
		Title:       "Вторая задача",
		Description: "Проверка API",
		Done:        true,
	}

	return svc
}

func (s *Service) List() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Task, 0, len(s.tasks))
	for i := 1; i <= s.counter; i++ {
		id := fmt.Sprintf("t_%03d", i)
		task, ok := s.tasks[id]
		if ok {
			result = append(result, task)
		}
	}

	return result
}

func (s *Service) Get(id string) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.tasks[id]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	return item, nil
}

func (s *Service) Create(input CreateInput) (Task, error) {
	if input.Title == "" {
		return Task{}, errors.New("title is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	item := Task{
		ID:          fmt.Sprintf("t_%03d", s.counter),
		Title:       input.Title,
		Description: input.Description,
		Done:        false,
	}

	s.tasks[item.ID] = item
	return item, nil
}

func (s *Service) Update(id string, input UpdateInput) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.tasks[id]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	if input.Title != nil {
		if *input.Title == "" {
			return Task{}, errors.New("title cannot be empty")
		}
		item.Title = *input.Title
	}

	if input.Description != nil {
		item.Description = *input.Description
	}

	if input.Done != nil {
		item.Done = *input.Done
	}

	s.tasks[id] = item
	return item, nil
}
