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

	"github.com/tharsis/dashboard-backend/go-crons/endpoints/constants"
	"github.com/tharsis/dashboard-backend/go-crons/endpoints/models"
	"github.com/tharsis/dashboard-backend/internal/requester"
)

var jwg sync.WaitGroup

func PingJrpc(endpoint string, c chan models.Endpoint) {
	defer jwg.Done()

	transactionURL := fmt.Sprintf("%s/tx?hash=0x0000000000000000000000000000000000000000000000000000000000000000", endpoint)

	// make request
	resp, err := requester.Client.Get(transactionURL)
	if err != nil {
		e := models.Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}

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

	var jsonTransactionRes models.JrpcTransactionErrorResponse
	_ = json.Unmarshal(body, &jsonTransactionRes)

	if strings.Contains(jsonTransactionRes.Error.Data, constants.IndexingDisabledError) {
		e := models.Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}

	// record start time to measure latency
	start := time.Now()

	// url := fmt.Sprintf("%s/status", endpoint)
	url := fmt.Sprintf("%s/status", endpoint)

	// make request
	resp, err = requester.Client.Get(url)

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

	body, err = io.ReadAll(resp.Body)

	if err != nil || len(string(body)) == 0 {
		e := models.Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
		return
	}

	var jsonRes models.JrpcStatusResponse
	_ = json.Unmarshal(body, &jsonRes)

	if jsonRes.Result.NodeInfo.Other.TxIndex == "on" {

		height, err := strconv.Atoi(jsonRes.Result.SyncInfo.LatestBlockHeight)
		if err != nil {
			height = -1
		}

		e := models.Endpoint{
			URL:     endpoint,
			Latency: duration,
			Height:  height,
		}

		c <- e

	} else {
		e := models.Endpoint{
			URL:     endpoint,
			Latency: -1,
			Height:  -1,
		}
		c <- e
	}
}

func ProcessJrpc(jrpcEndpoints []string) []models.Endpoint {
	// create a channel to receive results for each jrpc endpoint
	jrpcChannel := make(chan models.Endpoint, len(jrpcEndpoints))

	for _, v := range jrpcEndpoints {
		// ping jrpc endpoint & get results in a goroutine
		go PingJrpc(v, jrpcChannel)
		// add goroutine to jrpc wait group
		jwg.Add(1)
	}

	jrpcResults := make([]models.Endpoint, 0)

	done := make(chan struct{})

	// Loop over values sent via channel.
	// This has to be as a separate goroutine in order to keep channel listener open
	// while waiting for all endpoints to be pinged & processed
	go func() {
		for r := range jrpcChannel {
			jrpcResults = append(jrpcResults, r)
		}
		close(done)
	}()
	// wait for all endpoints to be pinged & processed
	jwg.Wait()
	// close channel
	close(jrpcChannel)
	// wait for all values to be read
	<-done

	return jrpcResults
}
