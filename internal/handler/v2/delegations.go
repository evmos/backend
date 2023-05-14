package v2

import (
	"github.com/valyala/fasthttp"
)

// DelegationsByAddress handles GET "/v2/delegations/{address}".
// Queries Numia and returns all delegations for a given address.
// Accepts both bech32 and hex addresses.
// Returns
// [
//
//	{
//	  "validatorAddress": "evmosvaloper1....",
//	  "delegated": {
//	    "denom": "aevmos",
//	    "amount": "5000....."
//	  },
//	  "unclaimed": {
//	    "denom": "aevmos",
//	    "amount": "5000....."
//	  }
//	}
//
// ]
func (h *Handler) DelegationsByAddress(ctx *fasthttp.RequestCtx) {
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

	delegations, err := h.numiaRPCClient.QueryDelegations(address)
	if err != nil {
		ctx.Logger().Printf("Error querying delegations from Numia: %s", err.Error())
		sendInternalErrorResponse(ctx, "")
		return
	}

	sendSuccessfulJSONResponse(ctx, delegations)
	return
}
