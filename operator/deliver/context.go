package deliver

import (
	"context"
	"fmt"
	"time"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/text/localizer"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/store"
)

type commandContext struct {
	d              *Deliver
	cmd            *core.Command
	user           *core.User
	course         *core.Course
	task           *core.Task
	question       *core.Question
	conversationID string
	traceID        string
	language       string
}

func (c *commandContext) bindTask(task *core.Task, course *core.Course) {
	c.task = task
	c.course = course
	c.question, _ = course.Question(0)
	c.language = course.Language
	return
}

func (d *Deliver) prepareContext(ctx context.Context, cmd *core.Command) (*commandContext, error) {
	c := &commandContext{
		d:       d,
		cmd:     cmd,
		traceID: cmd.TraceID,
	}

	var err error
	c.user, err = d.users.FindMixinID(ctx, cmd.UserID)
	if store.IsErrNotFound(err) {
		c.user = &core.User{
			MixinID: cmd.UserID,
		}
		err = d.users.Create(ctx, c.user)
	}

	if err != nil {
		return nil, err
	}

	c.language = c.user.Language

	c.conversationID = bot.UniqueConversationId(c.user.MixinID, d.config.ClientID)

	if task, err := d.tasks.FindUser(ctx, c.user.MixinID); err == nil && task.IsActive() {
		if course, err := d.courses.Find(ctx, task.Course); err == nil {
			c.task = task
			c.course = course
			c.question, _ = c.course.Question(c.task.Question)
			c.language = course.Language
		}
	}

	return c, nil
}

func (c *commandContext) Localizer() *localizer.Localizer {
	return localizer.WithLanguage(c.d.localizer, c.language)
}

func (c *commandContext) handleCommand(ctx context.Context, cmd *core.Command) ([]*bot.MessageRequest, error) {
	var requests []*bot.MessageRequest

	// 设置语言
	switch cmd.Action {
	case core.ActionSwitchChinese, core.ActionSwitchEnglish:
		if cmd.Action != c.user.Language {
			c.user.Language = cmd.Action
			c.language = c.user.Language
			if err := c.d.users.Update(ctx, c.user); err != nil {
				return nil, fmt.Errorf("update user failed: %w", err)
			}
		}

		req := c.languageSwitched(ctx)
		requests = append(requests, req)
		return requests, nil
	}

	// 还没有设置语言
	if c.language == "" {
		req := c.selectLanguage(ctx, cmd)
		requests = append(requests, req)
		return requests, nil
	}

	switch cmd.Action {
	case core.ActionUsage:
		requests = append(requests, c.showUsage(ctx))
	case core.ActionRandomCourse:
		course, err := c.d.pickRandomCourse(ctx, c.user)
		if err != nil {
			if store.IsErrNotFound(err) {
				break
			}
			return nil, fmt.Errorf("pick random course failed: %w", err)
		}

		task := &core.Task{
			Version:    0,
			Language:   course.Language,
			UserID:     c.user.MixinID,
			Creator:    "system",
			Course:     course.ID,
			State:      core.TaskStateCourse,
			BlockUntil: time.Now(),
		}
		if err := c.d.tasks.Create(ctx, task); err != nil {
			return nil, fmt.Errorf("create task failed: %w", err)
		}

		c.bindTask(task, course)
		cmd.Action = core.ActionShowCourse
		return c.handleCommand(ctx, cmd)
	case core.ActionShowCourse:
		requests = append(requests, c.showCourseContent(ctx))
		requests = append(requests, c.showCourseButtons(ctx))
	case core.ActionShowQuestion:
		c.task.State = core.TaskStateQuestion
		requests = append(requests, c.showQuestionContent(ctx))
		requests = append(requests, c.showQuestionChoiceButtons(ctx))
	case core.ActionAnswerQuestion:
		task := c.task

		right := c.question.Answer == cmd.Answer

		if right {
			requests = append(requests, c.showAnswerFeedback(ctx, true))
			task.Question += 1
			// 下一题
			if c.question, _ = c.course.Question(task.Question); c.question != nil {
				requests = append(requests, c.showQuestionContent(ctx))
				requests = append(requests, c.showQuestionChoiceButtons(ctx))
			} else {
				requests = append(requests, c.showFinishCourse(ctx))
				task.State = core.TaskStateFinish
			}
		} else {
			if blocked, _ := task.IsBlocked(); !blocked && task.IsMandatory() {
				task.BlockUntil = time.Now().Add(c.d.config.BlockDuration)
			}

			requests = append(requests, c.showAnswerFeedback(ctx, false))
		}

	default:
		logger.FromContext(ctx).Warnf("unknown action %s", cmd.Action)
	}

	return requests, nil
}
