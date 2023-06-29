// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package db

import (
	"time"
)

// func buildNetwokrKey(url string) string {
// 	var sb strings.Builder
// 	sb.WriteString("githubcache")
// 	sb.WriteString(url)
// 	return sb.String()
// }

var networkConfig = "git-network-config-directory"

func RedisSetNetworkConfig(result string) {
	err := rdb.Set(ctxRedis, networkConfig, result, time.Duration(expiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetNetworkConfig() (string, error) {
	val, err := rdb.Get(ctxRedis, networkConfig).Result()
	return formatRedisResponse(val, err)
}

func RedisSetNetworkConfigByName(name string, result string) {
	err := rdb.Set(ctxRedis, networkConfig+name, result, time.Duration(expiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetNetworkConfigByName(name string) (string, error) {
	val, err := rdb.Get(ctxRedis, networkConfig+name).Result()
	return formatRedisResponse(val, err)
}
