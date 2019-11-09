package hub

import (
	"context"
	"errors"
	"time"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/asaskevich/govalidator"
	"github.com/fox-one/pkg/logger"
	"github.com/yiplee/blockquiz/core"
)

type Hub struct {
	users    core.UserStore
	commands core.CommandStore
	config   Config
}

func New(
	users core.UserStore,
	commands core.CommandStore,
	config Config,
) *Hub {
	if _, err := govalidator.ValidateStruct(config); err != nil {
		panic(err)
	}

	return &Hub{
		users:    users,
		commands: commands,
		config:   config,
	}
}

func (h *Hub) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("operator", "hub")
	ctx = logger.WithContext(ctx, log)

	blaze := bot.NewBlazeClient(h.config.ClientID, h.config.SessionID, h.config.SessionKey)

	for {
		if err := blaze.Loop(ctx, h); err != nil {
			log.WithError(err).Error("blaze loop")

			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return err
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}
