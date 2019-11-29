package deliver

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/bwmarrin/snowflake"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/text/localizer"
	jsoniter "github.com/json-iterator/go"
	"github.com/yiplee/blockquiz/core"
	"golang.org/x/sync/errgroup"
)

const (
	limit         = 100
	checkpointKey = "quiz_commands_checkpoint_key"
)

type Deliver struct {
	users     core.UserStore
	commands  core.CommandStore
	parser    core.CommandParser
	shuffler  core.CourseShuffler
	courses   core.CourseStore
	wallets   core.WalletStore
	tasks     core.TaskStore
	messages  core.MessageStore
	property  core.PropertyStore
	localizer *localizer.Localizer
	config    Config

	fromID int64
	node   *snowflake.Node
}

func New(
	users core.UserStore,
	commands core.CommandStore,
	parser core.CommandParser,
	shuffler core.CourseShuffler,
	courses core.CourseStore,
	wallets core.WalletStore,
	tasks core.TaskStore,
	messages core.MessageStore,
	property core.PropertyStore,
	localizer *localizer.Localizer,
	config Config,
) *Deliver {
	if _, err := govalidator.ValidateStruct(config); err != nil {
		panic(err)
	}

	node, err := snowflake.NewNode(1)
	if err != nil {
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
		messages:  messages,
		property:  property,
		localizer: localizer,
		config:    config,
		node:      node,
	}
}

func (d *Deliver) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("operator", "deliver")
	ctx = logger.WithContext(ctx, log)

	value, err := d.property.Get(ctx, checkpointKey)
	if err != nil {
		log.Panic(err)
	}

	d.fromID = value.Int64()

	dur := time.Millisecond
	timer := time.NewTimer(dur)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			if num, err := d.poll(ctx); err != nil {
				dur = 200 * time.Millisecond
			} else if num == 0 {
				dur = 300 * time.Millisecond
			} else {
				dur = 1 * time.Millisecond
			}

			timer.Reset(dur)
		}
	}
}

func (d *Deliver) poll(ctx context.Context) (int, error) {
	log := logger.FromContext(ctx)

	commands, err := d.commands.ListPending(ctx, d.fromID, limit)
	if err != nil {
		log.WithError(err).Error("list pending commands")
		return 0, err
	}

	if len(commands) == 0 {
		return 0, nil
	}

	log.Infof("list %d pending commands", len(commands))
	groups, next := groupCommands(commands)

	start := time.Now()

	var g errgroup.Group
	for _, group := range groups {
		group := group
		g.Go(func() error {
			for _, cmd := range group {
				if err := d.handleCommand(ctx, cmd); err != nil {
					return err
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return 0, err
	}

	log.Infof("handle %d commands for %d users in %s", len(commands), len(groups), time.Since(start))

	if err := d.property.Save(ctx, checkpointKey, next); err != nil {
		log.WithError(err).Errorf("save checkpoint %s", checkpointKey)
	}

	d.fromID = next
	return len(commands), nil
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

	if skip := c.preHandleCommand(ctx, cmd); skip {
		log.Debugf("skip cmd %s", cmd.Action)
		return nil
	}

	requests, err := c.handleCommand(ctx, cmd)
	if err != nil {
		return err
	}

	if len(requests) == 0 {
		return nil
	}

	messages := make([]*core.Message, len(requests))
	for idx, req := range requests {
		body, _ := jsoniter.MarshalToString(req)
		msg := &core.Message{
			ID:     d.node.Generate().Int64(),
			UserID: req.RecipientId,
			Body:   body,
		}
		messages[idx] = msg
	}

	if err := d.messages.Creates(ctx, messages); err != nil {
		log.WithError(err).Error("create messages")
		return err
	}

	// update task
	if c.updateTask {
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
func groupCommands(list []*core.Command) (map[string][]*core.Command, int64) {
	groups := make(map[string][]*core.Command)
	var next int64

	for _, cmd := range list {
		user := cmd.UserID
		groups[user] = append(groups[user], cmd)
		next = cmd.ID
	}

	return groups, next
}
