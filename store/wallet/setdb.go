package wallet

import (
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/db"
)

func init() {
	db.RegisterMigrate(setTransfer)
}

func setTransfer(db *db.DB) error {
	tx := db.Update().Model(core.Transfer{})

	if err := tx.AutoMigrate(core.Transfer{}).Error; err != nil {
		return err
	}

	return nil
}
