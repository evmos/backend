// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"

	decimal "github.com/cosmos/cosmos-sdk/types"
	"github.com/tharsis/dashboard-backend/internal/v1/blockchain"
	"github.com/tharsis/dashboard-backend/internal/v1/constants"
	"github.com/tharsis/dashboard-backend/internal/v1/db"
	"github.com/tharsis/dashboard-backend/internal/v1/requester"
	"github.com/tharsis/dashboard-backend/internal/v1/resources"
	"github.com/tharsis/dashboard-backend/internal/v1/utils"
	"github.com/valyala/fasthttp"
	"golang.org/x/exp/slices"
)

type Pagination struct {
	NextKey interface{} `json:"next_key"`
	Total   string      `json:"total"`
}

type BalanceElement struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type BalanceJSON struct {
	Balances   []BalanceElement `json:"balances"`
	Pagination Pagination       `json:"pagination"`
}

type ERC20Entry struct {
	Name                string `json:"name"`
	Symbol              string `json:"symbol"`
	Decimals            int    `json:"decimals"`
	Erc20Balance        string `json:"erc20Balance"`
	CosmosBalance       string `json:"cosmosBalance"`
	TokenName           string `json:"tokenName"`
	TokenIdentifier     string `json:"tokenIdentifier"`
	Description         string `json:"description"`
	CoingeckoPrice      string `json:"coingeckoPrice"`
	ChainID             string `json:"chainId"`
	ChainIdentifier     string `json:"chainIdentifier"`
	HandledByExternalUI []struct {
		URL            string `json:"url"`
		HandlingAction string `json:"handlingAction"`
	} `json:"handledByExternalUI"`
	ERC20Address   string `json:"erc20Address"`
	PngSrc         string `json:"pngSrc"`
	Prefix         string `json:"prefix"`
	Price24HChange string `json:"price24HChange"`
}

type ModuleBalanceContainer struct {
	// mu             sync.Mutex // TODO: Mutex not used
	values         map[string]ERC20Entry
	cosmosBalances []BalanceElement
}

func getTotalBalance(balance ERC20Entry) decimal.Dec {
	erc20Balance := new(big.Int)
	erc20Balance.SetString(balance.Erc20Balance, 10)

	ibcBalance := new(big.Int)
	ibcBalance.SetString(balance.CosmosBalance, 10)
	res, err := utils.NumberToBiggerDenom(erc20Balance.Add(erc20Balance, ibcBalance).String(), uint64(balance.Decimals))
	if err != nil {
		return decimal.NewDec(10)
	}
	return res
}

func buildValuesResponse(values string) string {
	return "{\"values\":" + values + "}"
}

func ERC20ModuleEmptyBalance(ctx *fasthttp.RequestCtx) {
	container := ModuleBalanceContainer{
		values:         map[string]ERC20Entry{},
		cosmosBalances: []BalanceElement{},
	}

	erc20ModuleCoins, err := resources.GetERC20ModuleCoins()
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	networkConfigs, err := resources.GetNetworkConfigs()
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	index := 0
	for k, v := range erc20ModuleCoins {
		configIdx := slices.IndexFunc(networkConfigs, func(c resources.NetworkConfig) bool { return c.Prefix == v.ChainPrefix })
		networkConfig := networkConfigs[configIdx]
		mainnetConfig := resources.GetMainnetConfig(networkConfig)
		coingeckoPrice := GetCoingeckoPrice(v.CoingeckoID)
		coin24hChnage := GetCoingecko24HChange(v.CoingeckoID)
		container.values[k] = ERC20Entry{
			Name:                v.Name,
			Symbol:              v.Symbol,
			Decimals:            v.Decimals,
			Erc20Balance:        "0",
			CosmosBalance:       "0",
			TokenName:           v.TokenName,
			TokenIdentifier:     v.TokenRepresentation,
			Description:         v.Description,
			ChainID:             mainnetConfig.ChainID,
			ChainIdentifier:     mainnetConfig.Identifier,
			HandledByExternalUI: v.HandledByExternalUI,
			CoingeckoPrice:      coingeckoPrice,
			ERC20Address:        v.Erc20,
			PngSrc:              v.PngSrc,
			Prefix:              v.Prefix,
			Price24HChange:      coin24hChnage,
		}
		index++
	}

	balance := []ERC20Entry{}
	var evmosBalance ERC20Entry

	for _, v := range container.values {
		if strings.ToLower(v.TokenName) == constants.Evmos {
			evmosBalance = v
		} else {
			balance = append(balance, v)
		}
	}

	sort.SliceStable(balance, func(i, j int) bool {
		return strings.ToLower(balance[i].Symbol) < strings.ToLower(balance[j].Symbol)
	})

	jsonresponse, err := json.Marshal(map[string][]ERC20Entry{"balance": append([]ERC20Entry{evmosBalance}, balance...)})
	if err != nil {
		sendResponse("", err, ctx)
	}
	sendResponse(string(jsonresponse), nil, ctx)
}

