package core

import (
	"context"
	"time"
)

const (
	TaskStatePending   = "PENDING"   // 任务初始化
	TaskStateCourse    = "COURSE"    // 正在学习
	TaskStateQuestion  = "QUESTION"  // 正在答题
	TaskStateCancelled = "CANCELLED" // 任务已取消
	TaskStateFinish    = "FINISH"    // 课程结束，顺利毕业
)

type (
	Task struct {
		ID        int64     `gorm:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		Version   int64     `json:"version,omitempty"`
		Language  string    `gorm:"size:24" json:"language,omitempty"`
		UserID    string    `gorm:"size:36" json:"user_id,omitempty"`
		Creator   string    `gorm:"size:36" json:"creator,omitempty"`
		Info      string    `gorm:"size:512" json:"info,omitempty"`
		Course    int64     `json:"course,omitempty"`
		Question  int       `json:"question,omitempty"`
		State     string    `gorm:"size:36" json:"state,omitempty"`
		// 来自 luckycoin 的强制性任务，一旦答错了，一小时后才能继续答题
		BlockUntil time.Time `json:"block_until,omitempty"`
	}

	TaskStore interface {
		Create(ctx context.Context, task *Task) error
		Update(ctx context.Context, task *Task) error
		UpdateVersion(ctx context.Context, task *Task, version int64) error
		Find(ctx context.Context, id int64) (*Task, error)
		// FindUser return user's last task
		FindUser(ctx context.Context, userID string) (*Task, error)
	}
)

func (t *Task) IsMandatory() bool {
	return t.UserID != t.Creator
}

func (t *Task) IsBlocked() (blocked bool, remain time.Duration) {
	if dur := time.Until(t.BlockUntil); dur > 0 {
		return true, dur
	}

	return
}

func (t *Task) IsDone() bool {
	return t.State == TaskStateCancelled || t.State == TaskStateFinish
}

func (t *Task) IsPending() bool {
	return t.State == TaskStatePending
}

func (t *Task) IsActive() bool {
	return t.State == TaskStateCourse || t.State == TaskStateQuestion
}
