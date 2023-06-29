package v2

import (
	"encoding/json"

	"github.com/tharsis/dashboard-backend/internal/node/rest"

	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/valyala/fasthttp"
)

// BroadcastTxParams represents the parameters for the POST /v2/tx/broadcast endpoint.
type BroadcastTxParams struct {
	// which network should the transaction be broadcasted to
	Network string `json:"network"`
	// the signed transaction to be broadcasted
	TxBytes []byte `json:"tx_bytes"`
}

type BroadcastTxResponse struct {
	Code   uint32 `json:"code"`
	TxHash string `json:"tx_hash"`
	RawLog string `json:"raw_log"`
}

// BroadcastTx handles POST /tx/broadcast.
// It broadcasts a signed transaction synchronously to the specified network.
// Returns:
//
//	{
//	  "txhash": "3CB7FCC9F5FB31E530CC15665F3FD655AE6CB56CDACAD58D1395C68EDD50D0BB",
//	  "code": 0,
//	  "raw_log": "[]",
//	}
func (h *Handler) BroadcastTx(ctx *fasthttp.RequestCtx) {
	reqParams := BroadcastTxParams{}
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

	restClient, err := rest.NewClient(reqParams.Network)
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

	response := BroadcastTxResponse{
		Code:   txResponse.TxResponse.Code,
		TxHash: txResponse.TxResponse.TxHash,
		RawLog: txResponse.TxResponse.RawLog,
	}
	sendSuccessfulJSONResponse(ctx, &response)
}
