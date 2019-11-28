package hub

import (
	"context"
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pkg/logger"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/thirdparty/bot-api-go-client"
)

type Hub struct {
	commands core.CommandStore
	parser   core.CommandParser
	config   Config
}

func New(
	commands core.CommandStore,
	parser core.CommandParser,
	config Config,
) *Hub {
	if _, err := govalidator.ValidateStruct(config); err != nil {
		panic(err)
	}

	return &Hub{
		commands: commands,
		parser:   parser,
		config:   config,
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
