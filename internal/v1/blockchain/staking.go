// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package blockchain

import (
	"math/big"
)

type DelegationResponses struct {
	DelegationResponses []DelegationResponse `json:"delegation_responses"`
}

type DelegationResponse struct {
	Delegation struct {
		DelegatorAddress string   `json:"delegator_address"`
		ValidatorAddress struct{} `json:"validator_address"`
	} `json:"delegation"`
	Balance struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"balance"`
}

func GetTotalStake(delegations DelegationResponses) string {
	totalBalance := new(big.Int)
	for _, v := range delegations.DelegationResponses {
		balance := new(big.Int)
		balance.SetString(v.Balance.Amount, 10)
		totalBalance = totalBalance.Add(totalBalance, balance)
	}
	return totalBalance.String()
}
