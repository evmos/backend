// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"fmt"
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

func sendResponse(val string, err error, ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetContentType("application/json")
	// if there is an error, report to sentry
	if err != nil {
		metrics.Send(err.Error())
	}
	// if there is an error and no custom text, respond with standard text
	if err != nil && val == "" { //nolint:gocritic
		fmt.Fprint(ctx, "{\"error\":\"All Endpoints are failing\"}")
	} else if err != nil && val != "" {
		// respond with custom text
		fmt.Fprint(ctx, "{\"error\":\""+val+"\"}")
	} else {
		// send successful response
		fmt.Fprint(ctx, val)
	}
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
