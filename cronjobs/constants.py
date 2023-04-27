# Copyright Tharsis Labs Ltd.(Evmos)
# SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

from __future__ import annotations

import os
from enum import Enum
from typing import Any

EVMOS_PREFIX = os.getenv('EVMOS_PREFIX', 'evmos')


EVMOS = "EVMOS"

INDEXING_DISABLED_ERROR = "transaction indexing is disabled"

class Networks(Enum):
    GRAV = 'GRAVITY'
    EVMOS = 'EVMOS'
    OSMOSIS = 'OSMOSIS'
    COSMOS = 'COSMOS'
    JUNO = 'JUNO'
    AXELAR = 'AXELAR'


TOKENS = {
    'gWETH': 'gWETH',
    'gDAI': 'gDAI',
    'gWBTC': 'gWBTC',
    'GRAV': 'GRAV',
    'USDC.grv': 'USDC.grv',
    'gUSDT': 'gUSDT',
    'OSMO': 'OSMO',
    'ATOM': 'ATOM',
}

IBC_COINS = {
    TOKENS['gWETH']: {
        'sourceDenom': 'gravity0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2',
        'source': Networks.GRAV.value,
    },
    TOKENS['gDAI']: {
        'sourceDenom': 'gravity0x6B175474E89094C44Da98b954EedeAC495271d0F',
        'source': Networks.GRAV.value,
    },
    TOKENS['gWBTC']: {
        'sourceDenom': 'gravity0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599',
        'source': Networks.GRAV.value,
    },
    TOKENS['GRAV']: {
        'sourceDenom': 'ugraviton',
        'source': Networks.GRAV.value,
    },
    TOKENS['gUSDT']: {
        'sourceDenom': 'gravity0xdAC17F958D2ee523a2206206994597C13D831ec7',
        'source': Networks.GRAV.value,
    },
    TOKENS['USDC.grv']: {
        'sourceDenom': 'gravity0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
        'source': Networks.GRAV.value,
    },
    TOKENS['OSMO']: {
        'sourceDenom': 'uosmo',
        'source': Networks.OSMOSIS.value,
    },
    TOKENS['ATOM']: {
        'sourceDenom': 'uatom',
        'source': Networks.COSMOS.value,
    },
}

