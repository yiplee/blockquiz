package messenger

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pkg/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/thirdparty/bot-api-go-client"
)

const (
	limit      = 300
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

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			_ = m.run(ctx)
		}
	}
}

func (m *Messenger) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	start := time.Now()

	list, err := m.messages.ListPending(ctx, limit)
	if err != nil {
		log.WithError(err).Error("list pending messages")
		return err
	}

	if len(list) == 0 {
		return nil
	}

	log.Debugf("list %d pending messages in %s", len(list), time.Since(start))

	users := map[string]bool{}

	var idx int
	for _, msg := range list {
		if idx >= batchLimit {
			break
		}

		if users[msg.UserID] {
			continue
		}

		list[idx] = msg
		users[msg.UserID] = true
		idx++
	}

	list = list[:idx]

	requests := make([]*bot.MessageRequest, len(list))
	for idx, msg := range list {
		var req bot.MessageRequest
		if jsoniter.UnmarshalFromString(msg.Body, &req) == nil {
			requests[idx] = &req
		}
	}

	start = time.Now()
	if err := bot.PostMessages(ctx, requests, m.cfg.ClientID, m.cfg.SessionID, m.cfg.SessionKey); err != nil {
		log.WithError(err).Error("post messages")
		return err
	}

	log.Debugf("post %d messages in batch %s", len(requests), time.Since(start))

	start = time.Now()
	if err := m.messages.Deletes(ctx, list); err != nil {
		log.WithError(err).Error("delete messages")
		return err
	}

	log.Debugf("delete %d pending messages in %s", len(list), time.Since(start))

	return nil
}
