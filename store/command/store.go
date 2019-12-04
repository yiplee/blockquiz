package command

import (
	"context"

	"github.com/fox-one/pkg/store/db"
	"github.com/yiplee/blockquiz/core"
)

type store struct {
	db *db.DB
}

func New(db *db.DB) core.CommandStore {
	return &store{
		db: db,
	}
}

func (s *store) Create(ctx context.Context, command *core.Command) error {
	tx := s.db.Update()
	err := tx.Create(command).Error
	if err != nil {
		if err := tx.Where("trace_id = ?", command.TraceID).Last(command).Error; err == nil {
			return nil
		}
	}

	return err
}

func (s *store) Delete(ctx context.Context, command *core.Command) error {
	return s.Deletes(ctx, []*core.Command{command})
}

func (s *store) Deletes(ctx context.Context, commands []*core.Command) error {
	tx := s.db.Update()
	ids := make([]int64, 0, len(commands))
	for _, cmd := range commands {
		ids = append(ids, cmd.ID)
	}

	return tx.Where("id in (?)", ids).Delete(core.Command{}).Error
}

func (s *store) ListPending(ctx context.Context, fromID int64, limit int) ([]*core.Command, error) {
	var commands []*core.Command
	err := s.db.View().Where("id > ?", fromID).Limit(limit).Find(&commands).Error
	return commands, err
}
