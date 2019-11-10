package wallet

import (
	"context"

	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/db"
)

type store struct {
	db *db.DB
}

func New(db *db.DB) core.WalletStore {
	return &store{db: db}
}

func (s *store) CreateTransfer(ctx context.Context, transfer *core.Transfer) error {
	return s.db.Update().FirstOrCreate(transfer).Error
}

func (s *store) DeleteTransfers(ctx context.Context, traceIDs []string) error {
	return s.db.Update().Where("trace_id IN (?)", traceIDs).Delete(core.Transfer{}).Error
}

func (s *store) ListTransfers(ctx context.Context, limit int) ([]*core.Transfer, error) {
	var transfers []*core.Transfer
	err := s.db.View().Limit(limit).Find(&transfers).Error
	return transfers, err
}

func (s *store) CountTransfers(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.View().Model(core.Transfer{}).Count(&count).Error
	return count, err
}
