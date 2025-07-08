package service

import (
	"fmt"
	"task-tracker/internal/model"
	"task-tracker/internal/store"
	"time"
)

type StatusFilter string

const (
	FilterAll        StatusFilter = "all"
	FilterTodo       StatusFilter = "todo"
	FilterInProgress StatusFilter = "in-progress"
	FilterDone       StatusFilter = "done"
)

type TaskUsecase interface {
	Add(description string) (*model.Task, error)
	Update(id int, description string) (*model.Task, error)
	Delete(id int) error
	MarkInProgress(id int) (*model.Task, error)
	MarkDone(id int) (*model.Task, error)
	List(filter StatusFilter) ([]model.Task, error)
	GetByID(id int) (*model.Task, error)
}

type taskService struct {
	repo store.TaskRepository
}

func NewTaskUsecase(repo store.TaskRepository) TaskUsecase {
	return &taskService{repo: repo}
}

func (s *taskService) Add(description string) (*model.Task, error) {
	task := model.NewTask(description)
	if err := task.Validate(); err != nil {
		return nil, err
	}

	savedTask, err := s.repo.Save(task)
	if err != nil {
		return nil, fmt.Errorf("s.repo.Save(task): %w", err)
	}

	return &savedTask, nil
}

func (s *taskService) Update(id int, description string) (*model.Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := task.UpdateDescription(description); err != nil {
		return nil, err
	}

	task.UpdatedAt = time.Now()
	updatedTask, err := s.repo.Update(task)
	if err != nil {
		return nil, fmt.Errorf("s.repo.Update(task): %w", err)
	}

	return &updatedTask, nil
}

func (s *taskService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *taskService) MarkInProgress(id int) (*model.Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := task.MarkInProgress(); err != nil {
		return nil, err
	}

	task.UpdatedAt = time.Now()
	updatedTask, err := s.repo.Update(task)
	if err != nil {
		return nil, fmt.Errorf("s.repo.Update(task): %w", err)
	}

	return &updatedTask, nil
}

func (s *taskService) MarkDone(id int) (*model.Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := task.MarkDone(); err != nil {
		return nil, err
	}

	task.UpdatedAt = time.Now()
	updatedTask, err := s.repo.Update(task)
	if err != nil {
		return nil, fmt.Errorf("s.repo.Update(task): %w", err)
	}

	return &updatedTask, nil
}

func (s *taskService) List(filter StatusFilter) ([]model.Task, error) {
	tasks, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("s.repo.FindAll(): %w", err)
	}

	if filter == FilterAll || filter == "" {
		return tasks, nil
	}

	filtered := make([]model.Task, 0, len(tasks))
	for _, task := range tasks {
		if s.matchesFilter(task, filter) {
			filtered = append(filtered, task)
		}
	}

	return filtered, nil
}

func (s *taskService) GetByID(id int) (*model.Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *taskService) matchesFilter(task model.Task, filter StatusFilter) bool {
	switch filter {
	case FilterTodo:
		return task.Status == model.StatusTodo
	case FilterInProgress:
		return task.Status == model.StatusInProgress
	case FilterDone:
		return task.Status == model.StatusDone
	default:
		return true
	}
}
