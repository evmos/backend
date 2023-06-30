// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package db

import (
	"os"
	"time"
)

// func buildNetwokrKey(url string) string {
// 	var sb strings.Builder
// 	sb.WriteString("githubcache")
// 	sb.WriteString(url)
// 	return sb.String()
// }

func getNetworkConfigKey() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "production" {
		return "prod-git-network-config-directory"
	}
	return "git-network-config-directory"
}

func RedisSetNetworkConfig(result string) {
	err := rdb.Set(ctxRedis, getNetworkConfigKey(), result, time.Duration(expiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetNetworkConfig() (string, error) {
	val, err := rdb.Get(ctxRedis, getNetworkConfigKey()).Result()
	return formatRedisResponse(val, err)
}

func RedisSetNetworkConfigByName(name string, result string) {
	err := rdb.Set(ctxRedis, getNetworkConfigKey()+name, result, time.Duration(expiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetNetworkConfigByName(name string) (string, error) {
	val, err := rdb.Get(ctxRedis, getNetworkConfigKey()+name).Result()
	return formatRedisResponse(val, err)
}
