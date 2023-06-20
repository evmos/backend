// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"github.com/fasthttp/router"
)

func RegisterRoutes(r *router.Router) {
	r.GET("/RemainingEpochs", RemainingEpochs)
	r.GET("/VoteRecord/{chain}/{proposal_id}/{address}", VoteRecord)
	r.GET("/Epochs/{chain}", Epochs)
	r.GET("/BalanceByDenom/{chain}/{address}/{denom:*}", BalanceByDenom)
	r.GET("/TxStatus/{chain}/{tx_hash}", TxStatus)
	r.POST("/broadcast", Broadcast)
	r.POST("/simulate", Simulate)
	r.GET("/ERC20ModuleBalance", ERC20ModuleEmptyBalance)
	r.GET("/ERC20ModuleBalance/{evmos_address}/{eth_address}", ERC20ModuleBalance)
	r.GET("/isIBCExecuted/{tx_hash}/{chain}", isIBCExecuted)
	r.GET("/NetworkConfig", NetworkConfig)
	r.GET("/NetworkConfig/{name}", NetworkConfigByName)
	r.GET("/Announcements", GetAnnouncements)
	r.POST("/broadcastEip712", BroadcastMetamask)
	r.POST("/ibcTransfer", IBCTransfer)
	r.POST("/convertCoin", ConvertCoin)
	r.POST("/convertERC20", ConvertERC20)
	r.POST("/broadcastAmino", BroadcastAmino)
	r.POST("/delegate", Delegate)
	r.POST("/undelegate", Undelegate)
	r.POST("/redelegate", Redelegate)
	r.POST("/rewards", Rewards)
	r.POST("/vote", Vote)
	r.POST("/cancelUndelegation", CancelUndelegation)
	r.GET("/BalanceByNetworkAndDenom/{chain}/{token}/{address}", BalanceByNetworkAndDenom)
	r.GET("/EVMOSIBCBalance/{chain}/{address}", EVMOSIBCBalance)
	r.GET("/totalStakedByAddress/{address}", TotalStakingByAddress)
	r.GET("/stakingInfo/{address}", StakingInfo)
	r.GET("/V1Proposals", V1GovernanceProposals)
	r.GET("/AllValidators", AllValidators)
}
