package rest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	apiAddress string
}

// NewClient returns a new instance of a RestClient.
// It takes a network string as an argument, which is used to query a valid API address from redis
// for the desired network.
func NewClient(network string) (*Client, error) {
	apiAddress, err := getAPIAddress(network)
	if err != nil {
		return nil, fmt.Errorf("error while getting api address from redis: %w", err)
	}
	return &Client{
		apiAddress: apiAddress,
	}, nil
}

func getAPIAddress(_ string) (string, error) {
	// TODO: implement query to get domain from redis or decide what do to do from now on
	return "http://localhost:1317", nil
}

// PostRequest defines a wrapper around an HTTP POST request with a provided URL and data.
// An error is returned if the request or reading the body fails.
func (c *Client) postRequest(url string, body []byte) ([]byte, error) {
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
