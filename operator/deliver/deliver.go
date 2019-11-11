package deliver

import (
	"context"
	"math/rand"
	"time"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/text/localizer"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/store"
)

const limit = 100

type Deliver struct {
	users     core.UserStore
	commands  core.CommandStore
	parser    core.CommandParser
	courses   core.CourseStore
	wallets   core.WalletStore
	tasks     core.TaskStore
	localizer *localizer.Localizer
	config    Config
}

func New(
	users core.UserStore,
	commands core.CommandStore,
	parser core.CommandParser,
	courses core.CourseStore,
	wallets core.WalletStore,
	tasks core.TaskStore,
	localizer *localizer.Localizer,
	config Config,
) *Deliver {
	if _, err := govalidator.ValidateStruct(config); err != nil {
		panic(err)
	}

	return &Deliver{
		users:     users,
		commands:  commands,
		parser:    parser,
		courses:   courses,
		wallets:   wallets,
		tasks:     tasks,
		localizer: localizer,
		config:    config,
	}
}

func (d *Deliver) Run(ctx context.Context, dur time.Duration) error {
	log := logger.FromContext(ctx).WithField("operator", "deliver")
	ctx = logger.WithContext(ctx, log)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			_ = d.run(ctx)
		}
	}
}

func (d *Deliver) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	list, err := d.commands.ListPending(ctx, limit)
	if err != nil {
		log.WithError(err).Error("list pending commands")
		return err
	}

	for _, cmd := range list {
		if err := d.handleCommand(ctx, cmd); err != nil {
			log := log.WithField("action", cmd.Action)
			log.WithError(err).Error("handle command %d", cmd.ID)
			return err
		}
	}

	if err := d.commands.Deletes(ctx, list); err != nil {
		log.WithError(err).Error("delete commands")
		return err
	}

	return nil
}

func (d *Deliver) handleCommand(ctx context.Context, cmd *core.Command) error {
	c, err := d.prepareContext(ctx, cmd)
	if err != nil {
		return err
	}

	isTaskOperation := !govalidator.IsIn(cmd.Action,
		core.ActionSwitchChinese,
		core.ActionSwitchEnglish,
		core.ActionUsage,
		core.ActionRandomCourse,
	)

	if task := c.task; task == nil || task.IsDone() || task.IsPending() {
		if isTaskOperation {
			cmd.Action = core.ActionUsage
		}
	} else {
		if task.Version > cmd.ID {
			return nil
		}

		if blocked, _ := task.IsBlocked(); blocked {
			// 等待状态所有输入都当答错题处理
			cmd.Action = core.ActionAnswerQuestion
			cmd.Answer = -1
		} else if !isTaskOperation {
			if task.State == core.TaskStateCourse {
				cmd.Action = core.ActionShowCourse
			} else {
				cmd.Action = core.ActionShowQuestion
			}
		}
	}

	requests, err := c.handleCommand(ctx, cmd)
	if err != nil {
		return err
	}

	for _, req := range requests {
		if err := bot.PostMessage(
			ctx,
			req.ConversationId,
			req.RecipientId,
			req.MessageId,
			req.Category,
			req.Data,
			d.config.ClientID,
			d.config.SessionID,
			d.config.SessionKey,
		); err != nil {
			return err
		}
	}

	// update task
	if task := c.task; task != nil {
		if err := d.tasks.UpdateVersion(ctx, task, cmd.ID); err != nil {
			return err
		}
	}

	return nil
}

func (d *Deliver) pickRandomCourse(ctx context.Context, user *core.User) (*core.Course, error) {
	list, err := d.courses.ListLanguage(ctx, user.Language)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, store.ErrNotFound
	}

	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})

	return list[0], nil
}
