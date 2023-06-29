// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package blockchain

import (
	"fmt"

	sdkmath "cosmossdk.io/math"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type MessageSendParams struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
	Denom    string `json:"denom"`
}

func CreateMessageSend(sender string, receiver string, amount sdkmath.Int, denom string, prefix string) (sdk.Msg, error) {
	from, err := Bech32StringToAddress(sender, prefix)
	if err != nil {
		return &bankTypes.MsgSend{}, fmt.Errorf("error creating from address: %q", err)
	}
	to, err := Bech32StringToAddress(receiver, prefix)
	if err != nil {
		return &bankTypes.MsgSend{}, fmt.Errorf("error creating to address: %q", err)
	}

	msgSendSdk := bankTypes.NewMsgSend(from, to, SdkIntToCoins(amount, denom))

	return msgSendSdk, nil
}

func CreateMsgDelegate(amount sdkmath.Int, accountAddress string, validator string, denom string) (sdk.Msg, error) {
	delegateMsg := stakingtypes.MsgDelegate{
		DelegatorAddress: accountAddress,
		ValidatorAddress: validator,
		Amount:           SdkIntToCoin(amount, denom),
	}

	return &delegateMsg, nil
}

func CreateMsgUndelegate(amount sdkmath.Int, accountAddress string, validator string, denom string) (sdk.Msg, error) {
	undelegateMsg := stakingtypes.MsgUndelegate{
		DelegatorAddress: accountAddress,
		ValidatorAddress: validator,
		Amount:           SdkIntToCoin(amount, denom),
	}

	return &undelegateMsg, nil
}

func CreateMsgRedelegate(amount sdkmath.Int, accountAddress string, validator string, validatorDst string, denom string) (sdk.Msg, error) {
	beginRedelegateMsg := stakingtypes.MsgBeginRedelegate{
		DelegatorAddress:    accountAddress,
		ValidatorSrcAddress: validator,
		ValidatorDstAddress: validatorDst,
		Amount:              SdkIntToCoin(amount, denom),
	}

	return &beginRedelegateMsg, nil
}

func CreateMsgVote(proposalID int, option int, voter string) (sdk.Msg, error) {
	vo := govtypes.VoteOption(option)
	voteMsg := govtypes.MsgVote{
		ProposalId: uint64(proposalID),
		Option:     vo,
		Voter:      voter,
	}
	return &voteMsg, nil
}

func CreateMsgRewards(accountAddress string, validator string) (sdk.Msg, error) {
	beginRedelegateMsg := distributiontypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: accountAddress,
		ValidatorAddress: validator,
	}
	return &beginRedelegateMsg, nil
}

func CreateMsgCancelUndelegations(amount sdkmath.Int, accountAddress string, validator string, denom string, height int64) (sdk.Msg, error) {
	cancelUndelegationsMsg := stakingtypes.MsgCancelUnbondingDelegation{
		DelegatorAddress: accountAddress,
		ValidatorAddress: validator,
		Amount:           SdkIntToCoin(amount, denom),
		CreationHeight:   height,
	}

	return &cancelUndelegationsMsg, nil
}
