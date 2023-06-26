// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package db

import (
	"strings"
	"time"
)

func buildAirtableKey(key string, url string) string {
	var sb strings.Builder
	sb.WriteString(key)
	sb.WriteString(url)
	return sb.String()
}

func RedisGetAirtableRequest(path string) (string, error) {
	val, err := rdb.Get(ctxRedis, path).Result()
	return formatRedisResponse(val, err)
}

func RedisSetAirtableRequest(result string, path string) {
	err := rdb.Set(ctxRedis, path, result, time.Duration(15*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetAirtableFallbackRequest(path string) (string, error) {
	val, err := rdb.Get(ctxRedis, buildAirtableKey("airtableFallback", path)).Result()
	return formatRedisResponse(val, err)
}

func RedisSetAirtableFallbackRequest(result string, path string) {
	err := rdb.Set(ctxRedis, buildAirtableKey("airtableFallback", path), result, time.Duration(oneDayExpiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}
