package model

type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

func (s Status) CanTransitionTo(newStatus Status) bool {
	switch s {
	case StatusTodo:
		return newStatus == StatusInProgress
	case StatusInProgress:
		return newStatus == StatusDone
	case StatusDone:
		return false
	default:
		return false
	}
}
