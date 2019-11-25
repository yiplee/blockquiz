package task

import (
	"strconv"
	"time"

	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/uuid"
	"github.com/gin-gonic/gin"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/handler/api/render"
	"github.com/yiplee/blockquiz/handler/api/request"
	"github.com/yiplee/blockquiz/store"
)

func Required(tasks core.TaskStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		if id := c.Param("id"); uuid.IsUUID(id) {
			requiredByUser(c, tasks, id)
		} else {
			taskID, _ := strconv.ParseInt(id, 10, 64)
			requiredByTaskID(c, tasks, taskID)
		}
	}
}

func requiredByUser(c *gin.Context, tasks core.TaskStore, userID string) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	title := core.CourseTitleByDate(time.Now())
	task, err := tasks.FindUser(ctx, userID, title)
	if store.IsErrNotFound(err) {
		task = &core.Task{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    userID,
			State:     core.TaskStatePending,
		}
	} else if err != nil {
		log.WithError(err).Error("api: cannot find task")
		render.InternalErrorf(c, "find task failed: %w", err)
		return
	}

	request.WithTask(c, task)
}

func requiredByTaskID(c *gin.Context, tasks core.TaskStore, id int64) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	task, err := tasks.Find(ctx, id)
	if store.IsErrNotFound(err) {
		log.WithError(err).Warn("api: cannot find task")
		render.NotFoundf(c, "task with id %d not found", id)
		return
	} else if err != nil {
		log.WithError(err).Error("api: cannot find task")
		render.InternalErrorf(c, "find task failed: %w", err)
		return
	}

	request.WithTask(c, task)
}
