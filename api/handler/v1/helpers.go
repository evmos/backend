// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/tharsis/dashboard-backend/internal/v1/constants"
	"github.com/tharsis/dashboard-backend/internal/v1/db"
	"github.com/tharsis/dashboard-backend/internal/v1/metrics"
	"github.com/tharsis/dashboard-backend/internal/v1/requester"
	"github.com/valyala/fasthttp"
)

func getRequestRest(chain string, endpoint string) (string, error) {
	return getRequest(chain, "rest", endpoint)
}

func GetRequestJrpc(chain string, endpoint string) (string, error) {
	return getRequest(chain, "jrpc", endpoint)
}

func getRequest(chain string, endpointType string, endpoint string) (string, error) {
	val, err := db.RedisGetProxyResponse(chain, endpoint)
	if err != nil {
		val, err = requester.MakeGetRequest(chain, endpointType, endpoint)
		if err != nil {
			if val, err := db.RedisGetFallbackResponse(chain, endpoint); err == nil {
				return val, nil
			}
			return "", err
		}
		db.RedisSetProxyResponse(chain, endpoint, val)
		db.RedisSetFallbacResponse(chain, endpoint, val)
		return val, nil
	}
	return val, nil
}

func BuildTwoParamEndpoint(a, b string) string {
	var sb strings.Builder
	sb.WriteString(a)
	sb.WriteString(b)
	return sb.String()
}

func buildThreeParamEndpoint(a, b, c string) string {
	var sb strings.Builder
	sb.WriteString(a)
	sb.WriteString(b)
	sb.WriteString(c)
	return sb.String()
}

func BuildFourParamEndpoint(a, b, c, d string) string {
	var sb strings.Builder
	sb.WriteString(a)
	sb.WriteString(b)
	sb.WriteString(c)
	sb.WriteString(d)
	return sb.String()
}

func GetCoingeckoPrice(coingeckoID string) string {
	price := "0"
	val, err := db.RedisGetPrice(coingeckoID, "usd")
	if err == nil {
		price = val
	}
	return price
}

func paramToString(param string, ctx *fasthttp.RequestCtx) string {
	return fmt.Sprint(ctx.UserValue(param))
}

func getChain(ctx *fasthttp.RequestCtx) string {
	return paramToString("chain", ctx)
}

func enforceEvmos(ctx *fasthttp.RequestCtx) error {
	chain := paramToString("chain", ctx)
	if chain == constants.EVMOS {
		return nil
	}
	fmt.Fprintf(ctx, "{\"error\": \"This endpoint is EVMOS only\"}")
	return fmt.Errorf("network is not Evmos")
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func sendResponse(val string, err error, ctx *fasthttp.RequestCtx) {
	if err != nil {
		metrics.Send(err.Error())
	var msg string
		if val == "" {
			msg = "All Endpoints are failing"
		} else {
			msg = val
		}
		errResponse := ErrorResponse{Error: msg}
		sendJSONResponse(ctx, errResponse)
		return
	} else {
		sendJSONResponse(ctx, val)
		return
	}
}

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

func buildErrorResponse(a string) string {
	var sb strings.Builder
	sb.WriteString("{\"error\": \"")
	sb.WriteString(a)
	sb.WriteString("\"}")
	return sb.String()
}

func buildErrorBroadcast(a string) string {
	var sb strings.Builder
	sb.WriteString(`{"error":"`)
	sb.WriteString(a)
	sb.WriteString(`","tx_hash":null}`)
	return sb.String()
}
