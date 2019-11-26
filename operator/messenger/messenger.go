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

const limit = 100

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

	list, err := m.messages.ListPending(ctx, limit)
	if err != nil {
		log.WithError(err).Error("list pending messages")
		return err
	}

	if len(list) == 0 {
		return nil
	}

	requests := make([]*bot.MessageRequest, len(list))
	for idx, msg := range list {
		var req bot.MessageRequest
		if jsoniter.UnmarshalFromString(msg.Body, &req) == nil {
			requests[idx] = &req
		}
	}

	log.Debugf("post %d messages in batch", len(requests))
	if err := bot.PostMessages(ctx, requests, m.cfg.ClientID, m.cfg.SessionID, m.cfg.SessionKey); err != nil {
		log.WithError(err).Error("post messages")
		return err
	}

	if err := m.messages.Deletes(ctx, list); err != nil {
		log.WithError(err).Error("delete messages")
		return err
	}

	return nil
}
