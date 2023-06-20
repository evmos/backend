// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tharsis/dashboard-backend/internal/blockchain"
	"github.com/valyala/fasthttp"
)

type StakingInfoResponse struct {
	Delegations   []DelegationResponse   `json:"delegations"`
	Undelegations []UnbondingResponse    `json:"undelegations"`
	Rewards       StakingRewardsResponse `json:"rewards"`
}

type StakingReward struct {
	ValidatorAddress string `json:"validator_address"`
	Reward           []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"reward"`
}

type StakingRewardsResponse struct {
	Rewards []StakingReward `json:"rewards"`
	Total   []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"total"`
}

func TotalStakingByAddress(ctx *fasthttp.RequestCtx) {
	address := paramToString("address", ctx)

	addressSplitted := strings.Split(address, "evmos")
	if len(addressSplitted) < 2 && addressSplitted[0] != "evmos" {
		sendResponse("", fmt.Errorf("invalid wallet format"), ctx)
		return
	}

	endpoint := buildThreeParamEndpoint("/cosmos/staking/v1beta1/delegations/", address, "?pagination.limit=200")
	val, err := getRequestRest("EVMOS", endpoint)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	var stakingResponse blockchain.DelegationResponses
	_ = json.Unmarshal([]byte(val), &stakingResponse)

	totalStaked := blockchain.GetTotalStake(stakingResponse)

	res := "{\"value\":\"" + totalStaked + "\"}"
	sendResponse(res, err, ctx)
}

func StakingInfo(ctx *fasthttp.RequestCtx) {
	address := paramToString("address", ctx)

	delegationsURL := buildThreeParamEndpoint("/cosmos/staking/v1beta1/delegations/", address, "?pagination.limit=150")
	delegationsRes, err := getRequestRest("EVMOS", delegationsURL)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	var delegation DelegationResponsesResponse
	err = json.Unmarshal([]byte(delegationsRes), &delegation)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	undelegationsURL := buildThreeParamEndpoint("/cosmos/staking/v1beta1/delegators/", address, "/unbonding_delegations")
	undelegationsRes, err := getRequestRest("EVMOS", undelegationsURL)
	if err != nil {
		sendResponse("unable to get delegations", err, ctx)
		return
	}

	var unbodings UnbondingResponseAPI
	err = json.Unmarshal([]byte(undelegationsRes), &unbodings)
	if err != nil {
		sendResponse("unables to get unbonding data", err, ctx)
		return
	}

	valMap, err := GetValidatorsWithNoFilter("EVMOS")
	if err != nil {
		sendResponse("unable to get validators data", err, ctx)
		return
	}

	delegationsData := []DelegationResponse{}

	for _, d := range delegation.DelegationResponse {
		_, exists := valMap[d.Delegation.ValidatorAddress]
		if exists {
			d.Delegation.Validator = valMap[d.Delegation.ValidatorAddress]
		}
		delegationsData = append(delegationsData, d)
	}

	undelegationsData := []UnbondingResponse{}

	for _, u := range unbodings.UnbondingResponses {
		_, exists := valMap[u.ValidatorAddress]
		if exists {
			u.Validator = valMap[u.ValidatorAddress]
		}
		undelegationsData = append(undelegationsData, u)
	}

	endpoint := buildThreeParamEndpoint("/cosmos/distribution/v1beta1/delegators/", address, "/rewards")

	rewardsRes, err := getRequestRest("EVMOS", endpoint)
	if err != nil {
		sendResponse("unable to get rewards data", err, ctx)
		return
	}

	var rewards StakingRewardsResponse
	err = json.Unmarshal([]byte(rewardsRes), &rewards)
	if err != nil {
		sendResponse("unable to get rewards data", err, ctx)
		return
	}

	stakingInfo := StakingInfoResponse{
		Delegations:   delegationsData,
		Undelegations: undelegationsData,
		Rewards:       rewards,
	}

	res, err := json.Marshal(stakingInfo)
	if err != nil {
		sendResponse("unable to get validators data", err, ctx)
		return
	}

	sendResponse(string(res), nil, ctx)
}
