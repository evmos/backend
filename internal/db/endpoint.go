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
	return formatRedisResponse(val, err)
}

//func GetRestEndpointFromRedis(chain, endpoint string) (string, error) {
//	i := 1
//    var (
//      endpoint string
//      err error
//    )
//	for i < 4 {
//		endpoint, err = RedisGetEndpoint(chain, "rest", strconv.FormatInt(int64(i), 10))
//		if err != nil {
//			i++
//			continue
//		}
//
//	}
//}

func RedisSetEndpoint(chain, endpoint, index, url string) {
	key := buildKeyEndpoint(chain, endpoint, index)
	err := rdb.Set(ctxRedis, key, url, 0).Err()
	if err != nil {
		panic(err)
	}
}
