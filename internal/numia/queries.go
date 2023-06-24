// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)
package numia

type HeightResponse struct {
	LatestBlockHash   string `json:"latestBlockHash"`
	LatestBlockHeight string `json:"latestBlockHeight"`
	LatestBlockTime   string `json:"latestBlockTime"`
}

// QueryHeight queries the height of the latest block on the EVMOS blockchain.
// URL: "https://evmos.numia.xyz/height"
func (c *RPCClient) QueryHeight() (*HeightResponse, error) {
	// Unmarshal response into struct
	var data HeightResponse
	if err := c.get("/height", &data); err != nil {
		return nil, err
	}

	return &data, nil
}

type DelegationResponse struct {
	ValidatorAddress string `json:"validatorAddress"`
	Delegated        struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"delegated"`
	Unclaimed struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"unclaimed"`
}

// QueryDelegations queries the delegations of the requested address.
// It handles both Hex and Bech32 addresses.
// URL: "https://evmos.numia.xyz/evmos/delegations"
func (c *RPCClient) QueryDelegations(address string) ([]DelegationResponse, error) {
	var data []DelegationResponse
	if err := c.get("/evmos/delegations/"+address, &data); err != nil {
		return nil, err
	}

	return data, nil
}

type RewardsResponse struct {
	Month                 string  `json:"month"`
	Address               string  `json:"address"`
	WithdrawnRewardsUsd   float64 `json:"withdrawn_rewards_usd"`
	WithdrawnRewardsEvmos float64 `json:"withdrawn_rewards_evmos"`
}

// QueryRewards queries the rewards of the requested address.
// It handles both Hex and Bech32 addresses.
// URL: "https://evmos.numia.xyz/evmos/rewards"
func (c *RPCClient) QueryRewards(address string) ([]RewardsResponse, error) {
	var data []RewardsResponse
	if err := c.get("/evmos/rewards/"+address, &data); err != nil {
		return nil, err
	}

	return data, nil
}

type VestingAccountResponse struct {
	Locked []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"locked"`
	Unvested []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"unvested"`
	Vested []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"vested"`
	Account struct {
		Type               string `json:"@type"`
		BaseVestingAccount struct {
			BaseAccount struct {
				Address       string `json:"address"`
				PubKey        string `json:"pub_key"`
				AccountNumber string `json:"account_number"`
				Sequence      string `json:"sequence"`
			} `json:"base_account"`
			OriginalVesting []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"original_vesting"`
			EndTime string `json:"end_time"`
		} `json:"base_vesting_account"`
		FunderAddress string `json:"funder_address"`
		StartTime     string `json:"start_time"`
		LockupPeriods []struct {
			Length string `json:"length"`
			Amount []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"amount"`
		} `json:"lockup_periods"`
		VestingPeriods []struct {
			Length string `json:"length"`
			Amount []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"amount"`
		} `json:"vesting_periods"`
	} `json:"account"`
}

// QueryVestingAccount queries the vesting account information of the requested address.
// It handles both Hex and Bech32 addresses.
// URL: "https://evmos.numia.xyz/evmos/account/{address}/vesting_balances"
func (c *RPCClient) QueryVestingAccount(address string) (VestingAccountResponse, error) {
	var data VestingAccountResponse
	if err := c.get("/evmos/account/"+address+"/vesting_balances", &data); err != nil {
		return VestingAccountResponse{}, err
	}

	return data, nil
}
