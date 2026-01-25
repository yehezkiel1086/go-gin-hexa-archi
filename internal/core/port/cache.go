package port

import (
	"context"
	"time"
)

//go:generate mockery --name=CacheRepository --output=../../../mocks --outpkg=mocks
type CacheRepository interface {
	Set(ctx context.Context, key string, val []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	DeleteByPrefix(ctx context.Context, prefix string) error
	Close() error
}
