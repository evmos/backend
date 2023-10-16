// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package blockchain

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/tharsis/dashboard-backend/internal/v1/requester"
)

type V1GovernanceProposalsResponse struct {
	Proposals []V1GovernanceProposal `json:"proposals"`
}

type GovernanceProposal struct {
	ProposalID       string           `json:"proposal_id"`
	Content          ProposalContent  `json:"content"`
	Status           string           `json:"status"`
	FinalTallyResult FinalTallyResult `json:"final_tally_result"`
	SubmitTime       string           `json:"submit_time"`
	DepositEndTime   string           `json:"deposit_end_time"`
	TotalDeposit     []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"total_deposit"`
	VotingStartTime string `json:"voting_start_time"`
	VotingEndTime   string `json:"voting_end_time"`
}

type V1GovernanceProposal struct {
	ID               string              `json:"id"`
	Messages         []V1ProposalContent `json:"messages"`
	Status           string              `json:"status"`
	FinalTallyResult V1FinalTallyResult  `json:"final_tally_result"`
	SubmitTime       string              `json:"submit_time"`
	DepositEndTime   string              `json:"deposit_end_time"`
	TotalDeposit     []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"total_deposit"`
	VotingStartTime string `json:"voting_start_time"`
	VotingEndTime   string `json:"voting_end_time"`
	Title           string `json:"title"`
	Summary         string `json:"summary"`
}

type V1ProposalsResponse struct {
	Proposals   []V1GovernanceProposal `json:"proposals"`
	TallyParams V1TallyParams          `json:"tally_params"`
}

type V1FinalTallyResult struct {
	Yes        string `json:"yes_count"`
	No         string `json:"no_count"`
	Abstain    string `json:"abstain_count"`
	NoWithVeto string `json:"no_with_veto_count"`
}

type FinalTallyResult struct {
	Yes        string `json:"yes"`
	No         string `json:"no"`
	Abstain    string `json:"abstain"`
	NoWithVeto string `json:"no_with_veto"`
}

type ProposalContent struct {
	Type        string `json:"@type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Recipient   string `json:"recipient"`
	Amount      []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"amount"`
	Changes []struct {
		Subspace string `json:"subspace"`
		Key      string `json:"key"`
		Value    string `json:"value"`
	} `json:"changes"`
}

type V1ProposalContent struct {
	Type    string `json:"@type"`
	Content struct {
		Type        string `json:"@type"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Recipient   string `json:"recipient"`
		Amount      []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"amount"`
	} `json:"content"`
	Authority string `json:"authority"`
}

type V1TallyResponse struct {
	Tally V1FinalTallyResult `json:"tally"`
}

type V1TallyParams struct {
	Quorum        string `json:"quorum"`
	Threshold     string `json:"threshold"`
	VetoThreshold string `json:"veto_threshold"`
}

type V1TallyParamsResponse struct {
	TallyParams V1TallyParams `json:"tally_params"`
}

const ProposalStatusVotingPeriod = "PROPOSAL_STATUS_VOTING_PERIOD"

func GetGovURL(v1 []byte) string {
	useV1, err := strconv.ParseBool(string(v1))
	if err != nil {
		useV1 = false
	}
	url := "/cosmos/gov/v1beta1"
	if useV1 {
		url = "/cosmos/gov/v1"
	}
	return url
}

// TODO: tech debt - need to deprecate v1beta1 conversion on client to support v1 fully with multiple proposal msgs
func ConvertV1ToV1Beta(proposal V1GovernanceProposal) (GovernanceProposal, error) {
	if len(proposal.Messages) == 0 {
		return GovernanceProposal{}, errors.New("error unable to converts proposal with zero messages")
	}

	latestMessage := proposal.Messages[len(proposal.Messages)-1]

	v1beta1Prop := GovernanceProposal{
		ProposalID: proposal.ID,
		Content: ProposalContent{
			Type:        latestMessage.Content.Type,
			Title:       latestMessage.Content.Title,
			Description: latestMessage.Content.Description,
			Recipient:   latestMessage.Content.Recipient,
			Amount:      latestMessage.Content.Amount,
		},
		Status: proposal.Status,
		FinalTallyResult: FinalTallyResult{
			Yes:        proposal.FinalTallyResult.Yes,
			No:         proposal.FinalTallyResult.No,
			Abstain:    proposal.FinalTallyResult.Abstain,
			NoWithVeto: proposal.FinalTallyResult.NoWithVeto,
		},
		SubmitTime:      proposal.SubmitTime,
		DepositEndTime:  proposal.DepositEndTime,
		TotalDeposit:    proposal.TotalDeposit,
		VotingStartTime: proposal.VotingStartTime,
		VotingEndTime:   proposal.VotingEndTime,
	}
	return v1beta1Prop, nil
}

func GetProposalsTally(proposals []V1GovernanceProposal) ([]GovernanceProposal, error) {
	// We don't know the length of the proposals in voting period so we can't create an array with a fixed length
	proposalsWithTally := []GovernanceProposal{}
	for _, v := range proposals {
		// Get current proposal tally if proposal is in voting period
		if v.Status == ProposalStatusVotingPeriod {

			var sb strings.Builder
			sb.WriteString("/cosmos/gov/v1/proposals/")
			sb.WriteString(v.ID)
			sb.WriteString("/tally")
			endpoint := sb.String()

			val, err := requester.MakeGetRequest("EVMOS", "rest", endpoint)
			if err != nil {
				return nil, err
			}

			var jsonTallyRes V1TallyResponse
			_ = json.Unmarshal([]byte(val), &jsonTallyRes)
			v.FinalTallyResult = jsonTallyRes.Tally
		}
		// Convert v1 payload into v1beta1 payload
		v1beta1, err := ConvertV1ToV1Beta(v)
		if err != nil {
			return nil, err
		}

		proposalsWithTally = append(proposalsWithTally, v1beta1)

	}

	return proposalsWithTally, nil
}

func GetV1ProposalsTally(proposals []V1GovernanceProposal) (V1ProposalsResponse, error) {
	proposalsWithTally := []V1GovernanceProposal{}
	for _, v := range proposals {
		// Get current proposal tally if proposal is in voting period
		if v.Status == ProposalStatusVotingPeriod {

			var sb strings.Builder
			sb.WriteString("/cosmos/gov/v1/proposals/")
			sb.WriteString(v.ID)
			sb.WriteString("/tally")
			endpoint := sb.String()

			val, err := requester.MakeGetRequest("EVMOS", "rest", endpoint)
			if err != nil {
				return V1ProposalsResponse{}, err
			}

			var jsonTallyRes V1TallyResponse
			_ = json.Unmarshal([]byte(val), &jsonTallyRes)
			v.FinalTallyResult = jsonTallyRes.Tally
		}

		proposalsWithTally = append(proposalsWithTally, v)
	}

	pr := V1ProposalsResponse{
		Proposals: proposalsWithTally,
	}

	var sb strings.Builder
	sb.WriteString("/cosmos/gov/v1/params/tallying")
	endpoint := sb.String()

	val, err := requester.MakeGetRequest("EVMOS", "rest", endpoint)
	if err != nil {
		return V1ProposalsResponse{}, err
	}

	var jsonTallyParamsRes V1TallyParamsResponse
	_ = json.Unmarshal([]byte(val), &jsonTallyParamsRes)

	pr.TallyParams = jsonTallyParamsRes.TallyParams

	return pr, nil
}
