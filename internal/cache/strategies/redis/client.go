package redis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jtonynet/api-gin-rest/config"
	redis "github.com/redis/go-redis/v9"
)

/*
Fortemente baseado nos exemplos da lib go-redis
https://github.com/redis/go-redis
*/

type Client struct {
	db  *redis.Client
	ctx context.Context
	cfg config.Cache

	Expiration time.Duration
}

func NewClient(cfg config.Cache) (*Client, error) {
	strAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	db := redis.NewClient(&redis.Options{
		Addr:     strAddr,
		Password: cfg.Pass,
		DB:       cfg.DB,
		Protocol: cfg.Protocol,
	})

	Expiration := time.Duration(cfg.Expiration * int(time.Millisecond))

	return &Client{
		db:  db,
		ctx: context.Background(),
		cfg: cfg,

		Expiration: Expiration,
	}, nil
}

func (c *Client) Set(key string, value interface{}, expiration time.Duration) error {
	err := c.db.Set(c.ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Get(key string) (string, error) {
	slog.Info("++++++++++++++++++CONSULTANDO A CHAVE: ", key)
	val, err := c.db.Get(c.ctx, key).Result()
	if err != nil {
		slog.Error("Cannot get key: %v, CacheClient error: %v ", key, err)
		return "", err
	}
	if val == "" {
		slog.Error("-------------VALOR VAZIO: ", err)
		return "", errors.New("data empty")
	}

	return val, nil
}

func (c *Client) SetWithCtx(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := c.db.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetWithCtx(ctx context.Context, key string) (string, error) {
	slog.Info("++++++++++++++++++CONSULTANDO A CHAVE: ", key)
	val, err := c.db.Get(c.ctx, key).Result()
	if err != nil {
		slog.Error("Cannot get key: %v, CacheClient error: %v ", key, err)
		return "", err
	}
	if val == "" {
		slog.Error("-------------VALOR VAZIO: ", err)
		return "", errors.New("data empty")
	}

	return val, nil
}

func (c *Client) Delete(key string) error {
	err := c.db.Del(c.ctx, key).Err()
	if err != nil {
		slog.Error("Cannot delete key: %v, CacheClient error: %v", key, err)
		return err
	}
	return nil
}

func (c *Client) IsConnected() bool {
	_, err := c.db.Ping(c.ctx).Result()
	return err == nil
}

func (c *Client) GetDefaultExpiration() time.Duration {
	return c.Expiration
}
