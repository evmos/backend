// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package db

import (
	"strings"
)

func buildKeyPrice(asset string, vsCurrency string) string {
	var sb strings.Builder
	sb.WriteString(asset)
	sb.WriteString("|")
	sb.WriteString(vsCurrency)
	sb.WriteString("|price")
	return sb.String()
}

func RedisGetPrice(asset string, vsCurrency string) (string, error) {
	key := buildKeyPrice(asset, vsCurrency)
	val, err := rdb.Get(ctxRedis, key).Result()
	return formatRedisResponse(val, err)
}

func RedisGet24HChange(asset string) (string, error) {
	val, err := rdb.Get(ctxRedis, asset+"|24h|change").Result()
	return formatRedisResponse(val, err)
}

func RedisSetPrice(asset string, vsCurrency string, price string) {
	key := buildKeyPrice(asset, vsCurrency)
	err := rdb.Set(ctxRedis, key, price, 0).Err()
	if err != nil {
		panic(err)
	}
}
