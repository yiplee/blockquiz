package task

import (
	"github.com/fox-one/pkg/store/db"
	"github.com/yiplee/blockquiz/core"
)

func init() {
	db.RegisterMigrate(setTask)
}

func setTask(db *db.DB) error {
	tx := db.Update().Model(core.Task{})

	if err := tx.AutoMigrate(core.Task{}).Error; err != nil {
		return err
	}

	if err := tx.AddIndex("idx_tasks_user_title", "user_id", "title").Error; err != nil {
		return err
	}

	return nil
}
