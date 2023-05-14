package handler

import (
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
)

type StatusResponse struct {
	Status string `json:"status"`
}

// Status handles GET /status.
// Dummy endpoint to check if the server is up and running.
func (h *Handler) Status(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
	resp := StatusResponse{
		Status: "OK",
	}
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		ctx.Logger().Printf("Error encoding response: %s", err.Error())
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Write(jsonResponse)
}
