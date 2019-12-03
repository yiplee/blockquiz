package messenger

import (
	"context"
	"encoding/json"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pkg/logger"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/thirdparty/bot-api-go-client"
	"golang.org/x/sync/errgroup"
)

const (
	limit      = 2000
	batchLimit = 70
)

type Messenger struct {
	messages   core.MessageStore
	credential *bot.Credential
}

func New(messages core.MessageStore, cfg Config) *Messenger {
	if _, err := govalidator.ValidateStruct(messages); err != nil {
		panic(err)
	}

	c, err := bot.NewCredential(cfg.ClientID, cfg.SessionID, cfg.SessionKey)
	if err != nil {
		panic(err)
	}

	return &Messenger{
		messages:   messages,
		credential: c,
	}
}

func (m *Messenger) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("operator", "messenger")
	ctx = logger.WithContext(ctx, log)

	dur := 20 * time.Millisecond
	timer := time.NewTimer(dur)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			if count, err := m.run(ctx); err != nil || count > 0 {
				timer.Reset(dur)
			} else {
				timer.Reset(200 * time.Millisecond)
			}
		}
	}
}

func (m *Messenger) run(ctx context.Context) (int, error) {
	log := logger.FromContext(ctx)

	list, err := m.messages.ListPending(ctx, limit)
	if err != nil {
		log.WithError(err).Error("list pending messages")
		return 0, err
	}

	if len(list) == 0 {
		return 0, nil
	}

	users := map[string]bool{}

	var idx int
	for i := range list {
		msg := list[i]
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
		return 0, err
	}

	return len(list), nil
}

func (m *Messenger) postMessages(ctx context.Context, messages []*core.Message) error {
	requests := make([]json.RawMessage, len(messages))
	for idx := range messages {
		requests[idx] = json.RawMessage(messages[idx].Body)
	}

	return bot.PostRawMessages(ctx, m.credential, requests)
}
