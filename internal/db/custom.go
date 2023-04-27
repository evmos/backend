// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package db

import "time"

var (
	allValidators                   = "allValidators"
	validatorWithNoFilterKey        = "valWithNoFilter"
	validatorWithRanks              = "valWithRanks"
	delegationsWithRanks            = "delWithRanks"
	byAddrWithRanks                 = "byAddrWithRanks"
	expirationValidatorWithNoFilter = 60
)

func RedisSetValidatorWithNoFilter(chain string, result string) {
	// Using 1 hour as cache, creating this object takes too much time and the api response will not change very often
	err := rdb.Set(ctxRedis, chain+validatorWithNoFilterKey, result, time.Duration(expirationValidatorWithNoFilter*int(time.Minute))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetValidatorWithNoFilter(chain string) (string, error) {
	val, err := rdb.Get(ctxRedis, chain+validatorWithNoFilterKey).Result()
	return formatRedisResponse(val, err)
}

func RedisSetUnbondingByAddressWithValidatorInfo(chain string, address string, result string) {
	err := rdb.Set(ctxRedis, chain+address+validatorWithNoFilterKey, result, time.Duration(expiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetUnbondingByAddressWithValidatorInfo(chain string, address string) (string, error) {
	val, err := rdb.Get(ctxRedis, chain+address+validatorWithNoFilterKey).Result()
	return formatRedisResponse(val, err)
}

func RedisSetValidatorWithRanks(chain string, result string) {
	// Using 1 hour as cache, creating this object takes too much time and the api response will not change very often
	err := rdb.Set(ctxRedis, chain+validatorWithRanks, result, time.Duration(expirationValidatorWithNoFilter*int(time.Minute))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetValidatorWithRanks(chain string) (string, error) {
	val, err := rdb.Get(ctxRedis, chain+validatorWithRanks).Result()
	return formatRedisResponse(val, err)
}

func RedisSetDelegationsByAddressWithValidatorRanks(chain string, address string, result string) {
	err := rdb.Set(ctxRedis, chain+address+delegationsWithRanks, result, time.Duration(expiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetDelegationsByAddressWithValidatorRanks(chain string, address string) (string, error) {
	val, err := rdb.Get(ctxRedis, chain+address+delegationsWithRanks).Result()
	return formatRedisResponse(val, err)
}

func RedisSetValidatorsByAddressWithValidatorRanks(chain string, address string, result string) {
	err := rdb.Set(ctxRedis, chain+address+byAddrWithRanks, result, time.Duration(expiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetValidatorsByAddressWithValidatorRanks(chain string, address string) (string, error) {
	val, err := rdb.Get(ctxRedis, chain+address+byAddrWithRanks).Result()
	return formatRedisResponse(val, err)
}

func RedisSetAllValidators(chain string, result string) {
	err := rdb.Set(ctxRedis, chain+allValidators, result, time.Duration(expiration*int(time.Second))).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGetAllValidators(chain string) (string, error) {
	val, err := rdb.Get(ctxRedis, chain+allValidators).Result()
	return formatRedisResponse(val, err)
}
