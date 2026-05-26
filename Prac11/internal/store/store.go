package store

import (
	"errors"
	"fmt"
	"sync"
)

var ErrTaskNotFound = errors.New("task not found")

type Task struct {
	ID          string
	Title       string
	Description *string
	Done        bool
}

type Store struct {
	mu    sync.RWMutex
	tasks []*Task
	seq   int
}

func New() *Store {
	desc1 := "Учебный пример"
	desc2 := "GraphQL API"
	return &Store{
		tasks: []*Task{
			{ID: "t_001", Title: "Первая задача", Description: &desc1, Done: false},
			{ID: "t_002", Title: "Вторая задача", Description: &desc2, Done: true},
		},
		seq: 2,
	}
}

func (s *Store) All() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Task, len(s.tasks))
	copy(out, s.tasks)
	return out
}

func (s *Store) GetByID(id string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, t := range s.tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, ErrTaskNotFound
}

func (s *Store) Create(title string, desc *string) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	t := &Task{
		ID:          fmt.Sprintf("t_%03d", s.seq),
		Title:       title,
		Description: desc,
		Done:        false,
	}
	s.tasks = append(s.tasks, t)
	return t
}

func (s *Store) Update(id string, title *string, desc *string, done *bool) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range s.tasks {
		if t.ID == id {
			if title != nil {
				t.Title = *title
			}
			if desc != nil {
				t.Description = desc
			}
			if done != nil {
				t.Done = *done
			}
			return t, nil
		}
	}
	return nil, ErrTaskNotFound
}

func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return true
		}
	}
	return false
}

func StrPtr(s string) *string { return &s }
