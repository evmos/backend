// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/tharsis/dashboard-backend/internal/v1/blockchain"
	"github.com/tharsis/dashboard-backend/internal/v1/requester"
	"github.com/valyala/fasthttp"
)

// Endpoints

func VoteRecord(ctx *fasthttp.RequestCtx) {
	url := blockchain.GetGovURL(ctx.QueryArgs().Peek("v1"))
	endpoint := BuildFourParamEndpoint(url+"/proposals/", paramToString("proposal_id", ctx), "/votes/", paramToString("address", ctx))
	val, err := getRequestRest(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

func Epochs(ctx *fasthttp.RequestCtx) {
	if err := enforceEvmos(ctx); err == nil {
		endpoint := "/evmos/epochs/v1/epochs"
		val, err := getRequestRest(getChain(ctx), endpoint)
		sendResponse(val, err, ctx)
	}
}

func EthGasPriceInternal() (string, error) {
	url := "https://eth.bd.evmos.org:8545"
	val, _ := requester.MakePostGasPrice(url)
	return val, nil
}

func FeeMarketParamsInternal(chain string) (string, error) {
	if chain == "EVMOS" {
		endpoint := "/evmos/feemarket/v1/params"
		val, err := getRequestRest(chain, endpoint)
		return val, err
	}

	return "", fmt.Errorf("network is not Evmos")
}

func BalanceByDenom(ctx *fasthttp.RequestCtx) {
	endpoint := BuildFourParamEndpoint("/cosmos/bank/v1beta1/balances/", paramToString("address", ctx), "/by_denom?denom=", paramToString("denom", ctx))
	val, err := getRequestRest(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

func GetValidators(status string, chain string) (string, error) {
	endpoint := buildThreeParamEndpoint("/cosmos/staking/v1beta1/validators?status=", status, "&pagination.limit=200")
	return getRequestRest(chain, endpoint)
}

func GetAllValidators(chain string) (string, error) {
	endpoint := "/cosmos/staking/v1beta1/validators?pagination.limit=500"
	return getRequestRest(chain, endpoint)
}

func AccountInternal(address string, chain string) (string, error) {
	endpoint := BuildTwoParamEndpoint("/cosmos/auth/v1beta1/accounts/", address)
	val, err := getRequestRest(chain, endpoint)
	if err != nil {
		return "", err
	}
	return val, nil
}

func IBCClientStatusInternal(chain string, clientID string) (string, error) {
	endpoint := BuildTwoParamEndpoint("/ibc/core/client/v1/client_status/", clientID)
	val, err := getRequestRest(chain, endpoint)
	if err != nil {
		return "", err
	}
	return val, nil
}

func TxStatus(ctx *fasthttp.RequestCtx) {
	endpoint := BuildTwoParamEndpoint("/tx?hash=0x", paramToString("tx_hash", ctx))
	val, err := GetRequestJrpc(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

type broadcastParams struct {
	Network string  `json:"network"`
	TxBytes []uint8 `json:"txBytes"`
	Sender  string  `json:"sender"`
}

type simulateParams struct {
	Network string  `json:"network"`
	TxBytes []uint8 `json:"txBytes"`
}

func parseErrorString(errorString string) string {
	if strings.Contains(errorString, "insufficient fees") {
		return "Insufficient fees, please try again"
	}
	if strings.Contains(errorString, "fee") {
		return "Fee changed while sending transaction, please try again"
	}
	if strings.Contains(errorString, "insufficient funds") {
		return "Insufficient funds, please try again"
	}
	if strings.Contains(errorString, "aevmos is smaller than") {
		return "Fee changed while sending transaction, please try again"
	}
	return errorString
}

func ConvertTxBytesToString(txBytes []uint8) string {
	var txBytesTemp strings.Builder
	for _, v := range txBytes {
		txBytesTemp.WriteString(fmt.Sprint(v))
		txBytesTemp.WriteString(",")
	}

	// Remove last ,
	localTxBytes := txBytesTemp.String()
	lenBytes := len(localTxBytes)
	if lenBytes > 2 {
		localTxBytes = localTxBytes[:lenBytes-1]
	}
	return localTxBytes
}

func SimulateInternal(network string, txBytes string) (bool, string) {
	var sb strings.Builder
	sb.WriteString(`{"tx_bytes":[`)
	sb.WriteString(txBytes)
	sb.WriteString(`]}`)
	jsonBody := []byte(sb.String())

	val, err := requester.MakePostRequest(network, "rest", "/cosmos/tx/v1beta1/simulate", jsonBody)
	if err != nil {
		return false, fmt.Sprint(err)
	}

	var res map[string]interface{}
	err = json.Unmarshal([]byte(val), &res)
	if err != nil {
		return false, fmt.Sprint(err)
	}

	if _, ok := res["error"]; ok {
		return false, fmt.Sprint(res["error"])
	}
	return true, "Transaction was simulated correctly"
}

func Simulate(ctx *fasthttp.RequestCtx) {
	m := simulateParams{}
	if err := json.Unmarshal(ctx.PostBody(), &m); err != nil {
		sendResponse("", err, ctx)
		return
	}
	txBytes := ConvertTxBytesToString(m.TxBytes)
	success, err := SimulateInternal(m.Network, txBytes)
	sendResponse("{\"status\": "+strconv.FormatBool(success)+", \"message\": \""+err+"\"}", nil, ctx)
}

func broadcastInternal(bytes []byte, network string) (string, error) {
	txBytes := ConvertTxBytesToString(bytes)

	var sb strings.Builder
	sb.WriteString(`{"tx_bytes":[`)
	sb.WriteString(txBytes)
	sb.WriteString(`], "mode": "BROADCAST_MODE_SYNC"}`)
	jsonBody := []byte(sb.String())

	// simulate
	if network != "EMONEY" {
		// emoney uses a cosmos sdk version that does not match with the simulate
		// that we are using
		if success, msg := SimulateInternal(network, txBytes); !success {
			errorString := parseErrorString(msg)
			var sb strings.Builder
			sb.WriteString(`{"error":"`)
			sb.WriteString(errorString)
			sb.WriteString(`","tx_hash":null}`)
			return sb.String(), nil
		}
	}
	val, err := requester.MakeLongPostRequest(network, "rest", "/cosmos/tx/v1beta1/txs", jsonBody)
	if err != nil {
		return "", err
	}

	var res map[string]interface{}
	err = json.Unmarshal([]byte(val), &res)
	if err != nil {
		return "", err
	}

	if txResponseRaw, ok := res["tx_response"]; ok {
		if txResponse, ok := txResponseRaw.(map[string]interface{}); ok {
			if code, ok := txResponse["code"]; ok {
				if code, ok := code.(float64); ok {
					// Error sending transaction
					if code != 0 {
						if rawLogRaw, ok := txResponse["raw_log"]; ok {
							if rawLog, ok := rawLogRaw.(string); ok {
								errorString := parseErrorString(rawLog)
								var sb strings.Builder
								sb.WriteString(`{"error":"`)
								sb.WriteString(errorString)
								sb.WriteString(`","tx_hash":null}`)
								return sb.String(), nil
							}
						}
					}
					// Valid transaction
					if txHashRaw, ok := txResponse["txhash"]; ok {
						if txHash, ok := txHashRaw.(string); ok {
							var sb strings.Builder
							sb.WriteString(`{"error":null, "tx_hash":"`)
							sb.WriteString(txHash)
							sb.WriteString(`"}`)
							return sb.String(), nil
						}
					}
				}
			}
		}
	}
	return "", fmt.Errorf("invalid transaction response")
}

func Broadcast(ctx *fasthttp.RequestCtx) {
	m := broadcastParams{}
	if err := json.Unmarshal(ctx.PostBody(), &m); err != nil {
		sendResponse("", err, ctx)
		return
	}
	val, err := broadcastInternal(m.TxBytes, m.Network)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}
	sendResponse(val, err, ctx)
}
