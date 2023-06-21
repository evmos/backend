// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"github.com/fasthttp/router"
)

func RegisterRoutes(r *router.Router) {
	// epoch
	r.GET("/RemainingEpochs", RemainingEpochs)
	r.GET("/Epochs/{chain}", Epochs)

	// announcements
	r.GET("/Announcements", GetAnnouncements)

	//config
	r.GET("/NetworkConfig", NetworkConfig)
	r.GET("/NetworkConfig/{name}", NetworkConfigByName)

	// ibc
	r.POST("/ibcTransfer", IBCTransfer)

	// bank
	r.GET("/BalanceByDenom/{chain}/{address}/{denom:*}", BalanceByDenom)
	r.GET("/BalanceByNetworkAndDenom/{chain}/{token}/{address}", BalanceByNetworkAndDenom)
	r.GET("/EVMOSIBCBalance/{chain}/{address}", EVMOSIBCBalance)

	// distribution
	r.POST("/rewards", Rewards)

	//tx
	r.GET("/isIBCExecuted/{tx_hash}/{chain}", isIBCExecuted)
	r.POST("/broadcastEip712", BroadcastMetamask)
	r.POST("/broadcastAmino", BroadcastAmino)
	r.POST("/broadcast", Broadcast)
	r.POST("/simulate", Simulate)
	r.GET("/TxStatus/{chain}/{tx_hash}", TxStatus)

	// staking
	r.GET("/totalStakedByAddress/{address}", TotalStakingByAddress)
	r.GET("/AllValidators", AllValidators)
	r.POST("/delegate", Delegate)
	r.POST("/undelegate", Undelegate)
	r.POST("/redelegate", Redelegate)
	r.GET("/stakingInfo/{address}", StakingInfo)
	r.POST("/cancelUndelegation", CancelUndelegation)

	//gov
	r.GET("/VoteRecord/{chain}/{proposal_id}/{address}", VoteRecord)
	r.GET("/V1Proposals", V1GovernanceProposals)
	r.POST("/vote", Vote)

	// erc20
	r.GET("/ERC20ModuleBalance", ERC20ModuleEmptyBalance)
	r.GET("/ERC20ModuleBalance/{evmos_address}/{eth_address}", ERC20ModuleBalance)
	r.POST("/convertCoin", ConvertCoin)
	r.POST("/convertERC20", ConvertERC20)
}
