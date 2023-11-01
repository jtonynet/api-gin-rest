package decorators

import (
	"fmt"
	"log/slog"

	"github.com/jtonynet/api-gin-rest/internal/interfaces"
	"github.com/tidwall/gjson"
)

type Cached struct {
	next            interfaces.Handler
	cacheClient     interfaces.CacheClient
	cacheKeySegment string
}

func NewCached(
	next interfaces.Handler,
	cacheClient interfaces.CacheClient,
	cacheKeySegment string,
) interfaces.Handler {
	return &Cached{
		next:            next,
		cacheClient:     cacheClient,
		cacheKeySegment: cacheKeySegment,
	}
}

func (c *Cached) Execute(msg string) (string, error) {
	msgValue, err := c.next.Execute(msg)

	if err != nil {
		return msgValue, err
	}

	if c.cacheClient != nil && c.cacheClient.IsConnected() {
		msgUUID := gjson.Get(msgValue, "uuid")
		msgKey := fmt.Sprintf("%s:%s", c.cacheKeySegment, msgUUID)

		c.cacheClient.Delete(msgKey)
		err = c.cacheClient.Set(msgKey, string(msgValue), c.cacheClient.GetDefaultExpiration())
		if err != nil {
			slog.Error("cmd:worker:main:messageBroker:RunConsumer:handle:b.cacheClient.Set error set%v", err)
			return msgValue, err
		}
	}

	return msgValue, nil
}