CHAIN_INFO: dict[str, Any] = {
    Networks.GRAV.value: {
        'rest': [
            'https://gravitychain.io:1317',
            'https://lcd.gravity-bridge.ezstaking.io',
            'https://api-gravitybridge-ia.notional.ventures',
            'https://api.gravity-bridge.nodestake.top',
        ],
        'jrpc': [
            'https://gravitychain.io:26657',
            'http://gravity-bridge-1-08.nodes.amhost.net:26657',
            'https://gravity-rpc.polkachu.com',
            'https://rpc.gravity-bridge.ezstaking.io',
            'https://rpc-gravitybridge-ia.notional.ventures',
            'https://rpc.gravity-bridge.nodestake.top',
        ],
        'prefix':
        'gravity',
        'ibcchannels': {
            'channel-65': Networks.EVMOS.value
        }
    },
    Networks.EVMOS.value: {
        'rest': [
            'https://lcd-evmos.whispernode.com',
            'https://rest.bd.evmos.org:1317',
            'https://lcd.evmos.ezstaking.io',
            'https://api-evmos-ia.cosmosia.notional.ventures',
            'https://lcd.evmos.posthuman.digital',
            'https://api.evmos.interbloc.org',
            # Removed because the gas requirement is bigger than the min value
            # 'https://rest-evmos.ecostake.com',
            'https://lcd.evmos.bh.rocks',
        ],
        'jrpc': [
            'https://rpc-evmos.whispernode.com',
            'https://tendermint.bd.evmos.org:26657',
            'https://rpc.evmos.ezstaking.io',
            'https://rpc-evmos-ia.cosmosia.notional.ventures:443',
            'https://rpc.evmos.posthuman.digital',
            'https://rpc.evmos.interbloc.org',
            'https://rpc.evmos.nodestake.top',
            'https://rpc-evmos.ecostake.com',
            'https://rpc.evmos.bh.rocks',
        ],
        'web3': [
            'https://jsonrpc-rpcaas-evmos-mainnet.ubiquity.blockdaemon.tech',
            'https://eth.bd.evmos.org:8545',
            'https://jsonrpc-evmos-ia.cosmosia.notional.ventures',
            'https://evmos-json-rpc.stakely.io',
            'https://jsonrpc.evmos.nodestake.top',
            'https://json-rpc.evmos.bh.rocks',
        ],
        'prefix':
        EVMOS_PREFIX,
        'ibcchannels': {
            'channel-0': Networks.OSMOSIS.value,
            'channel-8': Networks.GRAV.value,
            'channel-3': Networks.COSMOS.value,
        }
    },
    Networks.OSMOSIS.value: {
        'rest': [
            'https://osmosis-lcd.quickapi.com:443',
            'https://lcd-osmosis.whispernode.com',
            'https://lcd-osmosis.blockapsis.com',
            'https://rest-osmosis.ecostake.com',
            'https://api-osmosis-ia.notional.ventures',
            'https://lcd.osmosis.zone',
            'https://api.osmosis.interbloc.org',
        ],
        'jrpc': [
            'https://osmosis-rpc.quickapi.com:443',
            'https://rpc-osmosis.whispernode.com',
            'https://osmosis.validator.network',
            'https://rpc-osmosis.blockapsis.com',
            'https://rpc-osmosis.ecostake.com',
            'https://osmosis-rpc.polkachu.com',
            'https://rpc-osmosis-ia.notional.ventures',
            'https://rpc.osmosis.zone',
            'https://rpc.osmosis.interbloc.org',
        ],
        'prefix':
        'osmo',
        'ibcchannels': {
            'channel-204': Networks.EVMOS.value,
        }
    },
    Networks.COSMOS.value: {
        'rest': [
            'https://cosmos-lcd.quickapi.com:443',
            'https://lcd-cosmoshub.whispernode.com',
            'https://lcd-cosmoshub.blockapsis.com',
            'https://rest-cosmoshub.ecostake.com',
            'https://api.cosmoshub.pupmos.network',
            'https://lcd.cosmos.ezstaking.io',
            'https://api-cosmoshub-ia.notional.ventures/',
        ],
        'jrpc': [
            'https://cosmos-rpc.quickapi.com:443',
            'https://rpc-cosmoshub.whispernode.com',
            'https://rpc-cosmoshub.blockapsis.com',
            'https://cosmoshub.validator.network/',
            'https://rpc.cosmoshub.strange.love',
            'https://rpc.cosmos.network:443',
            'https://rpc-cosmoshub.ecostake.com',
            'https://rpc.cosmoshub.pupmos.network',
            'https://cosmos-rpc.polkachu.com',
            'https://rpc.cosmos.ezstaking.io',
            'https://rpc-cosmoshub-ia.notional.ventures/',
        ],
        'prefix':
        'cosmos',
        'ibcchannels': {
            'channel-292': Networks.EVMOS.value,
        }
    },
    Networks.JUNO.value: {
        'rest': [
            'https://lcd-juno.itastakers.com',
            'https://rest-juno.ecostake.com',
            'https://juno-api.lavenderfive.com:443',
            'https://api.juno.pupmos.network',
            'https://api-juno-ia.cosmosia.notional.ventures',
            'https://juno-api.polkachu.com',
        ],
        'jrpc': [
            'https://rpc-juno.itastakers.com',
            'https://rpc-juno.ecostake.com',
            'https://juno-rpc.polkachu.com',
            'https://juno-rpc.lavenderfive.com:443',
            'https://rpc-juno-ia.cosmosia.notional.ventures',
            'https://rpc.juno.chaintools.tech',
            'https://rpc.juno.pupmos.network',
        ],
        'prefix':
        'juno',
        'ibcchannels': {
            'channel-70': Networks.EVMOS.value,
        }
    },
    Networks.AXELAR.value: {
        'rest': [
            "https://lcd-axelar.imperator.co:443",
            "https://axelar-lcd.quickapi.com:443",
            "https://axelar-rest.chainode.tech:443",
            "https://axelar-lcd.qubelabs.io:443",
            "https://api-1.axelar.nodes.guru:443",
            "https://api-axelar-ia.cosmosia.notional.ventures/",
            "https://axelar-api.polkachu.com",
        ],
        'jrpc': [
            "https://rpc-axelar.imperator.co:443",
            "https://axelar-rpc.quickapi.com:443",
            "https://axelar-rpc.chainode.tech:443",
            "https://axelar-rpc.pops.one:443",
            "https://axelar-rpc.qubelabs.io:443",
            "https://rpc-1.axelar.nodes.guru:443",
            "https://rpc-axelar-ia.cosmosia.notional.ventures/",
            "https://axelar-rpc.polkachu.com",
        ],
        'prefix':
        'axelar',
        'ibcchannels': {
            'channel-22': Networks.EVMOS.value,
        }
    }
}

