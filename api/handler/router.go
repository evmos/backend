// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package handler

import (
	"github.com/fasthttp/router"
	v1 "github.com/tharsis/dashboard-backend/api/handler/v1"
)

func (h *Handler) RegisterRoutes(r *router.Router) {
	r.GET("/status", h.Status)
	// v2 endpoints
	r.GET("/v2/height", h.v2.Height)
	r.GET("/v2/delegations/{address}", h.v2.DelegationsByAddress)
	r.GET("/v2/rewards/{address}", h.v2.RewardsByAddress)
	r.GET("/v2/vesting/{address}", h.v2.VestingByAddress)

	// Tx endpoints
	r.POST("/v2/tx/broadcast", h.v2.BroadcastTx)
	r.POST("/v2/tx/amino/broadcast", h.v2.BroadcastAminoTx)

	// v1 endpoints to be deprecated
	// NOTE: v1 endpoints do not have a /v1 prefix for backwards compatibility
	v1.RegisterRoutes(r)
}
