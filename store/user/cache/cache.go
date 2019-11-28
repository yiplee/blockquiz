package usercache

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/yiplee/blockquiz/core"
)

func Cache(users core.UserStore) core.UserStore {
	return &cacheUsers{
		Cache: cache.New(12*time.Hour, time.Hour),
		users: users,
	}
}

type cacheUsers struct {
	*cache.Cache
	users core.UserStore
}

func (c *cacheUsers) set(user *core.User) {
	c.Set(user.MixinID, user, cache.DefaultExpiration)
}

func (c *cacheUsers) get(mixinID string) (*core.User, bool) {
	if v, ok := c.Get(mixinID); ok {
		return v.(*core.User), true
	}

	return nil, false
}

func (c *cacheUsers) Create(ctx context.Context, user *core.User) error {
	if err := c.users.Create(ctx, user); err != nil {
		return err
	}

	c.set(user)
	return nil
}

func (c *cacheUsers) Update(ctx context.Context, user *core.User) error {
	if err := c.users.Update(ctx, user); err != nil {
		return err
	}

	c.set(user)
	return nil
}

func (c *cacheUsers) FindMixinID(ctx context.Context, mixinID string) (*core.User, error) {
	if user, ok := c.get(mixinID); ok {
		return user, nil
	}

	user, err := c.users.FindMixinID(ctx, mixinID)
	if err != nil {
		return nil, err
	}

	c.set(user)
	return user, nil
}
