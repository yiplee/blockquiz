package core

import (
	"context"
	"time"
)

type User struct {
	ID        int64     `gorm:"PRIMARY_KEY" json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	MixinID   string    `gorm:"size:36" json:"id,omitempty"`
	Language  string    `gorm:"size:36" json:"language,omitempty"`
}

type UserStore interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	FindMixinID(ctx context.Context, mixinID string) (*User, error)
}
