package user

import (
	"context"
	"testing"

	"github.com/fox-one/pkg/store/db"
	"github.com/fox-one/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yiplee/blockquiz/core"
)

var ctx = context.Background()

func newStore(t *testing.T) core.UserStore {
	c, err := db.Open(db.SqliteInMemory())
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Migrate(c); err != nil {
		t.Fatal(err)
	}

	c.Update().Delete(core.User{})

	return New(c)
}

func TestUserStore(t *testing.T) {
	s := newStore(t)

	mixinID := uuid.New()

	t.Run("create", func(t *testing.T) {
		assert.Nil(t, s.Create(ctx, &core.User{
			MixinID:  mixinID,
			Language: "zh",
		}))
	})

	t.Run("find user", func(t *testing.T) {
		user, err := s.FindMixinID(ctx, mixinID)
		if assert.Nil(t, err) {
			assert.Equal(t, int64(1), user.ID)
		}
	})
}
