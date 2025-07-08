package model

import (
	"errors"
	"time"
)

var (
	ErrNoDescription  = errors.New("description is required")
	ErrCantTransition = errors.New("can't transition to this status")
	ErrInvalidStatus  = errors.New("invalid status")
	ErrInvalidDates   = errors.New("created date must be before updated date")
)

type Task struct {
	ID          int
	Description string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTask(description string) Task {
	return Task{
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (t *Task) UpdateDescription(newDesc string) error {
	if newDesc == "" {
		return ErrNoDescription
	}

	t.Description = newDesc
	return nil
}

func (t *Task) MarkInProgress() error {
	if !t.Status.CanTransitionTo(StatusInProgress) {
		return ErrCantTransition
	}
	t.Status = StatusInProgress
	return nil
}

func (t *Task) MarkDone() error {
	if !t.Status.CanTransitionTo(StatusDone) {
		return ErrCantTransition
	}
	t.Status = StatusDone
	return nil
}

func (t *Task) Validate() error {
	if t.Description == "" {
		return ErrNoDescription
	}

	if !t.Status.IsValid() {
		return ErrInvalidStatus
	}

	if t.CreatedAt.After(t.UpdatedAt) {
		return ErrInvalidDates
	}

	return nil
}
