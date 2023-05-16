// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v2

import "github.com/tharsis/dashboard-backend/internal/numia"

type Handler struct {
	numiaRPCClient *numia.RPCClient // client to make RPC queries to Numia
}

func NewHandler() (*Handler, error) {
	numiaRPCClient, err := numia.NewRPCClient()
	if err != nil {
		return nil, err
	}

	return &Handler{
		numiaRPCClient: numiaRPCClient,
	}, nil
}
