// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package db

import (
	"time"
)

var (
	oneDayExpiration  = 60 * 60 * 24
	twoDaysExpiration = 60 * 60 * 48
)

var validatorDirectoryKey = "validator-directory"

func RedisSetValidatorDirectory(result string) {
	err := rdb.Set(ctxRedis, validatorDirectoryKey, result, time.Duration(oneDayExpiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetValidatorDirectory() (string, error) {
	val, err := rdb.Get(ctxRedis, validatorDirectoryKey).Result()
	return formatRedisResponse(val, err)
}

func RedisSetValidatorDirectoryNoListed(status string, sort string, result string) {
	err := rdb.Set(ctxRedis, validatorDirectoryKey+status+sort, result, time.Duration(expiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetValidatorDirectoryNoListed(status string, sort string) (string, error) {
	val, err := rdb.Get(ctxRedis, validatorDirectoryKey+status+sort).Result()
	return formatRedisResponse(val, err)
}
