// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v2

import "github.com/valyala/fasthttp"

// RewardsByAddress handles GET /v2/rewards/{address}.
// It returns the rewards of the requested address.
// It handles both Hex and Bech32 addresses.
// Returns:
// [
//
//	{
//	  "month": "2023-03-01 00:00:00.000",
//	  "address": "evmos12aqyq9d4k7a8hzh5av2xgxp0njan48498dvj2s",
//	  "withdrawn_rewards_usd": 49896.57734220073,
//	  "withdrawn_rewards_evmos": 135899.043080658
//	}
//
// ]
func (h *Handler) RewardsByAddress(ctx *fasthttp.RequestCtx) {
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

	rewards, err := h.numiaRPCClient.QueryRewards(address)
	if err != nil {
		ctx.Logger().Printf("Error querying rewards from Numia: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	sendSuccessfulJSONResponse(ctx, rewards)
}
