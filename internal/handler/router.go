package handler

import (
	"github.com/fasthttp/router"
	"github.com/tharsis/dashboard-backend/internal/handler/v1"
)

func (h *Handler) RegisterRoutes(r *router.Router) {
	r.GET("/status", h.Status)
	// v2 endpoints
	r.GET("/v2/height", h.v2.Height)
	r.GET("/v2/delegations/{address}", h.v2.DelegationsForAddress)
    r.GET("/v2/rewards/{address}", h.v2.RewardsForAddress)

	// v1 endpoints
	// NOTE: v1 endpoints do not have a /v1 prefix for backwards compatibility
	v1.RegisterRoutes(r)
}
