package deliver

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/mq"
	"github.com/fox-one/pkg/text/localizer"
	jsoniter "github.com/json-iterator/go"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/thirdparty/bot-api-go-client"
	"golang.org/x/sync/errgroup"
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
	sub       mq.Sub
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
	sub mq.Sub,
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
		sub:       sub,
		localizer: localizer,
		config:    config,
	}
}

func (d *Deliver) Run(ctx context.Context, capacity int) error {
	log := logger.FromContext(ctx).WithField("operator", "deliver")
	ctx = logger.WithContext(ctx, log)

	var g errgroup.Group
	for i := 0; i < capacity; i++ {
		g.Go(func() error {
			return d.start(ctx)
		})
	}
	return g.Wait()
}

func (d *Deliver) start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_ = d.poll(ctx)
		}
	}
}

func (d *Deliver) poll(ctx context.Context) error {
	log := logger.FromContext(ctx)

	data, callback, err := d.sub.Receive(ctx, &mq.ReceiveOption{VisibilityTimeout: 10 * time.Second})
	if err != nil {
		log.WithError(err).Error("receive cmd message")
		return err
	}

	var cmd core.Command
	if err := jsoniter.UnmarshalFromString(data, &cmd); err != nil {
		return callback.Finish(ctx)
	}

	if cmd.UserID == "00000000-0000-0000-0000-000000000000" {
		return callback.Finish(ctx)
	}

	if err := d.commands.Create(ctx, &cmd); err != nil {
		log.WithError(err).Error("create command")
		return callback.Delay(ctx, time.Second)
	}

	if err := d.handleCommand(ctx, &cmd); err != nil {
		log.WithError(err).Errorf("handle command %d %s", cmd.ID, cmd.Action)
		return callback.Delay(ctx, time.Second)
	}

	return callback.Finish(ctx)
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
	if skip := c.preHandleCommand(ctx, cmd); skip {
		log.Debugf("skip cmd %s", cmd.Action)
		return nil
	}

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
	if c.updateTask && len(requests) > 0 {
		if err := d.tasks.UpdateVersion(ctx, c.task, cmd.ID); err != nil {
			return err
		}
	}

	return nil
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
