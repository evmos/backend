// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package db

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// erc20TokensDirectoryKey represents the Redis key for the ERC20 tokens directory
var erc20TokensDirectoryKey string

func init() {
	erc20TokensDirectoryKey = getErc20TokensDirectoryKey()
}

func getErc20TokensDirectoryKey() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "production" {
		return "prod-git-erc20-tokens-directory"
	}
	return "git-erc20-tokens-directory"
}

func getErc20TokensDirectoryKeyByName(name string) string {
	return fmt.Sprintf("%s-%s", erc20TokensDirectoryKey, name)
}

func buildKeyERC20Balance(chain string, contract string, address string) string {
	var sb strings.Builder
	sb.WriteString("ERC20")
	sb.WriteString(chain)
	sb.WriteString(contract)
	sb.WriteString(address)
	return sb.String()
}

func RedisSetERC20Balance(contract string, address string, balance string) {
	key := buildKeyERC20Balance("EVMOS", contract, address)
	err := rdb.Set(ctxRedis, key, balance, time.Duration(expiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetERC20Balance(contract string, address string) (string, error) {
	key := buildKeyERC20Balance("EVMOS", contract, address)
	val, err := rdb.Get(ctxRedis, key).Result()
	return formatRedisResponse(val, err)
}

func RedisSetERC20TokensDirectory(result string) {
	err := rdb.Set(ctxRedis, erc20TokensDirectoryKey, result, time.Duration(oneDayExpiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetERC20TokensDirectory() (string, error) {
	val, err := rdb.Get(ctxRedis, erc20TokensDirectoryKey).Result()
	return formatRedisResponse(val, err)
}

func RedisSetERC20TokensByName(name string, result string) {
	key := getErc20TokensDirectoryKeyByName(name)
	err := rdb.Set(ctxRedis, key, result, time.Duration(expiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetERC20TokensByName(name string) (string, error) {
	key := getErc20TokensDirectoryKeyByName(name)
	val, err := rdb.Get(ctxRedis, key).Result()
	return formatRedisResponse(val, err)
}
