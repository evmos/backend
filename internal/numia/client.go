// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package numia

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type RPCClient struct {
	domain string
	apiKey string
}

// NewRPCClient creates a new RPC client for the Numia API.
// The client is used to make RPC requests to the Numia API.
func NewRPCClient() (*RPCClient, error) {
	apiKey := os.Getenv("NUMIA_API_KEY")
	if apiKey == "" {
		//return nil, fmt.Errorf("NUMIA_API_KEY environment variable not set")
	}

	endpoint := os.Getenv("NUMIA_RPC_ENDPOINT")
	if endpoint == "" {
		//return nil, fmt.Errorf("NUMIA_RPC_ENDPOINT environment variable not set")
	}

	return &RPCClient{
		domain: endpoint,
		apiKey: apiKey,
	}, nil
}

// get makes a GET request to the Numia API.
func (c *RPCClient) get(url string, v any) error {
	req, err := http.NewRequest("GET", c.domain+url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err.Error())
	}

	// Set authorization header
	authHeader := fmt.Sprintf("Bearer %s", c.apiKey)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "application/json")

	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err.Error())
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return fmt.Errorf("error decoding response: %s", err.Error())
	}

	return nil
}
