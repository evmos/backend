package v2

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/cosmos/cosmos-sdk/types/tx"
)

type RestClient struct {
	apiAddress string
}

func NewRestClient(network string) (*RestClient, error) {
	apiAddress, err := getAPIAddress(network)
	if err != nil {
		return nil, fmt.Errorf("error while getting api address from redis: %w", err)
	}
	return &RestClient{
		apiAddress: apiAddress,
	}, nil
}

func getAPIAddress(_ string) (string, error) {
	// TODO: implement query to get domain from redis or decide what do to do from now on
	return "http://localhost:1317", nil
}

// PostRequest defines a wrapper around an HTTP POST request with a provided URL and data.
// An error is returned if the request or reading the body fails.
func (c *RestClient) post(url string, body []byte) ([]byte, error) {
	// join the node's address with the endpoint's path
	queryURL := fmt.Sprintf("%s/%s", c.apiAddress, url)

	res, err := http.Post(queryURL, "application/json", bytes.NewBuffer(body)) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("error while sending post request: %w", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	bz, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return bz, nil
}

// BroadcastTx broadcasts transaction bytes to a Tendermint node
// synchronously through the REST API.
func (c *RestClient) BroadcastTx(txBytes []byte) (tx.BroadcastTxResponse, error) {
	broadcastTxEndpoint := "cosmos/tx/v1beta1/txs"
	postResponse, err := c.post(broadcastTxEndpoint, txBytes)
	if err != nil {
		return tx.BroadcastTxResponse{}, err
	}
	encConfig := MakeEncodingConfig()

	jsonResponse := tx.BroadcastTxResponse{}
	err = encConfig.Codec.UnmarshalJSON(postResponse, &jsonResponse)
	if err != nil {
		return tx.BroadcastTxResponse{}, fmt.Errorf("error while unmarshalling response body: %w", err)
	}
	return jsonResponse, nil
}
