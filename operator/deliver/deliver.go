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
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

const limit = 500

type Deliver struct {
	users     core.UserStore
	commands  core.CommandStore
	parser    core.CommandParser
	shuffler  core.CourseShuffler
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
	shuffler core.CourseShuffler,
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
		shuffler:  shuffler,
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

	// 将 commands 按用户分组然后并行处理
	groups := groupCommands(list)
	var g errgroup.Group
	// 最多同时处理五个用户
	sem := semaphore.NewWeighted(5)
	for _, group := range groups {
		group := group // copy ref

		if err := sem.Acquire(ctx, 1); err != nil {
			return err
		}

		g.Go(func() error {
			defer sem.Release(1)

			for _, cmd := range group {
				if err := d.handleCommand(ctx, cmd); err != nil {
					log := log.WithField("action", cmd.Action)
					log.WithError(err).Error("handle command %d", cmd.ID)
					return err
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
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

	log := logger.FromContext(ctx).WithField("user_id", c.user.MixinID)
	if c.task != nil {
		log = log.WithField("task", c.task.ID).WithField("title", c.course.Title)
	}

	log.Debugf("pre handle cmd %s", cmd.Action)
	c.preHandleCommand(ctx, cmd)
	log.Debugf("handle cmd %s", cmd.Action)

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
	if task := c.task; task != nil && task.IsPending() {
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

func (d *Deliver) createTask(ctx context.Context, user *core.User, at time.Time) (*core.Task, *core.Course, error) {
	title := core.CourseTitleByDate(at)
	course, err := d.courses.Find(ctx, title, user.Language)
	if err != nil {
		return nil, nil, err
	}

	d.shuffler.Shuffle(course, user.MixinID, d.config.QuestionCount)

	task := &core.Task{
		Version:       0,
		UserID:        user.MixinID,
		Creator:       "system",
		Title:         course.Title,
		Language:      course.Language,
		State:         core.TaskStatePending,
		BlockDuration: d.config.BlockDuration,
		BlockUntil:    time.Now(),
	}

	if err := d.tasks.Create(ctx, task); err != nil {
		return nil, nil, err
	}

	return task, course, err
}

/*
   1. 按用户 id 分组
   2. 一个用户只保留一个 usage cmd
*/
func groupCommands(list []*core.Command) map[string][]*core.Command {
	groups := make(map[string][]*core.Command)
	usages := make(map[string]bool)

	for _, cmd := range list {
		user := cmd.UserID
		isUsage := cmd.Action == core.ActionUsage

		if isUsage && usages[user] {
			continue
		}

		groups[user] = append(groups[user], cmd)

		if isUsage {
			usages[user] = true
		}
	}

	return groups
}
