package command

import (
	"context"
	"time"

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

type Command struct {
	ID        int64      `gorm:"PRIMARY_KEY" json:"id,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	*core.Command
}

func (s *store) Create(ctx context.Context, command *core.Command) error {
	cmd := Command{
		Command: command,
	}

	return s.db.Update().Unscoped().Where("trace_id = ?", command.TraceID).FirstOrCreate(&cmd).Error
}

func (s *store) Delete(ctx context.Context, command *core.Command) error {
	return s.Deletes(ctx, []*core.Command{command})
}

func (s *store) Deletes(ctx context.Context, commands []*core.Command) error {
	tx := s.db.Update()
	if !s.softDelete {
		tx = tx.Unscoped()
	}

	traceIDs := make([]string, 0, len(commands))
	for _, cmd := range commands {
		traceIDs = append(traceIDs, cmd.TraceID)
	}

	return tx.Where("trace_id in (?)", traceIDs).Delete(Command{}).Error
}

func (s *store) ListPending(ctx context.Context, limit int) ([]*core.Command, error) {
	var commands []*core.Command
	err := s.db.View().Limit(limit).Find(&commands).Error
	return commands, err
}
