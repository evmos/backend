// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v2

import (
	"encoding/json"
	"net/http"

	"github.com/tharsis/dashboard-backend/internal/v2/encoding"

	"github.com/gogo/protobuf/proto"
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
}

// sendSuccesfulProtoJSONResponse sends a successful JSON response to the client.
// It encodes
// It sets the status code to 200.
func sendSuccesfulProtoJSONResponse(ctx *fasthttp.RequestCtx, response proto.Message) {
	encConfig := encoding.MakeEncodingConfig()
	jsonResponse, err := encConfig.Codec.MarshalJSON(response)
	if err != nil {
		ctx.Logger().Printf("Error encoding response: %s", err.Error())
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.Response.Header.SetContentType("application/json")
	ctx.SetBody(jsonResponse)
}

// --- Error Responses ---
// Whenever the response is an error, any HTTP code that is not 200,
// we send a JSON response with the following format:
// {
//   "error": "error message"
// }

type ErrorResponse struct {
	Error string `json:"error"`
}

// sendInternalErrorResponse sends an internal error response to the client.
// It sets the status code to 500.
func sendInternalErrorResponse(ctx *fasthttp.RequestCtx) {
	message := "Something went wrong, please try again later"
	ctx.SetStatusCode(http.StatusInternalServerError)
	sendErrorJSONResponse(ctx, message)
}

// sendBadRequestResponse sends a bad request response to the client.
// It sets the status code to 400.
func sendBadRequestResponse(ctx *fasthttp.RequestCtx, message string) {
	ctx.SetStatusCode(http.StatusBadRequest)
	sendErrorJSONResponse(ctx, message)
}

func sendErrorJSONResponse(ctx *fasthttp.RequestCtx, message string) {
	errorResponse := &ErrorResponse{
		Error: message,
	}
	sendJSONResponse(ctx, errorResponse)
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
	ctx.SetBody(jsonResponse)
}
