// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"encoding/json"

	"github.com/fasthttp/router"
	"github.com/tharsis/dashboard-backend/internal/blockchain"
	"github.com/tharsis/dashboard-backend/internal/db"

	"github.com/valyala/fasthttp"
)

func ProcessProposals(proposalsRes string, v1 bool) ([]byte, error) {
	var jsonProposalRes blockchain.V1GovernanceProposalsResponse
	err := json.Unmarshal([]byte(proposalsRes), &jsonProposalRes)
	if err != nil {
		return []byte{}, err
	}

	var proposalRes []byte

	if v1 {
		// Get current tally for proposals in voting period and overwrite final tally
		proposals, err := blockchain.GetV1ProposalsTally(jsonProposalRes.Proposals)
		if err != nil {
			return []byte{}, err
		}

		proposalRes, err = json.Marshal(proposals)
		if err != nil {
			return []byte{}, err
		}

	} else {
		// Get current tally for proposals in voting period and overwrite final tally
		proposals, err := blockchain.GetProposalsTally(jsonProposalRes.Proposals)
		if err != nil {
			return []byte{}, err
		}

		proposalRes, err = json.Marshal(proposals)
		if err != nil {
			return []byte{}, err
		}

	}

	return proposalRes, nil
}

// This endpoint returns a list of the latest 50 governance proposals. In order
// to support both v1 and v1beta1 versions this endpoint converts the v1 payload
// to be the same as the v1beta1 payload.
func GovernanceProposals(ctx *fasthttp.RequestCtx) {
	var proposalRes []byte
	if redisVal, err := db.RedisGetGovernanceProposals(); err == nil && redisVal != "null" {
		proposalRes = []byte(redisVal)
		if err != nil {
			sendResponse("Unable to fetch governance proposals", err, ctx)
			return
		}
	} else {
		endpoint := buildThreeParamEndpoint("/cosmos/gov/v1/proposals?pagination.limit=", "50", "&pagination.reverse=true")
		val, err := getRequestRest("EVMOS", endpoint)
		if err != nil {
			sendResponse("Unable to fetch governance proposals", err, ctx)
			return
		}

		// Process and convert v1 payload into v1beta1 payload version
		proposalRes, err = ProcessProposals(val, false)
		if err != nil {
			sendResponse("Unable to fetch governance proposals", err, ctx)
			return
		}

		db.RedisSetGovernanceProposals(string(proposalRes))
	}
	sendResponse(string(proposalRes), nil, ctx)
}

// nolint: revive
func V1GovernanceProposals(ctx *fasthttp.RequestCtx) {
	var proposalRes []byte
	if redisVal, err := db.RedisGetGovernanceV1Proposals(); err == nil && redisVal != "null" {
		proposalRes = []byte(redisVal)
		if err != nil {
			sendResponse("Unable to fetch governance proposals", err, ctx)
			return
		}
	} else {
		endpoint := buildThreeParamEndpoint("/cosmos/gov/v1/proposals?pagination.limit=", "50", "&pagination.reverse=true")
		val, err := getRequestRest("EVMOS", endpoint)
		if err != nil {
			sendResponse("Unable to fetch governance proposals", err, ctx)
			return
		}

		// Process and convert v1 payload into v1beta1 payload version
		proposalRes, err = ProcessProposals(val, true)
		if err != nil {
			sendResponse("Unable to fetch governance proposals", err, ctx)
			return
		}

		db.RedisSetGovernanceV1Proposals(string(proposalRes))
	}
	sendResponse(string(proposalRes), nil, ctx)
}

func AddGovernanceRoutes(r *router.Router) {
	r.GET("/Proposals", GovernanceProposals)
	r.GET("/V1Proposals", V1GovernanceProposals)
}
