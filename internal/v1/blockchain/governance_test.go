// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package blockchain

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestConvertV1ToV1Beta(t *testing.T) {
	v1Payload := ` 
	{
		"id": "113",
		"messages": [
		  {
			"@type": "/cosmos.gov.v1.MsgExecLegacyContent",
			"content": {
			  "@type": "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal",
			  "title": "[ECP-2A] Formalizing and Funding the Governance Council Workstream (Rev. 2)",
			  "description": "# ECP-2A: Formalizing",
			  "recipient": "evmos1c0z326g3haflz2gnk68q2vsfv5mtxpsqgqgd0v",
			  "amount": [
				{
				  "denom": "aevmos",
				  "amount": "950000000000000000000000"
				}
			  ]
			},
			"authority": "evmos10d07y265gmmuvt4z0w9aw880jnsr700jcrztvm"
		  }
		],
		"status": "PROPOSAL_STATUS_VOTING_PERIOD",
		"final_tally_result": {
		  "yes_count": "0",
		  "abstain_count": "0",
		  "no_count": "0",
		  "no_with_veto_count": "0"
		},
		"submit_time": "2023-01-30T16:25:51.456995277Z",
		"deposit_end_time": "2023-02-02T16:25:51.456995277Z",
		"total_deposit": [
		  {
			"denom": "aevmos",
			"amount": "3500100000000000000000"
		  }
		],
		"voting_start_time": "2023-01-30T18:28:40.002375697Z",
		"voting_end_time": "2023-02-04T18:28:40.002375697Z",
		"metadata": ""
	}`
	v1betaPayload := `
	{
		"proposal_id": "113",
		"content": {
		  "@type": "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal",
		  "title": "[ECP-2A] Formalizing and Funding the Governance Council Workstream (Rev. 2)",
		  "description": "# ECP-2A: Formalizing",
		  "recipient": "evmos1c0z326g3haflz2gnk68q2vsfv5mtxpsqgqgd0v",
		  "amount": [
			{
			  "denom": "aevmos",
			  "amount": "950000000000000000000000"
			}
		  ]
		},
		"status": "PROPOSAL_STATUS_VOTING_PERIOD",
		"final_tally_result": {
		  "yes": "0",
		  "abstain": "0",
		  "no": "0",
		  "no_with_veto": "0"
		},
		"submit_time": "2023-01-30T16:25:51.456995277Z",
		"deposit_end_time": "2023-02-02T16:25:51.456995277Z",
		"total_deposit": [
		  {
			"denom": "aevmos",
			"amount": "3500100000000000000000"
		  }
		],
		"voting_start_time": "2023-01-30T18:28:40.002375697Z",
		"voting_end_time": "2023-02-04T18:28:40.002375697Z"
	}
	`
	var v1Proposal V1GovernanceProposal

	_ = json.Unmarshal([]byte(v1Payload), &v1Proposal)

	convertedProposal, err := ConvertV1ToV1Beta(v1Proposal)
	if err != nil {
		t.Fatalf("Error while converting v1 to v1beta1")
	}

	var expectedConversion GovernanceProposal
	_ = json.Unmarshal([]byte(v1betaPayload), &expectedConversion)

	isEqual := reflect.DeepEqual(convertedProposal, expectedConversion)

	if !isEqual {
		t.Fatalf("Incorrect v1 to v1beta1 conversion")
	}

	v1EmptyPayload := ` 
	{
		"id": "113",
		"messages": [],
		"status": "PROPOSAL_STATUS_VOTING_PERIOD",
		"final_tally_result": {
		  "yes_count": "0",
		  "abstain_count": "0",
		  "no_count": "0",
		  "no_with_veto_count": "0"
		},
		"submit_time": "2023-01-30T16:25:51.456995277Z",
		"deposit_end_time": "2023-02-02T16:25:51.456995277Z",
		"total_deposit": [
		  {
			"denom": "aevmos",
			"amount": "3500100000000000000000"
		  }
		],
		"voting_start_time": "2023-01-30T18:28:40.002375697Z",
		"voting_end_time": "2023-02-04T18:28:40.002375697Z",
		"metadata": ""
	}`

	var v1EmptyProposal V1GovernanceProposal

	_ = json.Unmarshal([]byte(v1EmptyPayload), &v1EmptyProposal)

	_, err = ConvertV1ToV1Beta(v1EmptyProposal)
	if err != nil {
		if err.Error() != "error unable to converts proposal with zero messages" {
			t.Fatalf("Incorrect error while converting v1 to v1beta1")
		}
	} else {
		t.Fatalf("Incorrect error while converting v1 to v1beta1")
	}
}
