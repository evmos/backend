# Copyright Tharsis Labs Ltd.(Evmos)
# SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

import signal
import time

import requests

from github import get_tokens
from helpers import get_erc20_coins
from redis_functions import redisSetPrice, redisSetEvmosChange


def get_evmos_change():
    try:
        url = 'https://api.coingecko.com/api/v3/coins/evmos'
        resp = requests.get(f'{url}')
        json_resp = resp.json()
        redisSetEvmosChange(json_resp["market_data"]["price_change_percentage_24h"])
        return
    except Exception:
        return None

def get_prices(vs_currency: str, erc20_module_coins):
    asset_ids = []

    for coin in erc20_module_coins:
        if (coin["coingeckoId"] and coin["coingeckoId"] != ""):
            asset_ids.append(coin["coingeckoId"])

    delim = ","
    try:
        url = 'https://api.coingecko.com/api/v3/simple/price?'
        resp = requests.get(f'{url}ids={delim.join(asset_ids)}&vs_currencies={vs_currency}')
        return resp.json()
    except Exception:
        return None


def process_assets(prices):
    for asset in prices:
        price = prices[asset].get('usd', None)
        if price is not None:
            redisSetPrice(asset, 'usd', price)
            print(f'Price {price} for {asset}')

running = True


def main():
    global running
    while running:
        tracked_tokens = get_tokens()
        erc20_module_coins = get_erc20_coins(tracked_tokens)
        print('Getting prices...')
        prices = get_prices("usd", erc20_module_coins)
        get_evmos_change()
        process_assets(prices)
        time.sleep(300)


def signal_handler(sig, frame):
    global running
    _ = sig
    _ = frame
    print('You pressed Ctrl+C!')
    running = False


signal.signal(signal.SIGINT, signal_handler)

if __name__ == '__main__':
    main()
