package task

import (
	"strconv"

	"github.com/fox-one/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/handler/api/render"
	"github.com/yiplee/blockquiz/handler/api/request"
	"github.com/yiplee/blockquiz/store"
)

func TaskRequired(tasks core.TaskStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		log := logger.FromContext(ctx)

		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
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
}
