package cashier

import (
	"context"
	"time"

	"github.com/fox-one/pkg/logger"
	"github.com/yiplee/blockquiz/core"
)

const (
	limit = 100
)

type Cashier struct {
	wallets core.WalletStore
	walletz core.WalletService
}

func New(
	wallets core.WalletStore,
	walletz core.WalletService,
) *Cashier {
	return &Cashier{
		wallets: wallets,
		walletz: walletz,
	}
}

func (c *Cashier) Work(ctx context.Context, dur time.Duration) error {
	log := logger.FromContext(ctx).WithField("operator", "cashier")
	ctx = logger.WithContext(ctx, log)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			_ = c.run(ctx)
		}
	}
}

func (c *Cashier) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	requests, err := c.wallets.ListTransfers(ctx, limit)
	if err != nil {
		log.WithError(err).Error("list pending transfers")
		return err
	}

	var ids []string

	for _, req := range requests {
		if err := c.walletz.Transfer(ctx, req); err != nil {
			log.WithError(err).Errorf("transfer %s", req.TraceID)
			continue
		}

		ids = append(ids, req.TraceID)
	}

	if len(ids) > 0 {
		log.Debugf("finish %d transfers", len(ids))

		if err := c.wallets.DeleteTransfers(ctx, ids); err != nil {
			log.WithError(err).Error("delete transfers")
			return err
		}
	}

	return nil
}
