// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package endpoints

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/fasthttp/router"
	"github.com/tharsis/dashboard-backend/internal/blockchain"
	"github.com/tharsis/dashboard-backend/internal/requester"
	"github.com/valyala/fasthttp"
)

// Endpoints

func ProposalByID(ctx *fasthttp.RequestCtx) {
	url := blockchain.GetGovURL(ctx.QueryArgs().Peek("v1"))
	endpoint := buildThreeParamEndpoint(url+"/proposals/", paramToString("proposal_id", ctx), "?pagination.limit=200")
	val, err := getRequestRest(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

func VoteRecord(ctx *fasthttp.RequestCtx) {
	url := blockchain.GetGovURL(ctx.QueryArgs().Peek("v1"))
	endpoint := BuildFourParamEndpoint(url+"/proposals/", paramToString("proposal_id", ctx), "/votes/", paramToString("address", ctx))
	val, err := getRequestRest(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

func ProposalTally(ctx *fasthttp.RequestCtx) {
	url := blockchain.GetGovURL(ctx.QueryArgs().Peek("v1"))
	endpoint := buildThreeParamEndpoint(url+"/proposals/", paramToString("proposal_id", ctx), "/tally")
	val, err := getRequestRest(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

func Proposal(ctx *fasthttp.RequestCtx) {
	url := blockchain.GetGovURL(ctx.QueryArgs().Peek("v1"))
	endpoint := BuildTwoParamEndpoint(url+"/proposals/", paramToString("proposal_id", ctx))
	val, err := getRequestRest(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

func InflationRate(ctx *fasthttp.RequestCtx) {
	if err := enforceEvmos(ctx); err == nil {
		endpoint := "/evmos/inflation/v1/inflation_rate"
		val, err := getRequestRest(getChain(ctx), endpoint)
		sendResponse(val, err, ctx)
	}
}

func TotalUnclaimed(ctx *fasthttp.RequestCtx) {
	if err := enforceEvmos(ctx); err == nil {
		endpoint := "/evmos/claims/v1/total_unclaimed"
		val, err := getRequestRest(getChain(ctx), endpoint)
		sendResponse(val, err, ctx)
	}
}

func ClaimsParams(ctx *fasthttp.RequestCtx) {
	if err := enforceEvmos(ctx); err == nil {
		endpoint := "/evmos/claims/v1/params"
		val, err := getRequestRest(getChain(ctx), endpoint)
		sendResponse(val, err, ctx)
	}
}

func Epochs(ctx *fasthttp.RequestCtx) {
	if err := enforceEvmos(ctx); err == nil {
		endpoint := "/evmos/epochs/v1/epochs"
		val, err := getRequestRest(getChain(ctx), endpoint)
		sendResponse(val, err, ctx)
	}
}

func InflationSupply(ctx *fasthttp.RequestCtx) {
	if err := enforceEvmos(ctx); err == nil {
		endpoint := "/evmos/inflation/v1/circulating_supply"
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

func FeeMarketParams(ctx *fasthttp.RequestCtx) {
	if err := enforceEvmos(ctx); err == nil {
		val, err := FeeMarketParamsInternal(getChain(ctx))
		sendResponse(val, err, ctx)
	}
}

func ClaimsRecordsByAddress(ctx *fasthttp.RequestCtx) {
	if err := enforceEvmos(ctx); err == nil {
		endpoint := BuildTwoParamEndpoint("/evmos/claims/v1/claims_records/", paramToString("address", ctx))
		val, err := getRequestRest(getChain(ctx), endpoint)
		sendResponse(val, err, ctx)
	}
}

func StakingParams(ctx *fasthttp.RequestCtx) {
	endpoint := "/cosmos/staking/v1beta1/params"
	val, err := getRequestRest(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

func Tallying(ctx *fasthttp.RequestCtx) {
	url := blockchain.GetGovURL(ctx.QueryArgs().Peek("v1"))
	endpoint := url + "/params/tallying"
	val, err := getRequestRest(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

func IBCClientStates(ctx *fasthttp.RequestCtx) {
	endpoint := "/ibc/core/client/v1/client_states"
	val, err := getRequestRest(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

func StakingRewards(ctx *fasthttp.RequestCtx) {
	endpoint := buildThreeParamEndpoint("/cosmos/distribution/v1beta1/delegators/", paramToString("address", ctx), "/rewards")
	val, err := getRequestRest(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

func BalanceByDenom(ctx *fasthttp.RequestCtx) {
	endpoint := BuildFourParamEndpoint("/cosmos/bank/v1beta1/balances/", paramToString("address", ctx), "/by_denom?denom=", paramToString("denom", ctx))
	val, err := getRequestRest(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

func DelegationsByAddress(ctx *fasthttp.RequestCtx) {
	endpoint := buildThreeParamEndpoint("/cosmos/staking/v1beta1/delegations/", paramToString("address", ctx), "?pagination.limit=200")
	val, err := getRequestRest(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

func ValidatorsByAddress(ctx *fasthttp.RequestCtx) {
	endpoint := buildThreeParamEndpoint("/cosmos/staking/v1beta1/delegators/", paramToString("address", ctx), "/validators")
	val, err := getRequestRest(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

func UnbondingByAddress(ctx *fasthttp.RequestCtx) {
	endpoint := buildThreeParamEndpoint("/cosmos/staking/v1beta1/delegators/", paramToString("address", ctx), "/unbonding_delegations")
	val, err := getRequestRest(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

func DelegatorInfoByValidator(ctx *fasthttp.RequestCtx) {
	endpoint := BuildFourParamEndpoint("/cosmos/staking/v1beta1/validators/", paramToString("validator_address", ctx), "/delegations/", paramToString("delegator_address", ctx))
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

func Validators(ctx *fasthttp.RequestCtx) {
	val, err := GetValidators(paramToString("validator_status", ctx), getChain(ctx))
	sendResponse(val, err, ctx)
}

func Proposals(ctx *fasthttp.RequestCtx) {
	url := blockchain.GetGovURL(ctx.QueryArgs().Peek("v1"))
	endpoint := buildThreeParamEndpoint(url+"/proposals?pagination.limit=", paramToString("pagination_limit", ctx), "&pagination.reverse=true")
	val, err := getRequestRest(getChain(ctx), endpoint)
	sendResponse(val, err, ctx)
}

func AccountInternal(address string, chain string) (string, error) {
	endpoint := BuildTwoParamEndpoint("/cosmos/auth/v1beta1/accounts/", address)
	val, err := getRequestRest(chain, endpoint)
	if err != nil {
		return "", err
	}
	return val, nil
}

func Account(ctx *fasthttp.RequestCtx) {
	val, err := AccountInternal(paramToString("address", ctx), getChain(ctx))
	sendResponse(val, err, ctx)
}

func IBCClientStatusInternal(chain string, clientID string) (string, error) {
	endpoint := BuildTwoParamEndpoint("/ibc/core/client/v1/client_status/", clientID)
	val, err := getRequestRest(chain, endpoint)
	if err != nil {
		return "", err
	}
	return val, nil
}

func IBCClientStatus(ctx *fasthttp.RequestCtx) {
	val, err := IBCClientStatusInternal(getChain(ctx), paramToString("client_id", ctx))
	sendResponse(val, err, ctx)
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

func AddProxyRoutes(r *router.Router) {
	r.GET("/ProposalById/{chain}/{proposal_id}", ProposalByID)
	r.GET("/VoteRecord/{chain}/{proposal_id}/{address}", VoteRecord)
	r.GET("/ProposalTally/{chain}/{proposal_id}", ProposalTally)
	r.GET("/Proposal/{chain}/{proposal_id}", Proposal)
	r.GET("/InflationRate/{chain}", InflationRate)
	r.GET("/FeeMarketParams/{chain}", FeeMarketParams)
	r.GET("/ClaimsRecordsByAddress/{chain}/{address}", ClaimsRecordsByAddress)
	r.GET("/TotalUnclaimed/{chain}", TotalUnclaimed)
	r.GET("/ClaimsParams/{chain}", ClaimsParams)
	r.GET("/Epochs/{chain}", Epochs)
	r.GET("/StakingParams/{chain}", StakingParams)
	r.GET("/Tallying/{chain}", Tallying)
	r.GET("/IBCClientStates/{chain}", IBCClientStates)
	r.GET("/StakingRewards/{chain}/{address}", StakingRewards)
	r.GET("/BalanceByDenom/{chain}/{address}/{denom:*}", BalanceByDenom)
	r.GET("/InflationSupply/{chain}", InflationSupply)
	r.GET("/DelegationsByAddress/{chain}/{address}", DelegationsByAddress)
	r.GET("/ValidatorsByAddress/{chain}/{address}", ValidatorsByAddress)
	r.GET("/UnbondingByAddress/{chain}/{address}", UnbondingByAddress)
	r.GET("/DelegatorInfoByValidator/{chain}/{validator_address}/{delegator_address}", DelegatorInfoByValidator)
	r.GET("/Validators/{chain}/{validator_status}", Validators)
	r.GET("/Proposals/{chain}/{pagination_limit}", Proposals)
	r.GET("/Account/{chain}/{address}", Account)
	r.GET("/IBCClientStatus/{chain}/{client_id}", IBCClientStatus)
	r.GET("/TxStatus/{chain}/{tx_hash}", TxStatus)
	r.POST("/broadcast", Broadcast)
	r.POST("/simulate", Simulate)
}
