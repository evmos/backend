# Copyright Tharsis Labs Ltd.(Evmos)
# SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

import signal
import time

import requests

from github import get_tokens
from helpers import get_erc20_coins
from redis_functions import redisSetPrice


def get_price(asset: str, vs_currency: str) -> float | None:
    try:
        url = 'https://api.coingecko.com/api/v3/simple/price?'
        resp = requests.get(f'{url}ids={asset}&vs_currencies={vs_currency}')
        print(resp)
        return float(resp.json()[asset][vs_currency])
    except Exception:
        return None


def process_assets(erc20_module_coins):
    for coin in erc20_module_coins:
        print(f'Getting price for {coin["tokenName"]}')
        price = get_price(coin['coingeckoId'], 'usd')
        if price is not None:
            redisSetPrice(coin['coingeckoId'], 'usd', price)
            print(f'Price {price} for {coin["tokenName"]}')
        time.sleep(2)


running = True


def main():
    global running
    while running:
        tracked_tokens = get_tokens()
        erc20_module_coins = get_erc20_coins(tracked_tokens)
        print('Getting prices...')
        process_assets(erc20_module_coins)
        time.sleep(5)


def signal_handler(sig, frame):
    global running
    _ = sig
    _ = frame
    print('You pressed Ctrl+C!')
    running = False


signal.signal(signal.SIGINT, signal_handler)

if __name__ == '__main__':
    main()
