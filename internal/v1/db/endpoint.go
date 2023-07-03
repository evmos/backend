// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)
package db

import (
	"strings"
)

func buildKeyEndpoint(chain, endpoint, index string) string {
	var sb strings.Builder
	sb.WriteString(strings.ToUpper(chain))
	sb.WriteString("|")
	sb.WriteString(endpoint)
	sb.WriteString("|")
	sb.WriteString(index)
	return sb.String()
}

func RedisGetEndpoint(chain, endpoint, index string) (string, error) {
	key := buildKeyEndpoint(chain, endpoint, index)
	val, err := rdb.Get(ctxRedis, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func RedisGetEndpoints(chain, serverType string) ([]string, error) {
	key := buildKeyEndpoint(chain, serverType, "*")
	keys, _, err := rdb.Scan(ctxRedis, 0, key, int64(0)).Result()
	if err != nil {
		return nil, err
	}

	nodes := make([]string, len(keys))
	for _, key := range keys {
		rd, err := rdb.Get(ctxRedis, key).Result()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, rd)
	}
	return nodes, nil
}

func RedisSetEndpoint(chain, endpoint, index, url string) {
	key := buildKeyEndpoint(chain, endpoint, index)
	err := rdb.Set(ctxRedis, key, url, 0).Err()
	if err != nil {
		panic(err)
	}
}
