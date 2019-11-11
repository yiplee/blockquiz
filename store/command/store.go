package command

import (
	"context"

	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/db"
)

type store struct {
	db         *db.DB
	softDelete bool
}

func New(db *db.DB, softDelete bool) core.CommandStore {
	return &store{
		db:         db,
		softDelete: softDelete,
	}
}

func (s *store) Create(ctx context.Context, command *core.Command) error {
	return s.db.Update().Unscoped().Where("trace_id = ?", command.TraceID).FirstOrCreate(command).Error
}

func (s *store) Delete(ctx context.Context, command *core.Command) error {
	return s.Deletes(ctx, []*core.Command{command})
}

func (s *store) Deletes(ctx context.Context, commands []*core.Command) error {
	tx := s.db.Update()
	if !s.softDelete {
		tx = tx.Unscoped()
	}

	ids := make([]int64, 0, len(commands))
	for _, cmd := range commands {
		ids = append(ids, cmd.ID)
	}

	return tx.Where("id in (?)", ids).Delete(core.Command{}).Error
}

func (s *store) ListPending(ctx context.Context, limit int) ([]*core.Command, error) {
	var commands []*core.Command
	err := s.db.View().Where("deleted_at IS NULL").Limit(limit).Find(&commands).Error
	return commands, err
}
