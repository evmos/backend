package handler

import (
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
)

type StatusResponse struct {
	Status string `json:"status"`
}

// HandleFastHTTP request handler in net/http style, i.e. method bound to MyHandler struct.
func (h *Handler) Status(ctx *fasthttp.RequestCtx) {
	// notice that we may access MyHandler properties here - see h.foobar.
	ctx.Logger().Printf("This is a test log")
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
