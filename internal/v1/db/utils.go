// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package db

import (
	"strings"
	"time"
)

func buildKeyChainHeight(chain string) string {
	var sb strings.Builder
	sb.WriteString("CHAINHEIGHT")
	sb.WriteString(chain)
	return sb.String()
}

func RedisSetChainHeight(chain string, response string) {
	key := buildKeyChainHeight(chain)
	err := rdb.Set(ctxRedis, key, response, time.Duration(expiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetChainHeight(chain string) (string, error) {
	key := buildKeyChainHeight(chain)
	val, err := rdb.Get(ctxRedis, key).Result()
	return formatRedisResponse(val, err)
}
