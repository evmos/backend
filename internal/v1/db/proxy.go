// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package db

import (
	"strings"
	"time"
)

func buildKeyProxy(chain string, endpoint string, index string) string {
	var sb strings.Builder
	sb.WriteString(chain)
	sb.WriteString(endpoint)
	sb.WriteString(index)
	return sb.String()
}

func RedisSetProxyResponse(chain string, url string, response string) {
	key := buildKeyProxy("proxy", chain, url)
	err := rdb.Set(ctxRedis, key, response, time.Duration(expiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetProxyResponse(chain string, url string) (string, error) {
	key := buildKeyProxy("proxy", chain, url)
	val, err := rdb.Get(ctxRedis, key).Result()
	return formatRedisResponse(val, err)
}

func RedisSetFallbacResponse(chain string, url string, response string) {
	key := buildKeyProxy("fallback", chain, url)
	err := rdb.Set(ctxRedis, key, response, time.Duration(expiration*int(time.Minute))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetFallbackResponse(chain string, url string) (string, error) {
	key := buildKeyProxy("fallback", chain, url)
	val, err := rdb.Get(ctxRedis, key).Result()
	return formatRedisResponse(val, err)
}
