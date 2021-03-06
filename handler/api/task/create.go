package task

import (
	"math/rand"
	"time"

	"github.com/fox-one/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/handler/api/render"
	"github.com/yiplee/blockquiz/handler/api/request"
	"github.com/yiplee/blockquiz/handler/api/view"
)

func HandleCreate(tasks core.TaskStore, courses core.CourseStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		log := logger.FromContext(ctx)

		var form struct {
			Language      string `json:"language,omitempty" valid:"in(en|zh),required"`
			UserID        string `json:"user_id,omitempty" valid:"uuid,required"`
			Creator       string `json:"creator,omitempty" valid:"stringlength(0|36)"`
			BlockDuration int64  `json:"block_duration,omitempty"` // default is 60*60+1
		}

		if err := request.BindJSON(c, &form); err != nil {
			render.BadRequest(c, err)
			return
		}

		list, err := courses.ListLanguage(ctx, form.Language)
		if err != nil {
			log.WithError(err).Error("list courses by language")
			render.InternalError(c, err)
			return
		}

		if len(list) == 0 {
			log.Errorf("missing courses with language %s", form.Language)
			render.InternalErrorf(c, "missing courses with language %s", form.Language)
			return
		}

		rand.Shuffle(len(list), func(i, j int) {
			list[i], list[j] = list[j], list[i]
		})

		course := list[0]
		blockDuration := form.BlockDuration
		if blockDuration <= 0 {
			blockDuration = 60*60 + 1
		}

		creator := form.Creator
		if creator == "" {
			creator = "api"
		}

		task := &core.Task{
			Language:      form.Language,
			UserID:        form.UserID,
			Creator:       creator,
			Info:          "",
			Course:        course.ID,
			Question:      0,
			State:         core.TaskStatePending,
			BlockDuration: blockDuration,
			BlockUntil:    time.Now(),
		}

		if err := tasks.Create(ctx, task); err != nil {
			log.WithError(err).Error("create task")
			render.InternalErrorf(c, "create task failed: %w", err)
			return
		}

		render.OK(c, view.TaskView(task, course))
	}
}
