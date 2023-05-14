package v2

import (
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
)

// This file contains private functions for sending responses to the client.
// This ensures we have a consistent way of sending responses.
// This also allows us to easily change the response format in the future.

// --- Success Responses ---

// sendSuccessfulJSONResponse sends a successful JSON response to the client.
// It sets the status code to 200.
func sendSuccessfulJSONResponse(ctx *fasthttp.RequestCtx, response interface{}) {
	ctx.SetStatusCode(http.StatusOK)
	sendJSONResponse(ctx, response)
	return
}

// --- Error Responses ---

type ErrorResponse struct {
	Error string `json:"error"`
}

// sendInternalErrorResponse sends an internal error response to the client.
// It sets the status code to 500.
func sendInternalErrorResponse(ctx *fasthttp.RequestCtx, message string) {
	// if message is empty, set default message for internal errors
	if message != "" {
		message = "Something went wrong, please try again later"
	}
	ctx.SetStatusCode(http.StatusInternalServerError)
	sendErrorJSONResponse(ctx, message)
	return
}

func sendErrorJSONResponse(ctx *fasthttp.RequestCtx, message string) {
	errorResponse := &ErrorResponse{
		Error: message,
	}
	sendJSONResponse(ctx, errorResponse)
	return
}

// --- JSON Responses ---

func sendJSONResponse(ctx *fasthttp.RequestCtx, response interface{}) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		ctx.Logger().Printf("Error encoding response: %s", err.Error())
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Write(jsonResponse)
	return
}
