// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tharsis/dashboard-backend/internal/db"
	"github.com/tharsis/dashboard-backend/internal/requester"
	"github.com/tharsis/dashboard-backend/internal/resources"
)

var running = true

func processAssets(erc20ModuleCoins []resources.CoinConfig) {
	for _, v := range erc20ModuleCoins {
		fmt.Println("Getting price for ", v.Name)

		res, err := requester.GetRequestPrice(v.CoingeckoID, "usd")
		if err != nil {
			// TODO: report to sentry? Set price as zero?
			panic(err)
		}

		var jsonRes map[string]map[string]float64
		_ = json.Unmarshal([]byte(res), &jsonRes)

		price := jsonRes[v.CoingeckoID]["usd"]
		stringPrice := fmt.Sprintf("%f", price)

		db.RedisSetPrice(v.CoingeckoID, "usd", stringPrice)

		fmt.Printf("Price %s for %s", stringPrice, v.Name)

		time.Sleep(5 * time.Second)
	}
}

func main() {
	for running {
		fmt.Println("Fetching ERC20 tokens...")

		erc20ModuleCoins, err := resources.GetERC20Tokens()
		if err != nil {
			// TODO: Add retries and report error to sentry??
			panic(err)
		}

		fmt.Println("Getting prices...")

		processAssets(erc20ModuleCoins)

		time.Sleep(5 * time.Second)
	}
}
