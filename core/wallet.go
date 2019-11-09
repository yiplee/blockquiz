package core

import (
	"context"
)

type Transfer struct {
	TraceID    string `gorm:"size:36;PRIMARY_KEY" json:"trace_id,omitempty"`
	OpponentID string `gorm:"size:36" json:"opponent_id,omitempty"`
	AssetID    string `gorm:"size:36" json:"asset_id,omitempty"`
	Amount     string `gorm:"size:64" json:"amount,omitempty"`
	Memo       string `gorm:"size:256" json:"memo,omitempty"`
}

type WalletStore interface {
	Create(ctx context.Context, transfer *Transfer) error
	Deletes(ctx context.Context, traceIDs []string) error
	List(ctx context.Context, limit int) ([]*Transfer, error)
}

type WalletService interface {
	Transfer(ctx context.Context, req *Transfer) error
}
