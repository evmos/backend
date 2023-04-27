# Copyright Tharsis Labs Ltd.(Evmos)
# SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

from constants import EVMOS


def get_mainnet_config(configuation):
    return next(item for item in configuation if item["configurationType"] == "mainnet")


def get_networks(data):
    return [chain.get('prefix').upper() for chain in data]


def get_chains_info(data):
    chains_info = {}
    destinationChannels = {}
    for chain in data:
        mainnet_config = get_mainnet_config(chain.get('configurations'))
        _chain_info = {
            "prefix": chain.get('prefix'),
            "ibcchannels": {
                f"{mainnet_config.get('source').get('sourceChannel')}": EVMOS
            },
            "rest": mainnet_config.get('rest') if not isinstance(mainnet_config.get('rest'), str) else [],
            "jrpc": mainnet_config.get('jrpc', []),
            "web3": mainnet_config.get('web3', [])
        }
        if mainnet_config.get('identifier', '').upper() != 'EVMOS':
            destinationChannels[f"{mainnet_config.get('source').get('destinationChannel')}"] = mainnet_config.get(
                'identifier', '').upper()
        chains_info[mainnet_config.get('identifier', '').upper()] = _chain_info

    evmos_chain_info = chains_info.get('EVMOS')
    evmos_chain_info["ibcchannels"] = destinationChannels

    return chains_info


def get_erc20_coins(data):
    erc20_coins = []

    for coin in data:
        _coin = {
            "denom": coin.get("cosmosDenom"),
            "erc20": coin.get("erc20Address"),
            "tokenName": coin.get("coinDenom"),
            "description": coin.get("description"),
            "coingeckoId": coin.get("coingeckoId")
        }
        erc20_coins.append(_coin)
    return erc20_coins
