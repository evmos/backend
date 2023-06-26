package db

import (
	"strings"
)

func buildKeyEndpoint(chain, endpoint, index string) string {
	var sb strings.Builder
	sb.WriteString(strings.ToUpper(chain))
	sb.WriteString("|")
	sb.WriteString(endpoint)
	sb.WriteString("|")
	sb.WriteString(index)
	return sb.String()
}

func RedisGetEndpoint(chain, endpoint, index string) (string, error) {
	key := buildKeyEndpoint(chain, endpoint, index)
	val, err := rdb.Get(ctxRedis, key).Result()
	return formatRedisResponse(val, err)
}

func RedisSetEndpoint(chain, endpoint, index, url string) {
	key := buildKeyEndpoint(chain, endpoint, index)
	err := rdb.Set(ctxRedis, key, url, 0).Err()
	if err != nil {
		panic(err)
	}
}
