// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tharsis/dashboard-backend/go-crons/endpoints/models"
	"github.com/tharsis/dashboard-backend/internal/v1/requester"
)

var rwg sync.WaitGroup

func PingNonTendermintRest(endpoint string, c chan models.Endpoint) {
	defer rwg.Done()

	url := fmt.Sprintf("%s/cosmos/auth/v1beta1/params", endpoint)

	// record start time to measure latency
	start := time.Now()

	// make request
	_, err := requester.Client.Get(url)
	if err != nil {
		e := models.Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}

	// compute latency
	duration := time.Since(start).Seconds()

	e := models.Endpoint{
		URL:     endpoint,
		Latency: duration,
		Height:  0,
	}

	c <- e
}

func PingRest(endpoint string, c chan models.Endpoint) {
	defer rwg.Done()

	url := fmt.Sprintf("%s/cosmos/base/tendermint/v1beta1/blocks/latest", endpoint)

	// record start time to measure latency
	start := time.Now()

	// make request
	resp, err := requester.Client.Get(url)
	if err != nil {
		e := models.Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}

	// compute latency
	duration := time.Since(start).Seconds()

	body, err := io.ReadAll(resp.Body)

	if err != nil || len(string(body)) == 0 {
		e := models.Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}

	// get block height from response
	var jsonRes models.RestResponse
	_ = json.Unmarshal(body, &jsonRes)

	height, err := strconv.Atoi(jsonRes.Block.Header.Height)
	if err != nil {
		height = -1
	}

	e := models.Endpoint{
		URL:     endpoint,
		Latency: duration,
		Height:  height,
	}

	c <- e
}

func ProcessRest(restEndpoints []string, chainIdentifier string) []models.Endpoint {
	// create a channel to receive results for each rest endpoint
	restChannel := make(chan models.Endpoint, len(restEndpoints))

	for _, v := range restEndpoints {
		// ping REST endpoint & get results in a goroutine
		if !strings.Contains(strings.ToUpper(chainIdentifier), "GRAVITY") {
			go PingRest(v, restChannel)
		} else {
			go PingNonTendermintRest(v, restChannel)
		}
		// add goroutine to rest wait group
		rwg.Add(1)
	}

	restResults := make([]models.Endpoint, 0)

	done := make(chan struct{})

	// Loop over values sent via channel.
	// This has to be as a separate goroutine in order to keep channel listener open
	// while waiting for all endpoints to be pinged & processed
	go func() {
		for r := range restChannel {
			restResults = append(restResults, r)
		}
		close(done)
	}()
	// wait for all endpoints to be pinged & processed
	rwg.Wait()
	// close channel
	close(restChannel)
	// wait for all values to be read
	<-done

	return restResults
}
