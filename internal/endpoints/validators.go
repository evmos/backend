// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package endpoints

import (
	"encoding/json"
	"sort"

	sdkmath "cosmossdk.io/math"
	"github.com/fasthttp/router"
	"github.com/tharsis/dashboard-backend/internal/db"
	"github.com/valyala/fasthttp"
)

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

func AddValidatorsRoutes(r *router.Router) {
	r.GET("/AllValidators", AllValidators)
}
