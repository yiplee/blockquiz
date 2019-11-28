package message

import (
	"context"

	"github.com/fox-one/pkg/store/db"
	"github.com/yiplee/blockquiz/core"
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

	return update.Create(msg).Error
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

	return m.db.Update().Where("id IN (?)", ids).Delete(core.Message{}).Error
}

func (m *messageStore) ListPending(ctx context.Context, limit int) ([]*core.Message, error) {
	var messages []*core.Message
	err := m.db.View().Limit(limit).Find(&messages).Error
	return messages, err
}
