package rest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/tharsis/dashboard-backend/internal/v1/db"
)

type Client struct {
	nodesEndpoints []string
	network        string
}

// NewClient returns a new instance of a RestClient.
// It takes a network string as an argument, which is used to collect available REST node's endpoints from redis
// for the desired network.
func NewClient(network string) (*Client, error) {
	nodes, err := getAvailableNodes(network)
	if err != nil {
		return nil, fmt.Errorf("error while getting available endpoints from redis: %w", err)
	}
	return &Client{
		nodesEndpoints: nodes,
		network:        network,
	}, nil
}

// post defines a wrapper around an HTTP POST request with a provided URL and body.
// An error is returned if the request or reading the body fails.
func (c *Client) post(endpoint string, body []byte) ([]byte, error) {
	res, err := c.requestWithRetries("POST", endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("error while making post request: %w", err)
	}

	bz, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return bz, nil
}

// get defines a wrapper around an HTTP GET request with a provided URL.
// An error is returned if the request or reading the body fails.
func (c *Client) get(endpoint string) ([]byte, error) {
	res, err := c.requestWithRetries("GET", endpoint, []byte{})
	if err != nil {
		return nil, fmt.Errorf("error while making get request: %w", err)
	}

	bz, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return bz, nil
}

type BadRequestError struct {
	Message string `json:"message"`
}

// postRequestWithRetries performs a POST request to the provided URL with the provided body.
// It will retry the request with the next available node if the request fails.
func (c *Client) requestWithRetries(method string, endpoint string, body []byte) (*http.Response, error) {
	// TODO: this should be in a config file
	client := http.Client{
		Timeout: time.Second * 5,
	}

	var errorMessages []string
	for i := range c.nodesEndpoints {
		queryURL := joinURL(c.nodesEndpoints[i], endpoint)

		var resp *http.Response
		var err error

		if method == "POST" {
			resp, err = client.Post(queryURL, "application/json", bytes.NewBuffer(body))
		} else {
			resp, err = client.Get(queryURL)
		}

		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil // success, no need to retry
		}

		// Collect errors in case no endpoint is available
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("node %v error: %v", c.nodesEndpoints[i], err))
		} else {
			errorMessages = append(errorMessages, fmt.Sprintf("node %v status code: %v", c.nodesEndpoints[i], resp.StatusCode))
		}
	}

	return nil, fmt.Errorf(
		"failed to post request at endpoint %v for network %v after %v attempts: %v",
		endpoint,
		c.network,
		len(c.nodesEndpoints),
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

// getAvailableNodes returns a list of available nodes for the provided network
// from redis.
func getAvailableNodes(network string) ([]string, error) {
	// If env variable env == "local" then the only option is localhost
	env := os.Getenv("ENV")
	if env == "local" {
		return []string{"http://localhost:1317"}, nil
	}

	endpoints, err := db.RedisGetEndpoints(network, "rest")
	if err != nil {
		return nil, err
	}

	return endpoints, nil
}
