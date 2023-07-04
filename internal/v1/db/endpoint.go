// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)
package db

import (
	"context"
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
	ctx := context.Background()
	match := buildKeyEndpoint(chain, serverType, "*")
	iter := rdb.Scan(ctx, 0, match, 0).Iterator()
	var nodes []string
	for iter.Next(ctx) {
		rd, err := rdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, rd)
	}

	if err := iter.Err(); err != nil {
		return nil, err
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
