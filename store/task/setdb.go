package task

import (
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/db"
)

func init() {
	db.RegisterMigrate(setTask)
}

func setTask(db *db.DB) error {
	tx := db.Update().Model(core.Task{})

	if err := tx.AutoMigrate(core.Task{}).Error; err != nil {
		return err
	}

	if err := tx.AddIndex("idx_tasks_user_id", "user_id").Error; err != nil {
		return err
	}

	return nil
}
