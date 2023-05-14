package v2

import (
	"github.com/valyala/fasthttp"
)

type HeightResponse struct {
	Height string `json:"height"`
}

// Height handles GET "/v2/height"
// Uses Numia API to get the current block height
// Returns
//
//	{
//	 "height": "13281459"
//	}
func (h *Handler) Height(ctx *fasthttp.RequestCtx) {
	data, err := h.numiaRPCClient.QueryHeight()
	if err != nil {
		ctx.Logger().Printf("Error querying height: %s", err.Error())
		sendInternalErrorResponse(ctx, "")
		return
	}

	response := &HeightResponse{
		Height: data.LatestBlockHeight,
	}
	sendSuccessfulJSONResponse(ctx, response)
}
