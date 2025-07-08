package store

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sync"
	"task-tracker/internal/model"
	"time"
)

type TaskRepository interface {
	Save(task model.Task) (model.Task, error)
	Update(task model.Task) (model.Task, error)
	Delete(id int) error
	FindAll() ([]model.Task, error)
	FindByID(id int) (model.Task, error)
}

type JSONStore struct {
	filePath string
	mu       sync.Mutex
}

func NewJSONStore(filePath string) TaskRepository {
	return &JSONStore{
		filePath: filePath,
	}
}

func (s *JSONStore) Save(task model.Task) (model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, err := s.loadTasks()
	if err != nil {
		return model.Task{}, fmt.Errorf("loadTasks(): %w", err)
	}

	task.ID = s.nextID(tasks)
	now := time.Now()
	task.UpdatedAt = now
	task.CreatedAt = now
	tasks = append(tasks, task)

	if err := s.writeTasks(tasks); err != nil {
		return model.Task{}, fmt.Errorf("writeTasks(): %w", err)
	}
	return task, nil
}

func (s *JSONStore) Update(task model.Task) (model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, err := s.loadTasks()
	if err != nil {
		return model.Task{}, fmt.Errorf("loadTasks(): %w", err)
	}

	found := false
	for i, t := range tasks {
		if t.ID == task.ID {
			task.CreatedAt = t.CreatedAt
			task.UpdatedAt = time.Now()
			tasks[i] = task
			found = true
			break
		}
	}

	if !found {
		return model.Task{}, fmt.Errorf("task not found")
	}

	if err := s.writeTasks(tasks); err != nil {
		return model.Task{}, fmt.Errorf("writeTasks(): %w", err)
	}
	return task, nil
}

func (s *JSONStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, err := s.loadTasks()
	if err != nil {
		return fmt.Errorf("loadTasks(): %w", err)
	}

	found := false
	for i, t := range tasks {
		if t.ID == id {
			tasks = slices.Delete(tasks, i, i+1)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task not found")
	}

	if err := s.writeTasks(tasks); err != nil {
		return fmt.Errorf("writeTasks(): %w", err)
	}
	return nil
}

func (s *JSONStore) FindAll() ([]model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, err := s.loadTasks()
	if err != nil {
		return nil, fmt.Errorf("loadTasks(): %w", err)
	}
	return tasks, nil
}

func (s *JSONStore) FindByID(id int) (model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, err := s.loadTasks()
	if err != nil {
		return model.Task{}, fmt.Errorf("loadTasks(): %w", err)
	}

	for _, t := range tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return model.Task{}, fmt.Errorf("task not found")
}

func (s *JSONStore) loadTasks() ([]model.Task, error) {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		if err := os.WriteFile(s.filePath, []byte("[]"), 0644); err != nil {
			return nil, fmt.Errorf("os.WriteFile(): %w", err)
		}
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile(): %w", err)
	}

	var tasks []model.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %w", err)
	}
	return tasks, nil
}

func (s *JSONStore) writeTasks(tasks []model.Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent(): %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("os.WriteFile(): %w", err)
	}
	return nil
}

func (s *JSONStore) nextID(tasks []model.Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}
