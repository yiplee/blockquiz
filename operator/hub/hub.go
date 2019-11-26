package hub

import (
	"context"
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/mq"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/thirdparty/bot-api-go-client"
	"golang.org/x/sync/semaphore"
)

type Hub struct {
	commands core.CommandStore
	parser   core.CommandParser
	pub      mq.Pub
	config   Config
	sem      *semaphore.Weighted
}

func New(
	commands core.CommandStore,
	parser core.CommandParser,
	pub mq.Pub,
	config Config,
) *Hub {
	if _, err := govalidator.ValidateStruct(config); err != nil {
		panic(err)
	}

	return &Hub{
		commands: commands,
		parser:   parser,
		pub:      pub,
		config:   config,
		sem:      semaphore.NewWeighted(10),
	}
}

func (h *Hub) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("operator", "hub")
	ctx = logger.WithContext(ctx, log)

	for {
		blaze := bot.NewBlazeClient(h.config.ClientID, h.config.SessionID, h.config.SessionKey)

		if err := blaze.Loop(ctx, h); err != nil {
			log.WithError(err).Error("blaze loop")

			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return err
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}
