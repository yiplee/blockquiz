package task

import (
	"context"

	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/db"
	"github.com/yiplee/blockquiz/store"
)

type taskStore struct {
	db *db.DB
}

func New(db *db.DB) core.TaskStore {
	return &taskStore{db: db}
}

func (t *taskStore) Create(ctx context.Context, task *core.Task) error {
	return t.db.Update().Create(task).Error
}

func (t *taskStore) Update(ctx context.Context, task *core.Task) error {
	return t.UpdateVersion(ctx, task, task.Version+1)
}

func toUpdateParams(task *core.Task) map[string]interface{} {
	return map[string]interface{}{
		"state":       task.State,
		"block_until": task.BlockUntil,
		"question":    task.Question,
	}
}

func (t *taskStore) UpdateVersion(ctx context.Context, task *core.Task, version int64) error {
	updates := toUpdateParams(task)
	updates["version"] = version

	tx := t.db.Update().Model(task).Where("version = ?", task.Version).Updates(updates)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return store.ErrOptimisticLock
	}

	return nil
}

func (t *taskStore) FindUser(ctx context.Context, userID string) (*core.Task, error) {
	var task core.Task
	err := t.db.View().Where("user_id = ?", userID).Last(&task).Error
	return &task, err
}
