package command

import (
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/db"
)

func init() {
	db.RegisterMigrate(setCommand)
}

func setCommand(db *db.DB) error {
	tx := db.Update().Model(core.Command{})

	if err := tx.AutoMigrate(core.Command{}).Error; err != nil {
		return err
	}

	if err := tx.AddUniqueIndex("idx_commands_trace_id", "trace_id").Error; err != nil {
		return err
	}

	if err := tx.AddIndex("idx_commands_deleted_at", "deleted_at").Error; err != nil {
		return err
	}

	return nil
}
