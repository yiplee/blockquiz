package messenger

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pkg/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/thirdparty/bot-api-go-client"
	"golang.org/x/sync/errgroup"
)

const (
	limit      = 2000
	batchLimit = 70
)

type Messenger struct {
	messages core.MessageStore
	cfg      Config
}

func New(messages core.MessageStore, cfg Config) *Messenger {
	if _, err := govalidator.ValidateStruct(messages); err != nil {
		panic(err)
	}

	return &Messenger{
		messages: messages,
		cfg:      cfg,
	}
}

func (m *Messenger) Run(ctx context.Context, dur time.Duration) error {
	log := logger.FromContext(ctx).WithField("operator", "messenger")
	ctx = logger.WithContext(ctx, log)

	timer := time.NewTimer(dur)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			_ = m.run(ctx)
			timer.Reset(dur)
		}
	}
}

func (m *Messenger) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	list, err := m.messages.ListPending(ctx, limit)
	if err != nil {
		log.WithError(err).Error("list pending messages")
		return err
	}

	if len(list) == 0 {
		return nil
	}

	users := map[string]bool{}

	var idx int
	for _, msg := range list {
		if users[msg.UserID] {
			continue
		}

		list[idx] = msg
		users[msg.UserID] = true
		idx++
	}

	list = list[:idx]

	var g errgroup.Group
	for idx := 0; idx < len(list); idx += batchLimit {
		r := idx + batchLimit
		if r >= len(list) {
			r = len(list)
		}

		messages := list[idx:r]

		g.Go(func() error {
			if err := m.postMessages(ctx, messages); err != nil {
				log.WithError(err).Error("post messages")
				return nil
			}

			if err := m.messages.Deletes(ctx, messages); err != nil {
				log.WithError(err).Error("delete messages")
				return err
			}

			log.Infof("post %d messages", len(messages))
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (m *Messenger) postMessages(ctx context.Context, messages []*core.Message) error {
	requests := make([]*bot.MessageRequest, len(messages))
	for idx, msg := range messages {
		var req bot.MessageRequest
		if jsoniter.UnmarshalFromString(msg.Body, &req) == nil {
			requests[idx] = &req
		}
	}

	return bot.PostMessages(ctx, requests, m.cfg.ClientID, m.cfg.SessionID, m.cfg.SessionKey)
}
