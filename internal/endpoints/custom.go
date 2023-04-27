// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package endpoints

import (
	"encoding/json"
	"math"
	"sort"

	sdkmath "cosmossdk.io/math"

	"github.com/fasthttp/router"
	"github.com/tharsis/dashboard-backend/internal/db"
	"github.com/valyala/fasthttp"
)

const EpochsPerPeriod = 365

type Entry struct {
	CreationHeight string `json:"creation_height"`
	CompletionTime string `json:"completion_time"`
	InitialBalance string `json:"initial_balance"`
	Balance        string `json:"balance"`
}

type UnbondingResponse struct {
	DelegatorAddress string  `json:"delegator_address"`
	ValidatorAddress string  `json:"validator_address"`
	Entries          []Entry `json:"entries"`
	// Custom value used to create the response
	Validator Validator `json:"validator"`
}

type UnbondingResponseAPI struct {
	UnbondingResponses []UnbondingResponse `json:"unbonding_responses"`
	Pagination         Pagination          `json:"pagination"`
}

type Delegation struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
	// Custom value
	ValidatorRank int       `json:"rank"`
	Validator     Validator `json:"validator"`
}

type DelegationResponse struct {
	Delegation Delegation     `json:"delegation"`
	Balance    BalanceElement `json:"balance"`
}

type DelegationResponsesResponse struct {
	DelegationResponse []DelegationResponse `json:"delegation_responses"`
	Pagination         Pagination           `json:"pagination"`
}

type SkippedEpochsResponse struct {
	SkippedEpochs int `json:"skipped_epochs,string"`
}

type CurrentEpochResponse struct {
	CurrentEpoch int `json:"current_epoch,string"`
}

type RemainingEpochsResponse struct {
	RemainingEpochs int `json:"remainingEpochs"`
}

func GetValidatorsWithRanks(chain string) (map[string]Validator, error) {
	if val, err := db.RedisGetValidatorWithRanks(chain); err == nil {
		var res map[string]Validator
		err := json.Unmarshal([]byte(val), &res)
		if err == nil {
			return res, nil
		}
	}

	// We need to make a request with just the bonded validators to get the ranks
	bondedRaw, err := GetAllValidators(chain)
	if err != nil {
		return nil, err
	}
	var bonded ValidatorAPIResponse
	err = json.Unmarshal([]byte(bondedRaw), &bonded)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(bonded.Validators, func(a int, b int) bool {
		valA, okA := sdkmath.NewIntFromString(bonded.Validators[a].Tokens)
		valB, okB := sdkmath.NewIntFromString(bonded.Validators[b].Tokens)
		if !okA || !okB {
			return false
		}
		return valA.GT(valB)
	})

	valMap := make(map[string]Validator)

	// Set the ranks
	for k, v := range bonded.Validators {
		item := v
		item.Rank = k + 1
		valMap[v.OperatorAddress] = item
	}

	if val, err := json.Marshal(valMap); err == nil {
		db.RedisSetValidatorWithRanks(chain, string(val))
	}

	return valMap, nil
}

func GetValidatorsWithNoFilter(chain string) (map[string]Validator, error) {
	if val, err := db.RedisGetValidatorWithNoFilter(chain); err == nil {
		var res map[string]Validator
		err := json.Unmarshal([]byte(val), &res)
		if err == nil {
			return res, nil
		}
	}

	endpoint := "/cosmos/staking/v1beta1/validators?pagination.limit=600"

	validators, _ := getRequestRest(chain, endpoint)

	var m ValidatorAPIResponse
	err := json.Unmarshal([]byte(validators), &m)
	if err != nil {
		return nil, err
	}

	valWithRanks, err := GetValidatorsWithRanks(chain)
	if err != nil {
		return nil, err
	}

	valMap := make(map[string]Validator)
	for _, v := range m.Validators {
		if val, ok := valWithRanks[v.OperatorAddress]; ok {
			v.Rank = val.Rank
		} else {
			v.Rank = -1
		}
		valMap[v.OperatorAddress] = v
	}

	val, err := json.Marshal(valMap)
	if err != nil {
		return nil, err
	}

	db.RedisSetValidatorWithNoFilter(chain, string(val))
	return valMap, nil
}

func UnbondingByAddressWithValidatorInfo(ctx *fasthttp.RequestCtx) {
	chain := getChain(ctx)
	address := paramToString("address", ctx)

	if val, err := db.RedisGetUnbondingByAddressWithValidatorInfo(chain, address); err == nil {
		sendResponse(val, nil, ctx)
		return
	}

	endpoint := buildThreeParamEndpoint("/cosmos/staking/v1beta1/delegators/", address, "/unbonding_delegations")
	val, err := getRequestRest(chain, endpoint)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	var unbodings UnbondingResponseAPI
	err = json.Unmarshal([]byte(val), &unbodings)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	valMap, err := GetValidatorsWithNoFilter(chain)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	res := []interface{}{}
	for _, u := range unbodings.UnbondingResponses {
		_, exists := valMap[u.ValidatorAddress]
		if exists {
			u.Validator = valMap[u.ValidatorAddress]
		}
		res = append(res, u)
	}

	valuesToSend, err := json.Marshal(res)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	value := "{\"values\":" + string(valuesToSend) + "}"

	db.RedisSetUnbondingByAddressWithValidatorInfo(chain, address, value)

	sendResponse(value, nil, ctx)
}

