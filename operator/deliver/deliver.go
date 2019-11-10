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
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

const limit = 100

type Deliver struct {
	users     core.UserStore
	commands  core.CommandStore
	parser    core.CommandParser
	lessons   core.LessonStore
	wallets   core.WalletStore
	localizer *localizer.Localizer
	config    Config
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

	// group by userID
	group := make(map[string][]*core.Command)
	for _, cmd := range list {
		group[cmd.UserID] = append(group[cmd.UserID], cmd)
	}

	var g errgroup.Group
	var sem = semaphore.NewWeighted(3)

	for userID, cmds := range group {
		if err := sem.Acquire(ctx, 1); err != nil {
			return err
		}

		userID, cmds := userID, cmds
		g.Go(func() error {
			defer sem.Release(1)
			return d.post(ctx, userID, cmds)
		})
	}

	return g.Wait()
}

func (d *Deliver) post(ctx context.Context, userID string, cmds []*core.Command) error {
	log := logger.FromContext(ctx)

	var requests []*bot.MessageRequest

	for _, cmd := range cmds {
		reqs, err := d.handleCommand(ctx, cmd)
		if err != nil {
			log.WithError(err).Error("handle command")
			return err
		}

		requests = append(requests, reqs...)
	}

	if err := bot.PostMessages(ctx, requests, d.config.ClientID, d.config.SessionID, d.config.SessionKey); err != nil {
		log.WithError(err).Error("post messages")
		return err
	}

	if err := d.commands.Deletes(ctx, cmds); err != nil {
		log.WithError(err).Error("delete commands")
		return err
	}

	return nil
}

func (d *Deliver) handleCommand(ctx context.Context, cmd *core.Command) ([]*bot.MessageRequest, error) {
	var requests []*bot.MessageRequest

	user, err := d.users.FindMixinID(ctx, cmd.UserID)
	if store.IsErrNotFound(err) {
		user = &core.User{
			MixinID:  cmd.UserID,
			Language: "",
		}
		if err := d.users.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("create user failed: %w", err)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("find user by mixin id %s failed: %w", cmd.UserID, err)
	}

	// 还没有设置语言
	if user.Language == "" {

	}
}
