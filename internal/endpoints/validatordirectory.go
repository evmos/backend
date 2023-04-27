// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package endpoints

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/fasthttp/router"
	"github.com/tharsis/dashboard-backend/internal/db"
	"github.com/tharsis/dashboard-backend/internal/requester"
	"github.com/valyala/fasthttp"
)

// Types

// TODO: add all the others properties from the github schema
type ValidatorsRegistery struct {
	OperatorAddress string `json:"operator_address"`
}

type ValidatorAPIResponse struct {
	Validators []Validator `json:"validators"`
	Pagination Pagination  `json:"pagination"`
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

// Sorters
func randomSort(array []Validator) []Validator {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(array), func(i, j int) { array[i], array[j] = array[j], array[i] })
	return array
}

func nameSort(array []Validator) []Validator {
	sort.Slice(array, func(i, j int) bool {
		return strings.ToLower(array[i].Description.Moniker) > strings.ToLower(array[j].Description.Moniker)
	})
	return array
}

func powerSort(array []Validator) []Validator {
	sort.Slice(array, func(i, j int) bool {
		a := new(big.Int)
		a, ok := a.SetString(array[i].Tokens, 10)
		if !ok {
			log.Default().Println("Error converting tokens: ", array[i].Tokens)
			a = big.NewInt(0)
		}

		b := new(big.Int)
		b, ok = b.SetString(array[j].Tokens, 10)
		if !ok {
			log.Default().Println("Error converting tokens: ", array[j].Tokens)
			b = big.NewInt(0)
		}

		return a.Cmp(b) == 1
	})

	return array
}

// Endpoints
func ValidatorDirectory(ctx *fasthttp.RequestCtx) {
	if val, err := db.RedisGetValidatorDirectory(); err == nil {
		sendResponse(val, nil, ctx)
		return
	}

	val, err := requester.GetValidatorDirectory()
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	values := []string{}
	for _, v := range val {
		values = append(values, v.Content)
	}

	res := "{\"values\":[" + strings.Join(values, ",") + "]}"

	db.RedisSetValidatorDirectory(res)

	sendResponse(res, nil, ctx)
}

func ValidatorDirectoryNotListed(ctx *fasthttp.RequestCtx) {
	// TODO: validate status string
	status := paramToString("status", ctx)
	sort := "random"
	res, err := ValidatorDirectoryNotListedSorted(status, sort)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}
	sendResponse(res, nil, ctx)
}

func ValidatorDirectoryNotListedWithFilter(ctx *fasthttp.RequestCtx) {
	// TODO: validate status string
	status := paramToString("status", ctx)
	sort := paramToString("sort", ctx)
	sort = strings.ToLower(sort)
	res, err := ValidatorDirectoryNotListedSorted(status, sort)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}
	sendResponse(res, nil, ctx)
}

func ValidatorDirectoryNotListedSorted(status string, sort string) (string, error) {
	if val, err := db.RedisGetValidatorDirectoryNoListed(status, sort); err == nil {
		return val, nil
	}

	val, err := requester.GetValidatorDirectory()
	if err != nil {
		return "", nil
	}

	values := make(map[string]interface{})
	for _, v := range val {
		var m ValidatorsRegistery
		err = json.Unmarshal([]byte(v.Content), &m)
		if err != nil {
			continue
		}
		values[m.OperatorAddress] = nil
	}

	validators, err := GetValidators(status, "EVMOS")
	if err != nil || strings.Contains(validators, "error") {
		return "", fmt.Errorf("error getting validators info")
	}

	var validatorsObject ValidatorAPIResponse
	err = json.Unmarshal([]byte(validators), &validatorsObject)

	if err != nil {
		return "", fmt.Errorf("error getting validators info")
	}

	res := []Validator{}
	// Remove validators in the directory list
	for _, validator := range validatorsObject.Validators {
		_, isListed := values[validator.OperatorAddress]
		if !isListed {
			res = append(res, validator)
		}
	}

	// TODO: add recently joined, uptime order

	switch sort {
	case "random":
		res = randomSort(res)
	case "name":
		res = nameSort(res)
	case "power":
		res = powerSort(res)
	}

	resJSON, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("error getting validators info")
	}
	resString := "{\"values\":" + string(resJSON) + "}"
	db.RedisSetValidatorDirectoryNoListed(status, sort, resString)
	return resString, nil
}

func AddValidatorDirectoryRoutes(r *router.Router) {
	r.GET("/ValidatorDirectory", ValidatorDirectory)
	r.GET("/ValidatorDirectory/NotListed/{status}", ValidatorDirectoryNotListed)
	r.GET("/ValidatorDirectory/NotListed/{status}/{sort}", ValidatorDirectoryNotListedWithFilter)
}
