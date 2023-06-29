// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"encoding/json"
	"math"
	"sort"

	sdkmath "cosmossdk.io/math"

	"github.com/tharsis/dashboard-backend/internal/v1/db"
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
