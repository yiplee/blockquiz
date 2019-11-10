package user

import (
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/db"
)

func init() {
	db.RegisterMigrate(setUser)
}

func setUser(db *db.DB) error {
	tx := db.Update().Model(core.User{})

	if err := tx.AutoMigrate(core.User{}).Error; err != nil {
		return err
	}

	if err := tx.AddUniqueIndex("idx_users_mixin_id", "mixin_id").Error; err != nil {
		return err
	}

	return nil
}
