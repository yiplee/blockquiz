package message

import (
	"context"

	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/db"
)

func New(db *db.DB) core.MessageStore {
	return &messageStore{db: db}
}

type messageStore struct {
	db *db.DB
}

func (m *messageStore) insert(tx *db.DB, msg *core.Message) error {
	update := tx.Update()
	update.Callback().Create().Remove("gorm:force_reload_after_create")

	err := update.Create(msg).Error
	if err != nil {
		if err := update.Where("message_id = ?", msg.MessageID).First(msg).Error; err == nil {
			return nil
		}
	}

	return err
}

func (m *messageStore) Create(ctx context.Context, message *core.Message) error {
	return m.insert(m.db, message)
}

func (m *messageStore) Creates(ctx context.Context, messages []*core.Message) error {
	if len(messages) == 0 {
		return nil
	}

	tx := m.db.Begin()
	defer tx.RollbackUnlessCommitted()

	for _, msg := range messages {
		if err := m.insert(tx, msg); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (m *messageStore) Deletes(ctx context.Context, messages []*core.Message) error {
	ids := make([]int64, len(messages))
	for idx, msg := range messages {
		ids[idx] = msg.ID
	}

	if len(ids) == 0 {
		return nil
	}

	return m.db.Update().Where("id IN (?)", ids).Delete(core.Message{}).Error
}

func (m *messageStore) ListPending(ctx context.Context, limit int) ([]*core.Message, error) {
	var messages []*core.Message
	if err := m.db.View().Limit(limit * 5).Find(&messages).Error; err != nil {
		return nil, err
	}

	var (
		users = map[string]bool{}
		idx   = 0
	)

	for _, msg := range messages {
		if idx >= limit {
			break
		}

		if users[msg.UserID] {
			continue
		}

		messages[idx] = msg
		users[msg.UserID] = true
		idx++
	}

	return messages[:idx], nil
}
