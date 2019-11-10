package cmd

import (
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/db"
	"github.com/yiplee/blockquiz/store/command"
	"github.com/yiplee/blockquiz/store/course"
	"github.com/yiplee/blockquiz/store/user"
	"github.com/yiplee/blockquiz/store/wallet"
)

func provideDB() *db.DB {
	return db.MustOpen(cfg.DB)
}

func provideUserStore(db *db.DB) core.UserStore {
	return user.New(db)
}

func provideCommandStore(db *db.DB) core.CommandStore {
	return command.New(db, true)
}

func provideCourseStore() core.CourseStore {
	return course.LoadCourses(cfg.Course.Path)
}

func provideWalletStore(db *db.DB) core.WalletStore {
	return wallet.New(db)
}
