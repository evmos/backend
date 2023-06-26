package rest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tharsis/dashboard-backend/internal/db"
)

type Client struct {
	nodes   []string
	network string
}

// NewClient returns a new instance of a RestClient.
// It takes a network string as an argument, which is used to collect available REST node's endpoints from redis
// for the desired network.
func NewClient(network string) (*Client, error) {
	nodes, err := getAvailableNodes(network)
	if err != nil {
		return nil, fmt.Errorf("error while getting available endpoints: %w", err)
	}
	return &Client{
		nodes:   nodes,
		network: network,
	}, nil
}

// post defines a wrapper around an HTTP POST request with a provided URL and body.
// An error is returned if the request or reading the body fails.
func (c *Client) post(endpoint string, body []byte) ([]byte, error) {
	res, err := c.postRequestWithRetries(endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("error while making post request: %w", err)
	}

	bz, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return bz, nil
}

// postRequestWithRetries performs a POST request to the provided URL with the provided body.
// It will retry the request with the next available node if the request fails.
func (c *Client) postRequestWithRetries(endpoint string, body []byte) (*http.Response, error) {
	maxRetries := len(c.nodes)
	// TODO: this should be in a config file
	client := http.Client{
		Timeout: time.Second * 5,
	}

	var errorMessages []string
	for i := 0; i < maxRetries; i++ {
		queryURL := joinURL(c.nodes[i], endpoint)
		resp, err := client.Post(queryURL, "application/json", bytes.NewBuffer(body))
		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil // success, no need to retry
		}

		// Collect errors in case no endpoint is available
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("node %v error: %v", c.nodes[i], err))
		} else {
			errorMessages = append(errorMessages, fmt.Sprintf("node %v status code: %v", c.nodes[i], resp.StatusCode))
		}
	}

	return nil, fmt.Errorf(
		"failed to post request at endpoint %v for network %v after %v attempts: %v",
		endpoint,
		c.network,
		maxRetries,
		strings.Join(errorMessages, ", "),
	)
}

// joinURL joins a base URL and a query path to form a valid URL.
func joinURL(baseURL string, queryPath string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	u.Path = queryPath
	return u.String()
}

// getAvailableNodes returns a list of available nodes for the provided network.
func getAvailableNodes(network string) ([]string, error) {
	// If env variable env == "local" then the only option is localhost
	env := os.Getenv("ENV")
	if env == "local" {
		return []string{"http://localhost:1317"}, nil
	}

	// In production, query redis for the most up to date rest nodes for the network
	// TODO: there should be a redis query that returns the array in one request
	amountOfAvailableNodesPerNetwork := 4
	nodes := make([]string, amountOfAvailableNodesPerNetwork-1)
	for i := 1; i <= amountOfAvailableNodesPerNetwork; i++ {
		endpoint, err := db.RedisGetEndpoint(network, "rest", strconv.Itoa(i))
		if err != nil {
			return []string{}, fmt.Errorf("error while getting endpoint %v for network %v from redis: %w", i, network, err)
		}
		nodes[i-1] = endpoint
	}
	return nodes, nil
}
