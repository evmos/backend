// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package endpoints

import (
	"encoding/json"

	"github.com/fasthttp/router"
	"github.com/tharsis/dashboard-backend/internal/blockchain"
	"github.com/tharsis/dashboard-backend/internal/resources"
	"github.com/valyala/fasthttp"
)

type EVMOSBalance struct {
	Chain        string `json:"chain"`
	EvmosBalance string `json:"evmosBalance"`
}

type EVMOSIBCBalancesReponse struct {
	Balances []EVMOSBalance `json:"balances"`
}

type BalanceResponse struct {
	Balance BalanceElement `json:"balance"`
}

func BalanceByNetworkAndDenom(ctx *fasthttp.RequestCtx) {
	token := paramToString("token", ctx)
	denom := ""

	coinConfigs, err := resources.GetERC20Tokens()
	if err != nil {
		sendResponse("", err, ctx)
		return
	}
	for _, v := range coinConfigs {
		if v.CoinDenom == token {
			denom = v.Ibc.SourceDenom
		}
	}
	if denom == "" {
		sendResponse("", err, ctx)
		return
	}

	endpoint := BuildFourParamEndpoint("/cosmos/bank/v1beta1/balances/", paramToString("address", ctx), "/by_denom?denom=", denom)
	val, err := getRequestRest(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

func EVMOSIBCBalances(ctx *fasthttp.RequestCtx) {
	pubKey := paramToString("pubkey", ctx)

	networkConfigs, err := resources.GetNetworkConfigs()
	if err != nil {
		sendResponse("Unable to get registry configurations", err, ctx)
		return
	}

	balances := []EVMOSBalance{}

	// TODO: We should send all the requests at the same time to speed up this call
	for _, v := range networkConfigs {
		derivedAddress, err := blockchain.DeriveCosmosAddress(pubKey, v.Prefix)
		if err != nil {
			// Unable to derive address
			continue
		}

		configuration := resources.GetMainnetConfig(v)
		evmosIbcDenom, err := GetDenom("EVMOS", configuration.Identifier)
		if err != nil {
			// Unable to get EVMOS denom for source chain
			continue
		}

		endpoint := BuildFourParamEndpoint("/cosmos/bank/v1beta1/balances/", derivedAddress, "/by_denom?denom=", evmosIbcDenom)
		val, err := getRequestRest(configuration.Identifier, endpoint)
		if err != nil {
			// Unable to get EVMOS balance in chain provided
			continue
		}

		var balance BalanceResponse

		err = json.Unmarshal([]byte(val), &balance)

		if err != nil {
			// Unable to get EVMOS balance in chain provided
			continue
		}

		balances = append(balances, EVMOSBalance{
			Chain:        configuration.Identifier,
			EvmosBalance: balance.Balance.Amount,
		})
	}

	res, err := json.Marshal(balances)
	if err != nil {
		sendResponse("Unable to get EVMOS balances", err, ctx)
		return
	}

	sendResponse("{\"values\":"+string(res)+"}", err, ctx)
}

func EVMOSIBCBalance(ctx *fasthttp.RequestCtx) {
	sourceChain := getChain(ctx)

	evmosIbcDenom, err := GetDenom("EVMOS", sourceChain)
	if err != nil {
		sendResponse("Unable to get EVMOS denom in source chain provided", err, ctx)
		return
	}

	endpoint := BuildFourParamEndpoint("/cosmos/bank/v1beta1/balances/", paramToString("address", ctx), "/by_denom?denom=", evmosIbcDenom)
	val, err := getRequestRest(sourceChain, endpoint)
	if err != nil {
		sendResponse("Unable to get EVMOS balance in chain provided", err, ctx)
	}
	sendResponse(val, err, ctx)
}

func AddBalancesRoutes(r *router.Router) {
	r.GET("/BalanceByNetworkAndDenom/{chain}/{token}/{address}", BalanceByNetworkAndDenom)
	r.GET("/EVMOSIBCBalances/{pubkey:*}", EVMOSIBCBalances)
	r.GET("/EVMOSIBCBalance/{chain}/{address}", EVMOSIBCBalance)
}
