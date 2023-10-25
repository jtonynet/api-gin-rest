package interfaces

import "time"

type CacheClient interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	IsConnected() bool
	GetDefaultExpiration() time.Duration
}
