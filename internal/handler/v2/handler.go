package v2

import "github.com/tharsis/dashboard-backend/internal/numia"

type Handler struct {
	numiaRPCClient *numia.RPCClient
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