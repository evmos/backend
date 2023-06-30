// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package db

import (
	"strings"
	"time"
)

func buildGithubKey(id string, url string) string {
	var sb strings.Builder
	sb.WriteString(id)
	sb.WriteString(url)
	return sb.String()
}

func RedisSetGithubResponse(url string, result string) {
	err := rdb.Set(ctxRedis, buildGithubKey("githubcache", url), result, time.Duration(oneDayExpiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetGithubResponse(url string) (string, error) {
	val, err := rdb.Get(ctxRedis, buildGithubKey("githubcache", url)).Result()
	return formatRedisResponse(val, err)
}

func RedisSetGithubFallbackResponse(url string, result string) {
	err := rdb.Set(ctxRedis, buildGithubKey("githubcachefallback", url), result, time.Duration(twoDaysExpiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetHithubFallbackResponse(url string) (string, error) {
	val, err := rdb.Get(ctxRedis, buildGithubKey("githubcachefallback", url)).Result()
	return formatRedisResponse(val, err)
}
