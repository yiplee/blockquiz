package user

import (
	"context"

	"github.com/fox-one/pkg/store/db"
	"github.com/yiplee/blockquiz/core"
)

type store struct {
	db *db.DB
}

func New(db *db.DB) core.UserStore {
	return &store{db: db}
}

func (s *store) Create(ctx context.Context, user *core.User) error {
	return s.db.Update().Create(user).Error
}

func toUpdateParams(user *core.User) map[string]interface{} {
	return map[string]interface{}{
		"language": user.Language,
	}
}

func (s *store) Update(ctx context.Context, user *core.User) error {
	return s.db.Update().Model(user).Updates(toUpdateParams(user)).Error
}

func (s *store) FindMixinID(ctx context.Context, mixinID string) (*core.User, error) {
	var user core.User
	err := s.db.View().Where("mixin_id = ?", mixinID).First(&user).Error
	return &user, err
}
