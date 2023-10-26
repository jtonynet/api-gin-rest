package redis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
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
	val, err := c.db.Get(c.ctx, key).Result()
	if err != nil {
		slog.Error("Cannot get key: %v, CacheClient error: %v ", key, err)
		return "", err
	}
	if val == "" {
		return "", errors.New("Get data empty")
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

func (c *Client) GetNameFromPath(path string) (string, error) {
	segments := strings.Split(path, "/")
	if len(segments) > 0 {
		return segments[len(segments)-1], nil
	}
	return "", errors.New("improperly formatted path.")
}

func (c *Client) GetNameAndKeyFromPath(path string) (string, string, error) {
	// TODO: Criar um metodo que receba o path e a queryString e devolva
	// a cacheKey adequada para Setar ou obter o dado, mas isso deve estar
	// no Middleware e nao na lib de cacheClient pois depende de parse
	// do gin

	var paramName, paramKey string = "", ""

	pathSplited := strings.Split(path, "/")
	if len(pathSplited) > 2 {
		paramName = pathSplited[len(pathSplited)-2]
	}

	paramSplited := strings.Split(path, ":")
	if len(paramSplited) > 1 {
		paramKey = paramSplited[1]
	}

	if paramName == "" || paramKey == "" {
		return paramName, paramKey, errors.New("improperly formatted path.")
	}

	return paramName, paramKey, nil
}
