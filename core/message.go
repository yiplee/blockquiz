package core

import (
	"context"
	"time"
)

type (
	Message struct {
		ID        int64     `gorm:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		MessageID string    `gorm:"size:36" json:"message_id,omitempty"`
		UserID    string    `gorm:"size:36" json:"user_id,omitempty"`
		Body      string    `gorm:"type:text" json:"body,omitempty"`
	}

	MessageStore interface {
		Create(ctx context.Context, msg *Message) error
		Deletes(ctx context.Context, messages []*Message) error
		ListPending(ctx context.Context, limit int) ([]*Message, error)
	}
)
