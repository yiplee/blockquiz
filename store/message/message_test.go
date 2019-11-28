package message

import (
	"context"
	"testing"

	"github.com/fox-one/pkg/store/db"
	"github.com/stretchr/testify/assert"
	"github.com/yiplee/blockquiz/core"
)

var ctx = context.Background()

func newStore(t *testing.T) core.MessageStore {
	c, err := db.Open(db.SqliteInMemory())
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Migrate(c); err != nil {
		t.Fatal(err)
	}

	c.Update().Delete(core.Message{})

	return New(c)
}

func TestMessageStore(t *testing.T) {
	s := newStore(t)

	t.Run("insert commands", func(t *testing.T) {
		err := s.Creates(ctx, []*core.Message{
			{
				UserID: "1",
			},
			{
				UserID: "1",
			},
			{
				UserID: "2",
			},
		})
		assert.Nil(t, err)
	})

	t.Run("list 1 pending commands", func(t *testing.T) {
		commands, err := s.ListPending(ctx, 1)
		assert.Nil(t, err)
		assert.Len(t, commands, 1)
		assert.Equal(t, "1", commands[0].UserID)
	})

	t.Run("list 2 pending commands", func(t *testing.T) {
		commands, err := s.ListPending(ctx, 2)
		assert.Nil(t, err)
		assert.Len(t, commands, 2)
		assert.Equal(t, "1", commands[0].UserID)
		assert.Equal(t, "2", commands[1].UserID)
	})
}
