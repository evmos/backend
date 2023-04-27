// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package constants

var IBCChannels = map[string]map[string]string{
	GRAV: {
		"channel-65": EVMOS,
	},
	EVMOS: {
		"channel-0":  OSMOSIS,
		"channel-8":  GRAV,
		"channel-3":  COSMOS,
		"channel-5":  JUNO,
		"channel-21": AXELAR,
	},
	OSMOSIS: {
		"channel-204": EVMOS,
	},
	COSMOS: {
		"channel-292": EVMOS,
	},
	JUNO: {
		"channel-70": EVMOS,
	},
	AXELAR: {
		"channel-22": EVMOS,
	},
}

var IBCCoins = map[string]map[string]string{
	Gweth: {
		"sourceDenom": "gravity0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		"source":      GRAV,
	},
	Gdai: {
		"sourceDenom": "gravity0x6B175474E89094C44Da98b954EedeAC495271d0F",
		"source":      GRAV,
	},
	Gwbtc: {
		"sourceDenom": "gravity0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599",
		"source":      GRAV,
	},
	Grav: {
		"sourceDenom": "ugraviton",
		"source":      GRAV,
	},
	Gusdt: {
		"sourceDenom": "gravity0xdAC17F958D2ee523a2206206994597C13D831ec7",
		"source":      GRAV,
	},
	Usdcgrv: {
		"sourceDenom": "gravity0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		"source":      GRAV,
	},
	Osmo: {
		"sourceDenom": "uosmo",
		"source":      OSMOSIS,
	},
	Atom: {
		"sourceDenom": "uatom",
		"source":      COSMOS,
	},
	Juno: {
		"sourceDenom": "ujuno",
		"source":      JUNO,
	},
	Axlusdc: {
		"sourceDenom": "uusdc",
		"source":      AXELAR,
	},
	Axlwbtc: {
		"sourceDenom": "wbtc-satoshi",
		"source":      AXELAR,
	},
	Axlweth: {
		"sourceDenom": "weth-wei",
		"source":      AXELAR,
	},
}
