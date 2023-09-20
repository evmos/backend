// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"encoding/json"

	"github.com/tharsis/dashboard-backend/internal/v1/blockchain"
	"github.com/tharsis/dashboard-backend/internal/v1/db"

	"github.com/valyala/fasthttp"
)

func ProcessProposals(proposalsRes string, v1 bool) ([]byte, error) {
	var jsonProposalRes blockchain.V1GovernanceProposalsResponse
	err := json.Unmarshal([]byte(proposalsRes), &jsonProposalRes)
	if err != nil {
		return []byte{}, err
	}

	var filteredProposals []blockchain.V1GovernanceProposal
	for _, proposal := range jsonProposalRes.Proposals {
		if proposal.Status != "PROPOSAL_STATUS_DEPOSIT_PERIOD" {
			filteredProposals = append(filteredProposals, proposal)
		}
	}

	var proposalRes []byte

	if v1 {
		// Get current tally for proposals in voting period and overwrite final tally
		proposals, err := blockchain.GetV1ProposalsTally(filteredProposals)
		if err != nil {
			return []byte{}, err
		}

		proposalRes, err = json.Marshal(proposals)
		if err != nil {
			return []byte{}, err
		}

	} else {
		// Get current tally for proposals in voting period and overwrite final tally
		proposals, err := blockchain.GetProposalsTally(filteredProposals)
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

func V1GovernanceProposals(ctx *fasthttp.RequestCtx) { //nolint: revive
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
