// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package db

import (
	"fmt"
	"os"
	"time"
)

// networkConfigKey represents the Redis key for the network config
var networkConfigKey string

func init() {
	networkConfigKey = getNetworkConfigKey()
}

func getNetworkConfigKey() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "production" {
		return "prod-git-network-config-directory"
	}
	return "git-network-config-directory"
}

func getNetworkConfigKeyByName(name string) string {
	return fmt.Sprintf("%s-%s", networkConfigKey, name)
}

func RedisSetNetworkConfig(result string) {
	err := rdb.Set(ctxRedis, networkConfigKey, result, time.Duration(expiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetNetworkConfig() (string, error) {
	val, err := rdb.Get(ctxRedis, networkConfigKey).Result()
	return formatRedisResponse(val, err)
}

func RedisSetNetworkConfigByName(name string, result string) {
	key := getNetworkConfigKeyByName(name)
	err := rdb.Set(ctxRedis, key, result, time.Duration(expiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetNetworkConfigByName(name string) (string, error) {
	key := getNetworkConfigKeyByName(name)
	val, err := rdb.Get(ctxRedis, key).Result()
	return formatRedisResponse(val, err)
}
