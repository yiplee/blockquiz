package view

import (
	"github.com/yiplee/blockquiz/core"
)

type Task struct {
	ID            int64   `json:"id,omitempty"`
	CreatedAt     int64   `json:"created_at,omitempty"`
	UpdatedAt     int64   `json:"updated_at,omitempty"`
	Language      string  `json:"language,omitempty"`
	UserID        string  `json:"user_id,omitempty"`
	Creator       string  `json:"creator,omitempty"`
	Course        *Course `json:"course,omitempty"`
	Question      int     `json:"question,omitempty"`
	TotalQuestion int     `json:"total_question,omitempty"`
	State         string  `json:"state,omitempty"`
	IsBlocked     bool    `json:"is_blocked,omitempty"`
	BlockUntil    int64   `json:"block_until,omitempty"`
	BlockDuration int64   `json:"block_duration,omitempty"`
}

func TaskView(task *core.Task, course *core.Course) *Task {
	view := &Task{
		ID:            task.ID,
		CreatedAt:     task.CreatedAt.Unix(),
		UpdatedAt:     task.UpdatedAt.Unix(),
		Language:      task.Language,
		UserID:        task.UserID,
		Creator:       task.Creator,
		Question:      task.Question,
		State:         task.State,
		BlockUntil:    task.BlockUntil.Unix(),
		BlockDuration: task.BlockDuration,
	}

	if course != nil {
		view.Course = CourseView(course)
		view.TotalQuestion = len(course.Questions)
	}

	view.IsBlocked, _ = task.IsBlocked()
	return view
}
