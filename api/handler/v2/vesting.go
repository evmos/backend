// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v2

import (
	"encoding/json"

	"github.com/evmos/evmos/v12/x/vesting/types"
	v1 "github.com/tharsis/dashboard-backend/api/handler/v1"
	"github.com/tharsis/dashboard-backend/internal/v2/node/rest"
	"github.com/valyala/fasthttp"
)

type VestingAccount struct {
	Account struct {
		Type               string                `json:"@type"`
		BaseVestingAccount v1.BaseVestingAccount `json:"base_vesting_account"`
		FunderAddress      string                `json:"funder_address"`
		StartTime          string                `json:"start_time"`
		LockupPeriods      []struct {
			Length string              `json:"length"`
			Amount []v1.BalanceElement `json:"amount"`
		} `json:"lockup_periods"`
		VestingPeriods []struct {
			Length string              `json:"length"`
			Amount []v1.BalanceElement `json:"amount"`
		} `json:"vesting_periods"`
	} `json:"account"`
}

type VestingByAddressResponse struct {
	VestingAccount
	types.QueryBalancesResponse
}

// VestingByAddress handles GET /v2/vesting/{address}.
// It returns the vesting information of the requested address.
// It handles both Hex and Bech32 addresses.
// Returns:
//
//	{
//		"locked": [
//		  {
//			"denom": "aevmos",
//			"amount": "10000000000000000"
//		  }
//		],
//		"unvested": [
//		  {
//			"denom": "aevmos",
//			"amount": "10000000000000000"
//		  }
//		],
//		"vested": [],
//		"account": {
//		  "@type": "/evmos.vesting.v1.ClawbackVestingAccount",
//		  "base_vesting_account": {
//			"base_account": {
//			  "address": "evmos1fwrmzh6kp2dh0wuevhzfsck0eeeqc54tpvkvc2",
//			  "pub_key": "",
//			  "account_number": "72669907",
//			  "sequence": "0"
//			},
//			"original_vesting": [
//			  {
//				"denom": "aevmos",
//				"amount": "10000000000000000"
//			  }
//			],
//			"end_time": "1812268746"
//		  },
//		  "funder_address": "evmos1khuaarl64qmnh96yhhlhm80qy3rgv0kcxqacha",
//		  "start_time": "2023-06-06T07:59:06Z",
//		  "lockup_periods": [
//			{
//			  "length": "31622400",
//			  "amount": [
//				{
//				  "denom": "aevmos",
//				  "amount": "10000000000000000"
//				}
//			  ]
//			}
//		  ],
//		  "vesting_periods": [
//			{
//			  "length": "31622400",
//			  "amount": [
//				{
//				  "denom": "aevmos",
//				  "amount": "2500000000000000"
//				}
//			  ]
//			}
//		  ]
//		}
//	}
func (h *Handler) VestingByAddress(ctx *fasthttp.RequestCtx) {
	address := ctx.UserValue("address").(string)
	if address == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		errorResponse := &ErrorResponse{
			Error: "Missing address",
		}
		sendJSONResponse(ctx, errorResponse)
		return
	}

	restClient, err := rest.NewClient("evmos")
	if err != nil {
		ctx.Logger().Printf("Error creating rest client: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	accountRes, err := restClient.Get("/cosmos/auth/v1beta1/accounts/" + address)

	if err != nil {
		ctx.Logger().Printf("Error querying vesting account from RPC: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	var account VestingAccount
	err = json.Unmarshal(accountRes, &account)
	if err != nil {
		ctx.Logger().Printf("Error decoding vesting account: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	rewardsRes, err := restClient.Get("/evmos/vesting/v1/balances/" + address)

	if err != nil {
		ctx.Logger().Printf("Error querying vesting balance from RPC: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	var vestingBalance types.QueryBalancesResponse
	err = json.Unmarshal(rewardsRes, &vestingBalance)
	if err != nil {
		ctx.Logger().Printf("Error decoding vesting account: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	res := VestingByAddressResponse{
		VestingAccount:        account,
		QueryBalancesResponse: vestingBalance,
	}

	sendSuccessfulJSONResponse(ctx, res)
}
