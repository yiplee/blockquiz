package core

import (
	"context"
)

type (
	Message struct {
		ID     int64  `gorm:"primary_key;auto_increment:false" json:"id,omitempty"`
		UserID string `gorm:"size:36" json:"user_id,omitempty"`
		Body   string `gorm:"type:text" json:"body,omitempty"`
	}

	MessageStore interface {
		Create(ctx context.Context, message *Message) error
		Creates(ctx context.Context, messages []*Message) error
		Deletes(ctx context.Context, messages []*Message) error
		ListPending(ctx context.Context, limit int) ([]*Message, error)
	}
)
