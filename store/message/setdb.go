package message

import (
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/db"
)

func init() {
	db.RegisterMigrate(setMessage)
}

func setMessage(db *db.DB) error {
	tx := db.Update().Model(core.Message{})

	if err := tx.AutoMigrate(core.Message{}).Error; err != nil {
		return err
	}

	if err := tx.AddUniqueIndex("idx_messages_message_id", "message_id").Error; err != nil {
		return err
	}

	return nil
}
