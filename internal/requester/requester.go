// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package requester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tharsis/dashboard-backend/internal/db"
	"github.com/tharsis/dashboard-backend/internal/metrics"
)

var Client = http.Client{
	Timeout: 2 * time.Second,
}

var clientLongRequest = http.Client{
	Timeout: 6 * time.Second,
}

const BadRequestError = `{"error": "Bad Request"}`

// Right now is not being used because python is saving the prices
func GetRequestPrice(asset string, vsCurrency string) (string, error) {
	var sb strings.Builder
	sb.WriteString("https://api.coingecko.com/api/v3/simple/price?ids=")
	sb.WriteString(asset)
	sb.WriteString("&vs_currencies=")
	sb.WriteString(vsCurrency)

	resp, err := Client.Get(sb.String())
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil || len(string(body)) == 0 {
		return "", err
	}

	return string(body), nil
}

func MakeGetRequest(chain string, endpointType string, url string) (string, error) {
	i := 1
	for i < 4 {
		endpoint, err := db.RedisGetEndpoint(chain, endpointType, strconv.FormatInt(int64(i), 10))
		if err != nil {
			i++
			continue
		}

		var sb strings.Builder
		sb.WriteString(endpoint)
		sb.WriteString(url)

		resp, err := Client.Get(sb.String())
		if err != nil {
			i++
			continue
		}

		// Handle 404 responses from cosmos api, it's actually element not found
		if resp.StatusCode == 404 {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			// endpoint error
			if strings.Contains(string(body), "Cannot GET") {
				i++
				continue
			}
			// node element not found
			return `{"error": "Element not found"}`, nil
		}

		if resp.StatusCode == 400 {
			return BadRequestError, nil
		}
		if resp.StatusCode != 200 {
			i++
			continue
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)

		if err != nil || len(string(body)) == 0 {
			i++
			continue
		}

		return string(body), nil
	}

	metrics.Send(fmt.Sprintln("All endpoints failed to get response(GET): ", chain, url))
	return "", fmt.Errorf("all endpoints are down")
}

type status400Params struct {
	Code    int8   `json:"code"`
	Message string `json:"message"`
}

// Uses a bigger timeout for broadcast transactions
func MakeLongPostRequest(chain string, endpointType string, url string, param []byte) (string, error) {
	return makePostRequestInternal(chain, endpointType, url, param, clientLongRequest)
}

func MakePostRequest(chain string, endpointType string, url string, param []byte) (string, error) {
	return makePostRequestInternal(chain, endpointType, url, param, Client)
}

func makePostRequestInternal(chain string, endpointType string, url string, param []byte, httpClient http.Client) (string, error) {
	// Post requests are not using a second cache to avoid returning the incorrect value after submiting a transaction

	i := 1
	if endpointType == "web3" {
		// We are using the best bd endpoint as index 0
		// Right now they only support web3
		i = 0
	}

	for i < 4 {
		endpoint, err := db.RedisGetEndpoint(chain, endpointType, strconv.FormatInt(int64(i), 10))
		if err != nil {
			i++
			continue
		}

		var sb strings.Builder
		sb.WriteString(endpoint)
		sb.WriteString(url)
		// It has to be created here because Post delets the buffer
		body := bytes.NewBuffer(param)

		resp, err := httpClient.Post(sb.String(), "application/json", body)
		if err != nil {
			i++
			continue
		}

		// Handle 404 responses from cosmos api, it's actually element not found
		if resp.StatusCode == 404 {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			// endpoint error
			if strings.Contains(string(body), "Cannot POST") {
				i++
				continue
			}
			// node element not found
			return `{"error": "Element not found"}`, nil
		}
		// Handle 400 responses from api, the txBytes are incorrect
		if resp.StatusCode == 400 {

			defer resp.Body.Close()

			bodyResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				return BadRequestError, nil
			}

			if len(string(bodyResponse)) == 0 {
				return BadRequestError, nil
			}

			m := status400Params{}

			err = json.Unmarshal(bodyResponse, &m)
			if err != nil {
				return BadRequestError, nil
			}

			if m.Code != 0 {
				return "{\"error\": \"" + m.Message + "\"}", nil
			}

			return BadRequestError, nil
		}

		if resp.StatusCode == 500 {
			// Case: when you send a tx with an incorrect sequence.
			return `{"error": "Couldn't broadcast tx, please try again"}`, nil
		}

		// Only 200 and 404 are valid status code responses
		if resp.StatusCode != 200 && resp.StatusCode != 404 {
			i++
			continue
		}

		defer resp.Body.Close()

		bodyResponse, err := io.ReadAll(resp.Body)

		if err != nil || len(string(bodyResponse)) == 0 {
			i++
			continue
		}

		return string(bodyResponse), nil
	}

	metrics.Send(fmt.Sprintln("All endpoints failed to get response(POST): ", chain, url))
	return "", fmt.Errorf("all endpoints are down")
}

func MakePostGasPrice(url string) (string, error) {
	// make request
	payload := bytes.NewBuffer([]byte(`{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":1}`))

	resp, err := Client.Post(url, "application/json", payload)
	if err != nil {
		return BadRequestError, nil
	}

	// Only 200 is valid status code response
	if resp.StatusCode != 200 {
		return BadRequestError, nil
	}

	defer resp.Body.Close()
	bodyResponse, err := io.ReadAll(resp.Body)

	if err != nil || len(string(bodyResponse)) == 0 {
		return BadRequestError, nil
	}

	return string(bodyResponse), nil
}
