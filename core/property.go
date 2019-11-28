package core

import (
	"context"

	"github.com/yiplee/blockquiz/property"
)

type PropertyStore interface {
	Get(ctx context.Context, key string) (property.Value, error)
	Save(ctx context.Context, key string, value interface{}) error
	List(ctx context.Context) (map[string]property.Value, error)
}
