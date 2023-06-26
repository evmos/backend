// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"encoding/json"
	"sort"

	sdkmath "cosmossdk.io/math"
	"github.com/tharsis/dashboard-backend/internal/v1/db"
	"github.com/valyala/fasthttp"
)

// Types

// TODO: add all the others properties from the github schema
type ValidatorsRegistery struct {
	OperatorAddress string `json:"operator_address"`
}

type ConsensusKey struct {
	TypeURL string `json:"type_url"`
	Value   string `json:"value"`
}

type Description struct {
	Moniker         string `json:"moniker"`
	Identity        string `json:"identity"`
	Website         string `json:"website"`
	SecurityContact string `json:"security_contact"`
	Details         string `json:"details"`
}

type CommissionRate struct {
	Rate          string `json:"rate"`
	MaxRate       string `json:"max_rate"`
	MaxChangeRate string `json:"max_change_rate"`
}

type Commission struct {
	CommissionRate CommissionRate `json:"commission_rates"`
	UpdateTime     string         `json:"update_time"`
}

type Validator struct {
	OperatorAddress   string       `json:"operator_address"`
	ConsensusKey      ConsensusKey `json:"consensus_pubkey"`
	Jailed            bool         `json:"jailed"`
	Status            string       `json:"status"`
	Tokens            string       `json:"tokens"`
	DelegatorShares   string       `json:"delegator_shares"`
	Description       Description  `json:"description"`
	UnbondingHeight   string       `json:"unbonding_height"`
	UnbondingTime     string       `json:"unbonding_time"`
	Commission        Commission   `json:"commission"`
	MinSelfDelegation string       `json:"min_self_delegation"`
	Rank              int          `json:"rank"`
}

type ValidatorAPIResponse struct {
	Validators []Validator `json:"validators"`
	Pagination Pagination  `json:"pagination"`
}

func AllValidators(ctx *fasthttp.RequestCtx) {
	if validators, err := db.RedisGetAllValidators("EVMOS"); err == nil {
		sendResponse("{\"values\":"+validators+"}", err, ctx)
		return
	}

	endpoint := BuildTwoParamEndpoint("/cosmos/staking/v1beta1/validators?", "pagination.limit=500")
	res, err := getRequestRest("EVMOS", endpoint)
	if err != nil {
		sendResponse(err.Error(), err, ctx)
	}

	var validatorsResponse ValidatorAPIResponse

	err = json.Unmarshal([]byte(res), &validatorsResponse)

	if err != nil {
		sendResponse(err.Error(), err, ctx)
	}

	sort.SliceStable(validatorsResponse.Validators, func(a int, b int) bool {
		valA, okA := sdkmath.NewIntFromString(validatorsResponse.Validators[a].Tokens)
		valB, okB := sdkmath.NewIntFromString(validatorsResponse.Validators[b].Tokens)
		if !okA || !okB {
			return false
		}
		return valA.GT(valB)
	})

	validators := make([]Validator, len(validatorsResponse.Validators))

	// Set the ranks
	for i := range validatorsResponse.Validators {
		validatorsResponse.Validators[i].Rank = i + 1
		validators[i] = validatorsResponse.Validators[i]
	}

	validatorsByte, err := json.Marshal(validators)
	if err != nil {
		sendResponse(err.Error(), err, ctx)
		return
	}

	validatorsJSON := string(validatorsByte)

	db.RedisSetAllValidators("EVMOS", validatorsJSON)
	sendResponse("{\"values\":"+validatorsJSON+"}", err, ctx)
}
