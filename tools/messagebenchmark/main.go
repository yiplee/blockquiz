package main

import (
	"time"

	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/store/db"
	"github.com/fox-one/pkg/uuid"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/store/message"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	log := logger.FromContext(ctx)

	database := db.MustOpen(db.Config{
		Dialect:  "mysql",
		Host:     "localhost",
		Port:     13306,
		User:     "root",
		Password: "yiplee",
		Database: "quiz",
	})

	if err := db.Migrate(database); err != nil {
		log.Panic(err)
	}

	messages := message.New(database)

	var g errgroup.Group

	for idx := 0; idx < 1; idx++ {
		g.Go(func() error {
			return insertMsg(ctx, messages)
		})
	}

	g.Go(func() error {
		return pollMsg(ctx, messages)
	})

	if err := g.Wait(); err != nil {
		log.WithError(err).Error(err)
	}
}

func insertMsg(ctx context.Context, messages core.MessageStore) error {
	log := logger.FromContext(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Millisecond):
			msg := core.Message{
				UserID: uuid.New(),
				Body:   `{"conversation_id":"88726e82-266e-3ee1-9b4e-bee452f2de46","recipient_id":"1a5d7c8b-4604-4dce-b9a0-91a3e4acf949","message_id":"7704aa1c-1878-5f5c-9e35-7ec4170edc32","category":"PLAIN_TEXT","data":"NC8xMCDkuIvpnaLlk6rkuKrkuI3mmK8gTWl4aW4g55qE5qC45b+D5Yqf6IO977yfCgpBIOetvuWIsOmihuWPliBCVEMKQiDnq6/liLDnq6/liqDlr4bogYrlpKkKQyDliJvlu7rlkITnp43lip/og73nmoTmnLrlmajkuroKRCDmnIDlronlhajmlrnkvr/nmoQgQlRDIOmSseWMhQo=","representative_id":"","quote_message_id":""}`,
			}

			if err := messages.Create(ctx, &msg); err != nil {
				return err
			}
		}
	}
}

func pollMsg(ctx context.Context, messages core.MessageStore) error {
	log := logger.FromContext(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Millisecond):
			start := time.Now()
			msgs, err := messages.ListPending(ctx, 100)
			if err != nil {
				return err
			}

			if len(msgs) == 0 {
				break
			}
			log.Infof("list   %d messages in %s", len(msgs), time.Since(start))

			// start = time.Now()
			if err := messages.Deletes(ctx, msgs); err != nil {
				return err
			}

			log.Infof("delete %d messages in %s", len(msgs), time.Since(start))
		}
	}
}
