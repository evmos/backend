// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package endpoints

import (
	"github.com/fasthttp/router"
)

func CreateRouter() *router.Router {
	r := router.New()
	AddCustomRoutes(r)
	AddProxyRoutes(r)
	AddPriceRoutes(r)
	AddERC20Routes(r)
	AddUtilsRoutes(r)
	AddNetworkRoutes(r)
	AddValidatorDirectoryRoutes(r)
	AddAirtableRoutes(r)
	AddTransactionRoutes(r)
	AddBalancesRoutes(r)
	AddStakingRoutes(r)
	AddGovernanceRoutes(r)
	AddValidatorsRoutes(r)
	return r
}
