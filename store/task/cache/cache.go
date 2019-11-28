package taskcache

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/yiplee/blockquiz/core"
)

func Cache(tasks core.TaskStore) core.TaskStore {
	return &cacheTask{
		Cache: cache.New(24*time.Hour, time.Hour),
		tasks: tasks,
	}
}

type cacheTask struct {
	*cache.Cache
	tasks core.TaskStore
}

func (c *cacheTask) Create(ctx context.Context, task *core.Task) error {
	if err := c.tasks.Create(ctx, task); err != nil {
		return err
	}

	c.set(task)
	return nil
}

func (c *cacheTask) Update(ctx context.Context, task *core.Task) error {
	if err := c.tasks.Update(ctx, task); err != nil {
		return err
	}

	c.set(task)
	return nil
}

func (c *cacheTask) UpdateVersion(ctx context.Context, task *core.Task, version int64) error {
	if err := c.tasks.UpdateVersion(ctx, task, version); err != nil {
		return err
	}

	c.set(task)
	return nil
}

func (c *cacheTask) Find(ctx context.Context, id int64) (*core.Task, error) {
	return c.tasks.Find(ctx, id)
}

func (c *cacheTask) FindUser(ctx context.Context, userID, title string) (*core.Task, error) {
	if task, ok := c.get(userID, title); ok {
		return task, nil
	}

	task, err := c.tasks.FindUser(ctx, userID, title)
	if err != nil {
		return nil, err
	}

	c.set(task)
	return task, nil
}

func (c *cacheTask) set(task *core.Task) {
	c.Set(task.UserID+task.Title, task, cache.DefaultExpiration)
}

func (c *cacheTask) get(userID, title string) (*core.Task, bool) {
	if v, ok := c.Get(userID + title); ok {
		return v.(*core.Task), true
	}

	return nil, false
}