func DelegationsByAddressWithValidatorRanks(ctx *fasthttp.RequestCtx) {
	chain := getChain(ctx)
	address := paramToString("address", ctx)

	if val, err := db.RedisGetDelegationsByAddressWithValidatorRanks(chain, address); err == nil {
		sendResponse(val, nil, ctx)
		return
	}

	endpoint := buildThreeParamEndpoint("/cosmos/staking/v1beta1/delegations/", address, "?pagination.limit=200")
	val, _ := getRequestRest(getChain(ctx), endpoint)

	// ValidatorsWithRank
	valWithRanks, err := GetValidatorsWithRanks(chain)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	// Parse the response to append the ranks
	var delegation DelegationResponsesResponse
	err = json.Unmarshal([]byte(val), &delegation)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	for i, v := range delegation.DelegationResponse {
		item := v.Delegation
		if val, ok := valWithRanks[item.ValidatorAddress]; ok {
			item.ValidatorRank = val.Rank
		} else {
			item.ValidatorRank = -1
		}
		delegation.DelegationResponse[i].Delegation = item
	}

	valuesToSend, err := json.Marshal(delegation)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	db.RedisSetDelegationsByAddressWithValidatorRanks(chain, address, string(valuesToSend))
	sendResponse(string(valuesToSend), err, ctx)
}

func ValidatorsByAddressWithValidatorRanks(ctx *fasthttp.RequestCtx) {
	chain := getChain(ctx)
	address := paramToString("address", ctx)

	if val, err := db.RedisGetValidatorsByAddressWithValidatorRanks(chain, address); err == nil {
		sendResponse(val, nil, ctx)
		return
	}

	endpoint := buildThreeParamEndpoint("/cosmos/staking/v1beta1/delegators/", address, "/validators")
	val, _ := getRequestRest(chain, endpoint)

	// ValidatorsWithRank
	valWithRanks, err := GetValidatorsWithRanks(chain)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	// Parse the response to append the ranks
	var validators ValidatorAPIResponse
	err = json.Unmarshal([]byte(val), &validators)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	for i, v := range validators.Validators {
		item := v
		if val, ok := valWithRanks[item.OperatorAddress]; ok {
			item.Rank = val.Rank
		} else {
			item.Rank = -1
		}
		validators.Validators[i] = item
	}

	valuesToSend, err := json.Marshal(validators)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	db.RedisSetValidatorsByAddressWithValidatorRanks(chain, address, string(valuesToSend))
	sendResponse(string(valuesToSend), err, ctx)
}

func RemainingEpochs(ctx *fasthttp.RequestCtx) {
	// query skipped epochs
	skippedEndpoint := "/evmos/inflation/v1/skipped_epochs"
	val, err := getRequestRest("EVMOS", skippedEndpoint)
	if err != nil {
		sendResponse("Failed to get remaining epochs from endpoint", err, ctx)
		return
	}

	var skippedEpochsRes SkippedEpochsResponse
	err = json.Unmarshal([]byte(val), &skippedEpochsRes)

	if err != nil {
		sendResponse("Failed to get remaining epochs from endpoint", err, ctx)
		return
	}

	// query current epochs
	currentEndpoint := "/evmos/epochs/v1/current_epoch?identifier=day"
	val, err = getRequestRest("EVMOS", currentEndpoint)
	if err != nil {
		sendResponse("Failed to get remaining epochs from endpoint", err, ctx)
		return
	}

	var currentEpochsRes CurrentEpochResponse
	err = json.Unmarshal([]byte(val), &currentEpochsRes)

	if err != nil {
		sendResponse("Failed to get remaining epochs from endpoint", err, ctx)
		return
	}

	skippedEpochs := skippedEpochsRes.SkippedEpochs
	currentEpochs := currentEpochsRes.CurrentEpoch

	// calculate epochs passes
	epochsPassed := currentEpochs - skippedEpochs

	// calculate how many periods have passed
	periodsPassed := func() int {
		return int(math.Floor(float64(epochsPassed) / float64(EpochsPerPeriod)))
	}()

	// calculate remaining epochs
	remainingEpochs := RemainingEpochsResponse{RemainingEpochs: EpochsPerPeriod*(1+periodsPassed) - epochsPassed}
	res, err := json.Marshal(remainingEpochs)
	if err != nil {
		sendResponse("Failed to get remaining epochs from endpoint", err, ctx)
		return
	}

	sendResponse(string(res), err, ctx)
}

func AddCustomRoutes(r *router.Router) {
	r.GET("/UnbondingByAddressWithValidatorInfo/{chain}/{address}", UnbondingByAddressWithValidatorInfo)
	r.GET("/ValidatorsByAddressWithValidatorRanks/{chain}/{address}", ValidatorsByAddressWithValidatorRanks)
	r.GET("/DelegationsByAddressWithValidatorRanks/{chain}/{address}", DelegationsByAddressWithValidatorRanks)
	r.GET("/RemainingEpochs", RemainingEpochs)
}
