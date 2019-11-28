package wallet

import (
	"github.com/fox-one/pkg/store/db"
	"github.com/yiplee/blockquiz/core"
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
