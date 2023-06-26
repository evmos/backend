// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package resources

type NetworkConfig struct {
	Prefix       string `json:"prefix"`
	GasPriceStep struct {
		Low     string `json:"low"`
		Average string `json:"average"`
		High    string `json:"high"`
	} `json:"gasPriceStep"`
	Bip44 struct {
		CoinType string `json:"coinType"`
	} `json:"bip44"`
	Configurations []ConfigurationEntry `json:"configurations"`
}

type ConfigurationEntry struct {
	ChainID    string   `json:"chainId"`
	ChainName  string   `json:"chainName"`
	Identifier string   `json:"identifier"`
	ClientID   string   `json:"clientId"`
	Rest       []string `json:"rest"`
	Jrpc       []string `json:"jrpc"`
	Web3       []string `json:"web3"`
	RPC        []string `json:"rpc"`
	Currencies []struct {
		CoinDenom    string `json:"coinDenom"`
		CoinMinDenom string `json:"coinMinDenom"`
		CoinDecimals string `json:"coinDecimals"`
	} `json:"currencies"`
	Source struct {
		SourceChannel         string   `json:"sourceChannel"`
		SourceIBCDenomToEvmos string   `json:"sourceIBCDenomToEvmos"`
		DestinationChannel    string   `json:"destinationChannel"`
		JSONRPC               []string `json:"jsonRPC"`
	} `json:"source"`
	ConfigurationType string `json:"configurationType"`
	ExplorerTxURL     string `json:"explorerTxUrl"`
}

type CoinConfig struct {
	CoinDenom           string `json:"coinDenom"`
	MinCoinDenom        string `json:"minCoinDenom"`
	ImgSrc              string `json:"imgSrc"`
	PngSrc              string `json:"pngSrc"`
	Type                string `json:"type"`
	Exponent            string `json:"exponent"`
	CosmosDenom         string `json:"cosmosDenom"`
	Description         string `json:"description"`
	Name                string `json:"name"`
	TokenRepresentation string `json:"tokenRepresentation"`
	Channel             string `json:"channel"`
	IsIBCEnabled        bool   `json:"isEnabled"`
	Erc20Address        string `json:"erc20Address"`
	Ibc                 struct {
		SourceDenom string `json:"sourceDenom"`
		Source      string `json:"source"`
	} `json:"ibc"`
	HiddenFromTestnet   bool `json:"hideFromTestnet"`
	HandledByExternalUI []struct {
		URL            string `json:"url"`
		HandlingAction string `json:"handlingAction"`
	} `json:"handledByExternalUI"`
	CoingeckoID      string `json:"coingeckoId"`
	Category         string `json:"category"`
	CoinSourcePrefix string `json:"coinSourcePrefix"`
}

type ERC20ModuleCoin struct {
	Denom               string
	Erc20               string
	TokenName           string
	Description         string
	CoingeckoID         string
	Name                string
	TokenRepresentation string
	Symbol              string
	Decimals            int
	ChainPrefix         string
	HandledByExternalUI []struct {
		URL            string `json:"url"`
		HandlingAction string `json:"handlingAction"`
	} `json:"handledByExternalUI"`
	PngSrc string `json:"pngSrc"`
	Prefix string `json:"prefix"`
}