func ERC20ModuleBalance(ctx *fasthttp.RequestCtx) {
	evmosAddress := paramToString("evmos_address", ctx)
	ethAddress := paramToString("eth_address", ctx)
	endpoint := BuildTwoParamEndpoint("/cosmos/bank/v1beta1/balances/", evmosAddress)
	val, err := getRequestRest("EVMOS", endpoint)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	var m BalanceJSON
	if err = json.Unmarshal([]byte(val), &m); err != nil {
		sendResponse("", err, ctx)
		return
	}

	container := ModuleBalanceContainer{
		values:         map[string]ERC20Entry{},
		cosmosBalances: m.Balances,
	}

	erc20ModuleCoins, err := resources.GetERC20ModuleCoins()
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	networkConfigs, err := resources.GetNetworkConfigs()
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	index := 0
	for k, v := range erc20ModuleCoins {
		// TODO: consider moving this to a work or remove the container mutex
		val, err := blockchain.GetERC20Balance(v.Erc20, ethAddress)
		balance := "0"
		if err == nil && val != "" {
			balance = val
		}

		cosmosBalance := "0"
		for _, vb := range container.cosmosBalances {
			if vb.Denom == v.Denom {
				cosmosBalance = vb.Amount
			}
		}

		configIdx := slices.IndexFunc(networkConfigs, func(c resources.NetworkConfig) bool { return c.Prefix == v.ChainPrefix })
		networkConfig := networkConfigs[configIdx]
		mainnetConfig := resources.GetMainnetConfig(networkConfig)
		coingeckoPrice := GetCoingeckoPrice(v.CoingeckoID)
		coin24hChnage := GetCoingecko24HChange(v.CoingeckoID)

		container.values[k] = ERC20Entry{
			Name:                v.Name,
			Symbol:              v.Symbol,
			Decimals:            v.Decimals,
			Erc20Balance:        balance,
			CosmosBalance:       cosmosBalance,
			TokenName:           v.TokenName,
			TokenIdentifier:     v.TokenRepresentation,
			Description:         v.Description,
			CoingeckoPrice:      coingeckoPrice,
			ChainID:             mainnetConfig.ChainID,
			ChainIdentifier:     mainnetConfig.Identifier,
			HandledByExternalUI: v.HandledByExternalUI,
			ERC20Address:        v.Erc20,
			PngSrc:              v.PngSrc,
			Prefix:              v.Prefix,
			Price24HChange:      coin24hChnage,
		}
		index++
	}

	// This will keep the same order in case we move to threads
	balance := []ERC20Entry{}
	var evmosBalance ERC20Entry
	zeroBalance := []ERC20Entry{}

	for _, v := range container.values {
		if strings.ToLower(v.TokenName) == "evmos" {
			evmosBalance = v
		} else {
			totalBalance := getTotalBalance(v)
			if totalBalance.IsZero() {
				zeroBalance = append(zeroBalance, v)
			} else {
				balance = append(balance, v)
			}
		}
	}

	sort.SliceStable(balance, func(i, j int) bool {
		iBalance := getTotalBalance(balance[i])
		jBalance := getTotalBalance(balance[j])
		if iBalance.Equal(jBalance) {
			return false
		}
		return iBalance.Abs().GT(jBalance)
	})

	sort.SliceStable(zeroBalance, func(i, j int) bool {
		return strings.ToLower(zeroBalance[i].Symbol) < strings.ToLower(zeroBalance[j].Symbol)
	})

	totalBalance := append(balance, zeroBalance...) //nolint:all

	jsonresponse, err := json.Marshal(map[string][]ERC20Entry{"balance": append([]ERC20Entry{evmosBalance}, totalBalance...)})
	if err != nil {
		sendResponse("", err, ctx)
	}
	sendResponse(string(jsonresponse), nil, ctx)
}

type TokensByNameIBC struct {
	SourceDenom string `json:"sourceDenom"`
	Source      string `json:"source"`
	// TODO: add denom from chain to evmos (sending evmos)
}

type TokensByNameConfig struct {
	CosmosDenom  string          `json:"cosmosDenom"`
	Ibc          TokensByNameIBC `json:"ibc"`
	ERC20Address string          `json:"erc20Address"`
}

type TokensByName struct {
	Values TokensByNameConfig `json:"values"`
}

func ERC20TokensByNameInternal(name string) (string, error) {
	// the name it has to be equal to the one that is in the github repo.
	if val, err := db.RedisGetERC20TokensByName(name); err == nil {
		return val, nil
	}

	val, err := requester.GetERC20TokensDirectory()
	if err != nil {
		return "", err
	}

	for _, v := range val {
		if strings.Contains(v.URL, name) {
			res := buildValuesResponse(v.Content)
			db.RedisSetERC20TokensByName(name, res)
			return res, nil
		}
	}
	return "", fmt.Errorf("invalid token, please try again")
}
