# Copyright Tharsis Labs Ltd.(Evmos)
# SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

from __future__ import annotations
import json

import os

import redis

REDIS_HOST = os.getenv('REDIS_HOST', None)
REDIS_PORT = os.getenv('REDIS_PORT', 6379)
SECONDS_PER_HOUR = 3600

ENVIRONMENT = os.getenv("ENVIRONMENT")


if REDIS_HOST:
    r = redis.Redis(host=REDIS_HOST, port=int(
        REDIS_PORT), decode_responses=True)
else:
    r = redis.Redis(decode_responses=True)


prod_prefix = "prod-" if ENVIRONMENT == "production" else ""

erc20TokensDirectoryKey = f"{prod_prefix}git-erc20-tokens-directory"
networkConfig = f"{prod_prefix}git-network-config-directory"


def redisSetPrice(asset: str, vs_currency: str, price: float):
    key = f'{asset}|{vs_currency}|price'
    r.mset({key: price})

def redisSetEvmosChange(change: float):
    key = f'evmos|24h|change'
    r.mset({key: change})

def redisGetPrice(asset: str, vs_currency: str) -> float | None:
    key = f'{asset}|{vs_currency}|price'
    value = r.get(key)
    if value is None:
        return None
    return float(value)


def redisSetEndpoint(chain: str, endpoint: str, order: int, url: str):
    key = f'{chain}|{endpoint}|{order}'
    r.mset({key: url})


def redisGetEndpoint(chain: str, endpoint: str, order: int) -> str | None:
    key = f'{chain}|{endpoint}|{order}'
    value = r.get(key)
    if not value:
        return None
    return str(value)


def setPrimaryEndpoint(chain: str, endpoint: str, url: str):
    return redisSetEndpoint(chain, endpoint, 1, url)


def getPrimaryEndpoint(chain: str, endpoint: str) -> str | None:
    return redisGetEndpoint(chain, endpoint, 1)


def setSecondaryEndpoint(chain: str, endpoint: str, url: str):
    return redisSetEndpoint(chain, endpoint, 2, url)


def getSecondaryEndpoint(chain: str, endpoint: str) -> str | None:
    return redisGetEndpoint(chain, endpoint, 2)


def setTertiaryEndpoint(chain: str, endpoint: str, url: str):
    return redisSetEndpoint(chain, endpoint, 3, url)


def getTertiaryEndpoint(chain: str, endpoint: str) -> str | None:
    return redisGetEndpoint(chain, endpoint, 3)


def setTokens(data):
    r.set(erc20TokensDirectoryKey, json.dumps(data), SECONDS_PER_HOUR*24)


def getTokens():
    value = r.get(erc20TokensDirectoryKey)
    if not value:
        return None
    return json.loads(value)


def setChains(data):
    r.set(networkConfig, json.dumps(data), SECONDS_PER_HOUR*24)


def getChains():
    value = r.get(networkConfig)
    if not value:
        return None
    return json.loads(value)