ERC20_MODULE_COINS = [
    # Note: aevmos and wevmos are not created using the erc20 module
    {
        'denom': 'aevmos',
        'erc20': '0xD4949664cD82660AaE99bEdc034a0deA8A0bd517',
        'tokenName': 'Evmos',
        'description': 'Evmos native coin',
        'coingeckoId': 'evmos',
    },
    {
        'denom': 'ibc/6B3FCE336C3465D3B72F7EFB4EB92FC521BC480FE9653F627A0BD0237DF213F3',
        'erc20': '0xc03345448969Dd8C00e9E4A85d2d9722d093aF8E',
        'tokenName': TOKENS['gWETH'],
        'description': 'Gravity Bridge WETH',
        'coingeckoId': 'weth',
    },
    {
        'denom': 'ibc/F96A7F81E8F82E4EE81F94D507CD257319EFB70FE46E23B4953F63B62E855603',
        'erc20': '0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75',
        'tokenName': TOKENS['gDAI'],
        'description': 'Gravity Bridge DAI',
        'coingeckoId': 'dai',
    },
    {
        'denom': 'ibc/350B6DC0FF48E3BDB856F40A8259909E484259ED452B3F4F39A0FEF874F30F61',
        'erc20': '0x1D54EcB8583Ca25895c512A8308389fFD581F9c9',
        'tokenName': TOKENS['gWBTC'],
        'description': 'Gravity Bridge BTC',
        'coingeckoId': 'wrapped-bitcoin',
    },
    {
        'denom': 'ibc/7F0C2CB6E79CC36D29DA7592899F98E3BEFD2CF77A94340C317032A78812393D',
        'erc20': '0x80b5a32E4F032B2a058b4F29EC95EEfEEB87aDcd',
        'tokenName': TOKENS['GRAV'],
        'description': 'Gravity Bridge native coin',
        'coingeckoId': 'graviton',
    },
    {
        'denom': 'ibc/DF63978F803A2E27CA5CC9B7631654CCF0BBC788B3B7F0A10200508E37C70992',
        'erc20': '0xecEEEfCEE421D8062EF8d6b4D814efe4dc898265',
        'tokenName': TOKENS['gUSDT'],
        'description': 'Gravity Bridge USDT',
        'coingeckoId': 'tether',
    },
    {
        'denom': 'ibc/693989F95CF3279ADC113A6EF21B02A51EC054C95A9083F2E290126668149433',
        'erc20': '0x5FD55A1B9FC24967C4dB09C513C3BA0DFa7FF687',
        'tokenName': TOKENS['USDC.grv'],
        'description': 'Gravity Bridge USDC',
        'coingeckoId': 'usd-coin',
    },
    {
        'denom': 'ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518',
        'erc20': '0xFA3C22C069B9556A4B2f7EcE1Ee3B467909f4864',
        'tokenName': TOKENS['OSMO'],
        'description': 'The native token of Osmosis',
        'coingeckoId': 'osmosis',
    },
    {
        'denom': 'ibc/A4DB47A9D3CF9A068D454513891B526702455D3EF08FB9EB558C561F9DC2B701',
        'erc20': '0xC5e00D3b04563950941f7137B5AfA3a534F0D6d6',
        'tokenName': TOKENS['ATOM'],
        'description': 'The native token of Cosmos Hub',
        'coingeckoId': 'cosmos',
    },
]
