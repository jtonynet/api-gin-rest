package handlers

import (
	"fmt"
	"log/slog"

	"github.com/jtonynet/api-gin-rest/internal/interfaces"
	"github.com/tidwall/gjson"
)

type InsertAlunoCached struct {
	next            iInsert
	cacheClient     interfaces.CacheClient
	cacheKeySegment string
}

func NewInsertAlunoCached(
	next iInsert,
	cacheClient interfaces.CacheClient,
	cacheKeySegment string,
) iInsert {
	return &InsertAlunoCached{
		next:            next,
		cacheClient:     cacheClient,
		cacheKeySegment: cacheKeySegment,
	}
}

func (i *InsertAlunoCached) InsertMethod(msg string) (string, error) {
	msgValue, err := i.next.InsertMethod(msg)
	if err != nil {
		return msgValue, err
	}

	if i.cacheClient != nil {
		msgUUID := gjson.Get(msgValue, "uuid")
		msgKey := fmt.Sprintf("%s:%s", i.cacheKeySegment, msgUUID)

		i.cacheClient.Delete(msgKey)
		err = i.cacheClient.Set(msgKey, string(msgValue), i.cacheClient.GetDefaultExpiration())
		if err != nil {
			slog.Error("cmd:worker:main:messageBroker:RunConsumer:handle:b.cacheClient.Set error set%v", err)
			return msgValue, err
		}
	}

	return msgValue, nil
}
