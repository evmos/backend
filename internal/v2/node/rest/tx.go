package rest

import (
	"fmt"

	"github.com/tharsis/dashboard-backend/internal/v2/encoding"

	"github.com/cosmos/cosmos-sdk/types/tx"
)

// All endpoints under /cosmos/tx/ path should be defined in this file

// BroadcastTx broadcasts transaction bytes to a Tendermint node
// through its REST API.
func (c *Client) BroadcastTx(txBytes []byte) (tx.BroadcastTxResponse, error) {
	broadcastTxEndpoint := "cosmos/tx/v1beta1/txs"
	postResponse, err := c.postRequest(broadcastTxEndpoint, txBytes)
	if err != nil {
		return tx.BroadcastTxResponse{}, err
	}

	encConfig := encoding.MakeEncodingConfig()
	jsonResponse := tx.BroadcastTxResponse{}
	err = encConfig.Codec.UnmarshalJSON(postResponse, &jsonResponse)
	if err != nil {
		return tx.BroadcastTxResponse{}, fmt.Errorf("error while unmarshalling response body: %w", err)
	}
	return jsonResponse, nil
}
