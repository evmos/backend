// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package blockchain

import (
	sdkmath "cosmossdk.io/math"
	ibc "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
)

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
) *ibc.MsgTransfer {
	timeoutHeight := clienttypes.Height{RevisionNumber: revisionNumber, RevisionHeight: revisionHeight}

	// TODO: add parameter for memo
	return ibc.NewMsgTransfer(sourcePort, sourceChannel, SdkIntToCoin(amount, denom), sender, receiver, timeoutHeight, timeoutTimestamp, "")
}
