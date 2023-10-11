package rest

import (
	"encoding/json"
	"fmt"

	"github.com/evmos/evmos/v12/x/vesting/types"
	v1 "github.com/tharsis/dashboard-backend/api/handler/v1"
)

type VestingAccount struct {
	Account struct {
		Type               string                `json:"@type"`
		BaseVestingAccount v1.BaseVestingAccount `json:"base_vesting_account"`
		FunderAddress      string                `json:"funder_address"`
		StartTime          string                `json:"start_time"`
		LockupPeriods      []struct {
			Length string              `json:"length"`
			Amount []v1.BalanceElement `json:"amount"`
		} `json:"lockup_periods"`
		VestingPeriods []struct {
			Length string              `json:"length"`
			Amount []v1.BalanceElement `json:"amount"`
		} `json:"vesting_periods"`
	} `json:"account"`
}

type VestingByAddressResponse struct {
	VestingAccount
	types.QueryBalancesResponse
}

func (c *Client) GetVestingAccount(address string) (VestingByAddressResponse, error) {
	accountRes, err := c.get("/cosmos/auth/v1beta1/accounts/" + address)
	if err != nil {
		return VestingByAddressResponse{}, fmt.Errorf("error querying vesting account from RPC: %s", err.Error())
	}

	var account VestingAccount
	err = json.Unmarshal(accountRes, &account)
	if err != nil {
		return VestingByAddressResponse{}, fmt.Errorf("error decoding vesting account: %s", err.Error())
	}

	rewardsRes, err := c.get("/evmos/vesting/v2/balances/" + address)
	if err != nil {
		return VestingByAddressResponse{}, fmt.Errorf("error querying vesting balance from RPC: %s", err.Error())
	}

	var vestingBalance types.QueryBalancesResponse
	err = json.Unmarshal(rewardsRes, &vestingBalance)
	if err != nil {
		return VestingByAddressResponse{}, fmt.Errorf("error decoding vesting account: %s", err.Error())
	}

	res := VestingByAddressResponse{
		VestingAccount:        account,
		QueryBalancesResponse: vestingBalance,
	}

	return res, nil
}
