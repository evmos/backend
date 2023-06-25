package v2

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/valyala/fasthttp"
)

type BroadcastParams struct {
	// which network should the transaction be broadcasted to
	Network string `json:"network"`
	// the signed transaction to be broadcasted
	TxBytes []byte `json:"tx_bytes"`
}

type BroadcastTxResponse struct {
	tx.BroadcastTxResponse
}

// BroadcastTx handles POST /tx/broadcast.
// It broadcasts a signed transaction to the specified network.
// Returns:
//
//	{
//	  "tx_response": {
//	    "height": "0",
//	    "txhash": "3CB7FCC9F5FB31E530CC15665F3FD655AE6CB56CDACAD58D1395C68EDD50D0BB",
//	    "codespace": "",
//	    "code": 0,
//	    "data": "",
//	    "raw_log": "[]",
//	    "logs": [],
//	    "info": "",
//	    "gas_wanted": "0",
//	    "gas_used": "0",
//	    "tx": null,
//	    "timestamp": "",
//	    "events": []
//	  }
//	}
func (h *Handler) BroadcastTx(ctx *fasthttp.RequestCtx) {
	reqParams := BroadcastParams{}
	if err := json.Unmarshal(ctx.PostBody(), &reqParams); err != nil {
		ctx.Logger().Printf("Error decoding request body: %s", err.Error())
		sendBadRequestResponse(ctx, "Invalid request body")
		return
	}

	txRequest := tx.BroadcastTxRequest{
		TxBytes: reqParams.TxBytes,
		Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
	}

	jsonTxRequest, err := json.Marshal(txRequest)
	if err != nil {
		ctx.Logger().Printf("Error marshaling txRequest: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	restClient, err := NewRestClient(reqParams.Network)
	if err != nil {
		ctx.Logger().Printf("Error creating rest client: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	txResponse, err := restClient.BroadcastTx(jsonTxRequest)
	if err != nil {
		ctx.Logger().Printf("Error broadcasting tx: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	sendSuccessfulJSONResponse(ctx, txResponse)
}
