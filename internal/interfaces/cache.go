package interfaces

import (
	"time"
)

type CacheClient interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	IsConnected() bool
	GetNameFromPath(path string) (string, error)
	GetNameAndKeyFromPath(path string) (string, string, error)
	GetDefaultExpiration() time.Duration
}
