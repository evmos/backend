// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v2

import "github.com/valyala/fasthttp"

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
	// TODO - validate address before querying numia
	if address == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		errorResponse := &ErrorResponse{
			Error: "Missing address",
		}
		sendJSONResponse(ctx, errorResponse)
		return
	}

	rewards, err := h.numiaRPCClient.QueryVestingAccount(address)
	if err != nil {
		ctx.Logger().Printf("Error querying vesting account from Numia: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	sendSuccessfulJSONResponse(ctx, rewards)
}
