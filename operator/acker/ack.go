package acker

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pkg/logger"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/thirdparty/bot-api-go-client"
	"golang.org/x/sync/errgroup"
)

const (
	batchLimit    = 80
	checkpointKey = "quiz_ack_checkpoint_key"
)

func New(commands core.CommandStore, property core.PropertyStore, cfg Config) *Acker {
	if _, err := govalidator.ValidateStruct(cfg); err != nil {
		panic(err)
	}

	c, err := bot.NewCredential(cfg.ClientID, cfg.SessionID, cfg.SessionKey)
	if err != nil {
		panic(err)
	}

	return &Acker{
		commands:   commands,
		property:   property,
		credential: c,
	}
}

type Acker struct {
	commands core.CommandStore
	property core.PropertyStore

	credential *bot.Credential
	fromID     int64
}

func (a *Acker) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("operator", "acker")
	ctx = logger.WithContext(ctx, log)

	value, err := a.property.Get(ctx, checkpointKey)
	if err != nil {
		log.Panic(err)
	}

	a.fromID = value.Int64()

	dur := 12 * time.Millisecond
	timer := time.NewTimer(dur)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			if count, err := a.run(ctx); err != nil || count > 0 {
				timer.Reset(dur)
			} else {
				timer.Reset(200 * time.Millisecond)
			}
		}
	}
}

func (a *Acker) run(ctx context.Context) (int, error) {
	log := logger.FromContext(ctx)

	list, err := a.commands.ListPending(ctx, a.fromID, batchLimit*10)
	if err != nil {
		log.WithError(err).Error("list pending commands")
		return 0, err
	}

	if len(list) == 0 {
		return 0, nil
	}

	var g errgroup.Group
	for idx := 0; idx < len(list); idx += batchLimit {
		r := idx + batchLimit
		if r >= len(list) {
			r = len(list)
		}

		commands := list[idx:r]

		g.Go(func() error {
			if err := a.ack(ctx, commands); err != nil {
				log.WithError(err).Error("ack commands")
				return err
			}

			log.Infof("ack %d commands", len(commands))
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return 0, err
	}

	next := list[len(list)-1].ID
	if err := a.property.Save(ctx, checkpointKey, next); err != nil {
		log.WithError(err).Error("update checkpoint")
		// skip
	}

	a.fromID = next
	return len(list), nil
}

func (a *Acker) ack(ctx context.Context, commands []*core.Command) error {
	acks := make([]*bot.AcknowledgementRequest, len(commands))
	for idx, cmd := range commands {
		acks[idx] = &bot.AcknowledgementRequest{
			MessageID: cmd.TraceID,
			Status:    "READ",
		}
	}

	return bot.PostAcknowledgements(ctx, a.credential, acks)
}
