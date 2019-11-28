package cmd

import (
	"context"

	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/db"
	"github.com/yiplee/blockquiz/store/command"
	"github.com/yiplee/blockquiz/store/course"
	"github.com/yiplee/blockquiz/store/message"
	"github.com/yiplee/blockquiz/store/property"
	"github.com/yiplee/blockquiz/store/task"
	"github.com/yiplee/blockquiz/store/user"
	"github.com/yiplee/blockquiz/store/wallet"
)

func provideDB() *db.DB {
	database := db.MustOpen(cfg.DB)
	if err := db.Migrate(database); err != nil {
		panic(err)
	}

	return database
}

func provideUserStore(db *db.DB) core.UserStore {
	return user.New(db)
}

func provideCommandStore(db *db.DB) core.CommandStore {
	return command.New(db)
}

func provideCourseStore() core.CourseStore {
	courses := course.LoadCourses(cfg.Course.Path)
	list, err := courses.ListAll(context.Background())
	if err != nil {
		panic(err)
	}

	if len(list) == 0 {
		panic("no courses")
	}

	return courses
}

func provideWalletStore(db *db.DB) core.WalletStore {
	return wallet.New(db)
}

func provideTaskStore(db *db.DB) core.TaskStore {
	return task.New(db)
}

func provideMessageStore(db *db.DB) core.MessageStore {
	return message.New(db)
}

func providePropertyStore(db *db.DB) core.PropertyStore {
	return property.New(db)
}
