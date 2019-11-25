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
	updateTask     bool
	conversationID string
	traceID        string
}

func (c *commandContext) bindTask(task *core.Task, course *core.Course, question int) {
	c.task = task
	c.course = course
	c.question, _ = course.Question(question)
	c.updateTask = c.task.IsPending()
	return
}

func (d *Deliver) createConversation(ctx context.Context, participantId string) error {
	conversationID := bot.UniqueConversationId(participantId, d.config.ClientID)
	participants := []bot.Participant{{
		UserId: participantId,
		Role:   "",
	}}

	_, err := bot.CreateConversation(ctx, "CONTACT", conversationID, participants, d.config.ClientID, d.config.SessionID, d.config.SessionKey)
	return err
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

		// new user and from api
		if cmd.Source == core.CommandSourceAPI {
			if err := d.createConversation(ctx, cmd.UserID); err != nil {
				return nil, fmt.Errorf("create conversation failed: %w", err)
			}
		}
	}

	if err != nil {
		return nil, err
	}

	c.conversationID = bot.UniqueConversationId(c.user.MixinID, d.config.ClientID)

	title := core.CourseTitleByDate(cmd.CreatedAt)
	if task, err := d.tasks.FindUser(ctx, c.user.MixinID, title); err == nil {
		if course, err := d.courses.Find(ctx, task.Title, task.Language); err == nil {
			d.shuffler.Shuffle(course, c.user.MixinID, d.config.QuestionCount)
			c.bindTask(task, course, task.Question)
		}
	}

	return c, nil
}

func (c *commandContext) Language() string {
	if c.task != nil && c.task.IsPending() {
		return c.task.Language
	}

	return c.user.Language
}

func (c *commandContext) Localizer() *localizer.Localizer {
	return localizer.WithLanguage(c.d.localizer, c.Language())
}

func (c *commandContext) preHandleCommand(ctx context.Context, cmd *core.Command) {
	task := c.task

	if task == nil {
		switch cmd.Action {
		case core.ActionAnswerQuestion:
			cmd.Action = core.ActionUsage
		}

		return
	}

	if task.IsPending() {
		if blocked, _ := task.IsBlocked(); blocked {
			// block 状态所有输入都当答错题处理
			cmd.Action = core.ActionAnswerQuestion
			cmd.Answer = -1
		} else if cmd.Action != core.ActionAnswerQuestion {
			cmd.Action = core.ActionShowQuestion
		}

		return
	}

	if task.IsDone() {
		switch cmd.Action {
		case core.ActionAnswerQuestion, core.ActionShowQuestion:
			cmd.Action = core.ActionUsage
		}

		return
	}
}

func (c *commandContext) handleCommand(ctx context.Context, cmd *core.Command) ([]*bot.MessageRequest, error) {
	log := logger.FromContext(ctx)

	var requests []*bot.MessageRequest

	// 设置语言
	switch cmd.Action {
	case core.ActionSwitchChinese, core.ActionSwitchEnglish:
		if cmd.Action != c.user.Language {
			c.user.Language = cmd.Action
			if err := c.d.users.Update(ctx, c.user); err != nil {
				return nil, fmt.Errorf("update user failed: %w", err)
			}
		}

		req := c.languageSwitched(ctx)
		requests = append(requests, req)
		return requests, nil
	}

	// 还没有设置语言
	if c.Language() == "" {
		req := c.selectLanguage(ctx, cmd)
		requests = append(requests, req)
		return requests, nil
	}

	switch cmd.Action {
	case core.ActionSwitchLanguage:
		requests = append(requests, c.selectLanguage(ctx, nil))
	case core.ActionUsage:
		finish := c.task.IsDone()
		requests = append(requests, c.showUsage(ctx, finish))
		requests = append(requests, c.showUsageButtons(ctx, finish))
	case core.ActionShowQuestion:
		if c.task == nil {
			task, course, err := c.d.createTask(ctx, c.user, cmd.CreatedAt)
			if err != nil {
				if store.IsErrNotFound(err) {
					requests = append(requests, c.showMissingCourse(ctx))
					break
				}

				return nil, err
			}

			c.bindTask(task, course, 0)
		}

		requests = append(requests, c.showQuestionContent(ctx))
		requests = append(requests, c.showQuestionChoiceButtons(ctx))
	case core.ActionAnswerQuestion:
		task := c.task

		if right := c.question.Answer == cmd.Answer; right {
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
			if blocked, _ := task.IsBlocked(); !blocked && task.BlockDuration > 0 {
				dur := time.Duration(task.BlockDuration) * time.Second
				task.BlockUntil = time.Now().Add(dur)
			}

			requests = append(requests, c.showAnswerFeedback(ctx, false))
		}
	default:
		logger.FromContext(ctx).Warnf("unknown action %s", cmd.Action)
	}

	return requests, nil
}
