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
	ibctransfer "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
)

type MessageSendParams struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
	Denom    string `json:"denom"`
}

func CreateMessageSend(sender string, receiver string, amount sdkmath.Int, denom string, prefix string) (sdk.Msg, error) {
	from, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return &bankTypes.MsgSend{}, fmt.Errorf("error creating from address: %s", err)
	}
	to, err := sdk.AccAddressFromBech32(receiver)
	if err != nil {
		return &bankTypes.MsgSend{}, fmt.Errorf("error creating to address: %s", err)
	}

	msg := bankTypes.NewMsgSend(from, to, sdk.Coins{{Denom: denom, Amount: amount}})

	return msg, msg.ValidateBasic()
}

func CreateMsgDelegate(amount sdkmath.Int, accountAddress string, validator string, denom string) (sdk.Msg, error) {
	msg := stakingtypes.MsgDelegate{
		DelegatorAddress: accountAddress,
		ValidatorAddress: validator,
		Amount:           sdk.Coin{Denom: denom, Amount: amount},
	}

	return &msg, msg.ValidateBasic()
}

func CreateMsgUndelegate(amount sdkmath.Int, accountAddress string, validator string, denom string) (sdk.Msg, error) {
	msg := stakingtypes.MsgUndelegate{
		DelegatorAddress: accountAddress,
		ValidatorAddress: validator,
		Amount:           sdk.Coin{Denom: denom, Amount: amount},
	}

	return &msg, msg.ValidateBasic()
}

func CreateMsgRedelegate(amount sdkmath.Int, accountAddress string, validator string, validatorDst string, denom string) (sdk.Msg, error) {
	msg := stakingtypes.MsgBeginRedelegate{
		DelegatorAddress:    accountAddress,
		ValidatorSrcAddress: validator,
		ValidatorDstAddress: validatorDst,
		Amount:              sdk.Coin{Denom: denom, Amount: amount},
	}

	return &msg, msg.ValidateBasic()
}

func CreateMsgVote(proposalID int, option int, voter string) (sdk.Msg, error) {
	vo := govtypes.VoteOption(option)
	msg := govtypes.MsgVote{
		ProposalId: uint64(proposalID),
		Option:     vo,
		Voter:      voter,
	}
	return &msg, msg.ValidateBasic()
}

func CreateMsgRewards(accountAddress string, validator string) (sdk.Msg, error) {
	msg := distributiontypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: accountAddress,
		ValidatorAddress: validator,
	}
	return &msg, msg.ValidateBasic()
}

func CreateMsgCancelUndelegations(amount sdkmath.Int, accountAddress string, validator string, denom string, height int64) (sdk.Msg, error) {
	msg := stakingtypes.MsgCancelUnbondingDelegation{
		DelegatorAddress: accountAddress,
		ValidatorAddress: validator,
		Amount:           sdk.Coin{Denom: denom, Amount: amount},
		CreationHeight:   height,
	}

	return &msg, msg.ValidateBasic()
}

func CreateMsgTransfer(
	sourcePort string,
	sourceChannel string,
	amount sdkmath.Int,
	denom string,
	sender string,
	receiver string,
	revisionNumber uint64,
	revisionHeight uint64,
	timeoutTimestamp uint64,
	memo string,
) (sdk.Msg, error) {
	timeoutHeight := clienttypes.Height{RevisionNumber: revisionNumber, RevisionHeight: revisionHeight}

	msg := ibctransfer.NewMsgTransfer(sourcePort, sourceChannel, amount, denom, sender, receiver, timeoutHeight, timeoutTimestamp, "")
	return msg, msg.ValidateBasic()
}
