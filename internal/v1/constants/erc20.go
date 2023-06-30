// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package constants

// ERC20Tokens
const (
	Gweth   = "gweth"
	Gdai    = "gdai"
	Gwbtc   = "gwbtc"
	Grav    = "grav"
	Usdcgrv = "usdc.grv"
	Gusdt   = "gusdt"
	Osmo    = "osmo"
	Atom    = "atom"
	Juno    = "juno"
	Axlusdc = "axlusdc"
	Axlwbtc = "axlwbtc"
	Axlweth = "axlweth"
	// Native coin
	Evmos = "evmos"
)

type ERC20ModuleCoin struct {
	Denom       string
	Erc20       string
	TokenName   string
	Description string
	CoingeckoID string
	Name        string
	Symbol      string
	Decimals    int
}

var ERC20ModuleCoins = map[string]ERC20ModuleCoin{
	// Note: aevmos and wevmos are not created using the Erc20 module
	Evmos: {
		Denom:       "aevmos",
		Erc20:       "0xD4949664cD82660AaE99bEdc034a0deA8A0bd517",
		TokenName:   "Evmos",
		Description: "Evmos native coin",
		CoingeckoID: "evmos",
		Name:        "Evmos",
		Symbol:      "EVMOS",
		Decimals:    18,
	},

	Gweth: {
		Denom:       "ibc/6B3FCE336C3465D3B72F7EFB4EB92FC521BC480FE9653F627A0BD0237DF213F3",
		Erc20:       "0xc03345448969Dd8C00e9E4A85d2d9722d093aF8E",
		TokenName:   "gWETH",
		Description: "Gravity Bridge WETH",
		CoingeckoID: "weth",
		Name:        "Wrapped Ether channel-8",
		Symbol:      "ibc G-WETH",
		Decimals:    18,
	},

	Gdai: {
		Denom:       "ibc/F96A7F81E8F82E4EE81F94D507CD257319EFB70FE46E23B4953F63B62E855603",
		Erc20:       "0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75",
		TokenName:   "gDAI",
		Description: "Gravity Bridge DAI",
		CoingeckoID: "dai",
		Name:        "Dai Stablecoin channel-8",
		Symbol:      "ibc G-DAI",
		Decimals:    18,
	},

	Gwbtc: {
		Denom:       "ibc/350B6DC0FF48E3BDB856F40A8259909E484259ED452B3F4F39A0FEF874F30F61",
		Erc20:       "0x1D54EcB8583Ca25895c512A8308389fFD581F9c9",
		TokenName:   "gWBTC",
		Description: "Gravity Bridge BTC",
		CoingeckoID: "wrapped-bitcoin",
		Name:        "Wrapped Bitcoin channel-8",
		Symbol:      "ibc G-WBTC",
		Decimals:    8,
	},

	Grav: {
		Denom:       "ibc/7F0C2CB6E79CC36D29DA7592899F98E3BEFD2CF77A94340C317032A78812393D",
		Erc20:       "0x80b5a32E4F032B2a058b4F29EC95EEfEEB87aDcd",
		TokenName:   "GRAV",
		Description: "Gravity Bridge native coin",
		CoingeckoID: "graviton",
		Name:        "Graviton channel-8",
		Symbol:      "ibc GRAV",
		Decimals:    6,
	},

	Gusdt: {
		Denom:       "ibc/DF63978F803A2E27CA5CC9B7631654CCF0BBC788B3B7F0A10200508E37C70992",
		Erc20:       "0xecEEEfCEE421D8062EF8d6b4D814efe4dc898265",
		TokenName:   "gUSDT",
		Description: "Gravity Bridge USDT",
		CoingeckoID: "tether",
		Name:        "Tether USD channel-8",
		Symbol:      "ibc G-USDT",
		Decimals:    6,
	},

	Usdcgrv: {
		Denom:       "ibc/693989F95CF3279ADC113A6EF21B02A51EC054C95A9083F2E290126668149433",
		Erc20:       "0x5FD55A1B9FC24967C4dB09C513C3BA0DFa7FF687",
		TokenName:   "USDC.grv",
		Description: "Gravity Bridge USDC",
		CoingeckoID: "usd-coin",
		Name:        "USD Coin channel-8",
		Symbol:      "ibc G-USDC",
		Decimals:    6,
	},

	Osmo: {
		Denom:       "ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518",
		Erc20:       "0xFA3C22C069B9556A4B2f7EcE1Ee3B467909f4864",
		TokenName:   "OSMO",
		Description: "The native token of Osmosis",
		CoingeckoID: "osmosis",
		Name:        "Osmosis",
		Symbol:      "OSMO",
		Decimals:    6,
	},

	Atom: {
		Denom:       "ibc/A4DB47A9D3CF9A068D454513891B526702455D3EF08FB9EB558C561F9DC2B701",
		Erc20:       "0xC5e00D3b04563950941f7137B5AfA3a534F0D6d6",
		TokenName:   "ATOM",
		Description: "The native token of Cosmos Hub",
		CoingeckoID: "cosmos",
		Name:        "Cosmos Hub",
		Symbol:      "ATOM",
		Decimals:    6,
	},
	Juno: {
		Denom:       "ibc/448C1061CE97D86CC5E86374CD914870FB8EBA16C58661B5F1D3F46729A2422D",
		Erc20:       "0x3452e23F9c4cC62c70B7ADAd699B264AF3549C19",
		TokenName:   "JUNO",
		Description: "The native token of Juno",
		CoingeckoID: "juno-network",
		Name:        "Juno",
		Symbol:      "JUNO",
		Decimals:    6,
	},
	Axlusdc: {
		Denom:       "ibc/63C53CBDF471D4E867366ABE2E631197257118D1B2BEAD1946C8A408F96464C3",
		Erc20:       "0x15C3Eb3B621d1Bff62CbA1c9536B7c1AE9149b57",
		TokenName:   "axlUSDC",
		Description: "Circle's stablecoin on Axelar",
		CoingeckoID: "usd-coin",
		Name:        "USD Coin by Axelar",
		Symbol:      "axlUSDC",
		Decimals:    6,
	},
	Axlwbtc: {
		Denom:       "ibc/C834CD421B4FD910BBC97E06E86B5E6F64EA2FE36D6AE0E4304C2E1FB1E7333C",
		Erc20:       "0xF5b24c0093b65408ACE53df7ce86a02448d53b25",
		TokenName:   "axlWBTC",
		Description: "Wrapped Bitcoin on Axelar",
		CoingeckoID: "wrapped-bitcoin",
		Name:        "Wrapped Bitcoin on Axelar",
		Symbol:      "axlWBTC",
		Decimals:    8,
	},
	Axlweth: {
		Denom:       "ibc/356EDE917394B2AEF7F915EB24FA683A0CCB8D16DD4ECCEDC2AD0CEC6B66AC81",
		Erc20:       "0x50dE24B3f0B3136C50FA8A3B8ebc8BD80a269ce5",
		TokenName:   "axlWETH",
		Description: "Wrapped Ether on Axelar",
		CoingeckoID: "weth",
		Name:        "Wrapped Ether on Axelar",
		Symbol:      "axlWETH",
		Decimals:    18,
	},
}
