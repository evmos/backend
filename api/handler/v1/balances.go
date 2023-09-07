// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"github.com/tharsis/dashboard-backend/internal/v1/resources"
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
	denom := token

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
