package interfaces

import (
	"context"
	"time"
)

type CacheClient interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	SetWithCtx(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetWithCtx(ctx context.Context, key string) (string, error)
	Delete(key string) error
	IsConnected() bool
	GetDefaultExpiration() time.Duration
}
