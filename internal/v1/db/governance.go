// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package db

import (
	"time"
)

var proposalsKey = "governance-props"

func RedisSetGovernanceProposals(proposals string) {
	err := rdb.Set(ctxRedis, proposalsKey+"-"+"v1beta1", proposals, time.Duration(60*15*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetGovernanceProposals() (string, error) {
	val, err := rdb.Get(ctxRedis, proposalsKey+"-"+"v1beta1").Result()
	return formatRedisResponse(val, err)
}

func RedisSetGovernanceV1Proposals(proposals string) {
	err := rdb.Set(ctxRedis, proposalsKey+"-"+"v1", proposals, time.Duration(60*15*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetGovernanceV1Proposals() (string, error) {
	val, err := rdb.Get(ctxRedis, proposalsKey+"-"+"v1").Result()
	return formatRedisResponse(val, err)
}
