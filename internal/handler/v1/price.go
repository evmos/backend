// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"encoding/json"
	"fmt"

	"github.com/fasthttp/router"
	"github.com/tharsis/dashboard-backend/internal/db"
	"github.com/tharsis/dashboard-backend/internal/resources"
	"github.com/valyala/fasthttp"
)

func EvmosPrice(ctx *fasthttp.RequestCtx) {
	price := "0"
	val, err := db.RedisGetPrice("evmos", "usd")
	if err == nil {
		price = val
	}

	sendResponse(fmt.Sprint("{\"denom\": \"evmos\", \"price\":", price, "}"), nil, ctx)
}

func CoingeckoPrice(ctx *fasthttp.RequestCtx) {
	asset := paramToString("asset", ctx)
	price := "0.0"
	val, err := db.RedisGetPrice(asset, "usd")
	if err == nil {
		price = val
	}

	sendResponse(fmt.Sprint("{\"denom\": \"", asset, "\", \"price\":", price, "}"), nil, ctx)
}

func CoingeckoPrices(ctx *fasthttp.RequestCtx) {
	tokens, err := resources.GetERC20Tokens()
	if err != nil {
		sendResponse("", err, ctx)
	}

	// We don't know the length of the tokens so we can't create an array with a fixed length
	erc20Tokens := []map[string]string{}

	for _, v := range tokens {
		price := make(map[string]string)
		price[v.CoingeckoID] = "0.0"
		val, err := db.RedisGetPrice(v.CoingeckoID, "usd")
		if err == nil {
			price[v.CoingeckoID] = val
		}
		erc20Tokens = append(erc20Tokens, price)
	}

	stringRes, err := json.Marshal(erc20Tokens)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	res := "{\"values\":" + string(stringRes) + "}"

	sendResponse(res, nil, ctx)
}

// TODO: implement assetsPrice (?)
func AddPriceRoutes(r *router.Router) {
	r.GET("/evmosPrice", EvmosPrice)
	r.GET("/coingeckoPrice/{asset}", CoingeckoPrice)
	r.GET("/coingeckoPrice", CoingeckoPrices)
}
