package core

import (
	"context"
	"time"
)

const (
	TaskStatePending = "PENDING" // 任务初始化
	TaskStateFinish  = "FINISH"  // 课程结束，顺利毕业
)

type (
	Task struct {
		ID            int64     `gorm:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt     time.Time `json:"created_at,omitempty"`
		UpdatedAt     time.Time `json:"updated_at,omitempty"`
		Version       int64     `json:"version,omitempty"`
		UserID        string    `gorm:"size:36" json:"user_id,omitempty"`
		Creator       string    `gorm:"size:36" json:"creator,omitempty"`
		Info          string    `gorm:"size:512" json:"info,omitempty"`
		Title         string    `gorm:"size:128" json:"title,omitempty"`
		Language      string    `gorm:"size:24" json:"language,omitempty"`
		Question      int       `json:"question,omitempty"`
		State         string    `gorm:"size:36" json:"state,omitempty"`
		BlockDuration int64     `json:"block_duration,omitempty"`
		BlockUntil    time.Time `json:"block_until,omitempty"`
	}

	TaskStore interface {
		Create(ctx context.Context, task *Task) error
		Update(ctx context.Context, task *Task) error
		UpdateVersion(ctx context.Context, task *Task, version int64) error
		Find(ctx context.Context, id int64) (*Task, error)
		FindUser(ctx context.Context, userID, title string) (*Task, error)
	}
)

func (t *Task) IsBlocked() (blocked bool, remain time.Duration) {
	if dur := time.Until(t.BlockUntil); dur > 0 {
		return true, dur
	}

	return
}

func (t *Task) IsDone() bool {
	return t != nil && t.State == TaskStateFinish
}

func (t *Task) IsPending() bool {
	return t == nil || t.State == TaskStatePending
}
