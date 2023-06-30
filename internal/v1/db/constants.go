// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package db

import (
	"context"
	"os"

	"github.com/go-redis/redis/v9"
)

func getRedisHost() *redis.Client {
	host, set := os.LookupEnv("REDIS_HOST")
	if !set {
		host = "localhost"
	}

	return redis.NewClient(&redis.Options{
		Addr:     host + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

var (
	ctxRedis = context.Background()
	rdb      = getRedisHost()
)

var expiration = 7

func formatRedisResponse(val string, err error) (string, error) {
	switch err {
	case redis.Nil:
		return "", err
	case nil:
		return val, nil
	default:
		return "", err
	}
}
