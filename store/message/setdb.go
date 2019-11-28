package message

import (
	"github.com/fox-one/pkg/store/db"
	"github.com/yiplee/blockquiz/core"
)

func init() {
	db.RegisterMigrate(setMessage)
}

func setMessage(db *db.DB) error {
	tx := db.Update().Model(core.Message{})

	if err := tx.AutoMigrate(core.Message{}).Error; err != nil {
		return err
	}

	return nil
}
