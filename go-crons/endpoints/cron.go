// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package main

import (
	"fmt"
	"sort"
	"sync"

	"github.com/tharsis/dashboard-backend/go-crons/endpoints/helpers"
	"github.com/tharsis/dashboard-backend/go-crons/endpoints/models"
	"github.com/tharsis/dashboard-backend/internal/db"
	"github.com/tharsis/dashboard-backend/internal/resources"
)

var (
	running = true
	wg      sync.WaitGroup
)

func sortEndpoints(endpoints []models.Endpoint) []models.Endpoint { //nolint:all
	sort.SliceStable(endpoints, func(i, j int) bool {
		if endpoints[i].Height == -1 {
			return true
		}

		if endpoints[i].Latency == -1 {
			return true
		}

		if endpoints[i].Height != endpoints[j].Height {
			return endpoints[i].Height < endpoints[j].Height
		}

		return endpoints[i].Latency > endpoints[j].Latency
	})

	return endpoints
}

func processNetwork(networkConfig resources.NetworkConfig) {
	// signal waiting group goroutine is done at the end of the function
	defer wg.Done()

	// get mainnet configuration
	config := resources.GetMainnetConfig(networkConfig)

	fmt.Printf("Processing %s network...\n", config.Identifier)

	// process REST endpoints
	restEndpoints := helpers.ProcessRest(config.Rest, config.Identifier)
	// sort endpoints based on latency and height
	sortEndpoints(restEndpoints)

	// process JRPC endpoints
	jrpcEndpoints := helpers.ProcessJrpc(config.Jrpc)
	// sort endpoints based on latency and height
	sortEndpoints(jrpcEndpoints)

	// process web3 endpoints if available
	var web3Endpoints []models.Endpoint
	if len(config.Web3) > 0 {
		web3Endpoints = helpers.ProcessWeb3(config.Web3)
		// sort endpoints based on latency and height
		sortEndpoints((web3Endpoints))
	}

	fmt.Printf("Finished processing %s network...\n", config.Identifier)

	// store rest result in redis
	if len(restEndpoints) > 2 {
		db.RedisSetEndpoint(config.Identifier, "rest", "1", restEndpoints[len(restEndpoints)-1].URL)
		db.RedisSetEndpoint(config.Identifier, "rest", "2", restEndpoints[len(restEndpoints)-2].URL)
		db.RedisSetEndpoint(config.Identifier, "rest", "3", restEndpoints[len(restEndpoints)-3].URL)
	}
	// store jrc result in redis
	if len(jrpcEndpoints) > 2 {
		db.RedisSetEndpoint(config.Identifier, "jrpc", "1", jrpcEndpoints[len(jrpcEndpoints)-1].URL)
		db.RedisSetEndpoint(config.Identifier, "jrpc", "2", jrpcEndpoints[len(jrpcEndpoints)-2].URL)
		db.RedisSetEndpoint(config.Identifier, "jrpc", "3", jrpcEndpoints[len(jrpcEndpoints)-3].URL)
	}
	// store web3 result in redis
	if len(web3Endpoints) > 2 {
		db.RedisSetEndpoint(config.Identifier, "web3", "1", web3Endpoints[len(web3Endpoints)-1].URL)
		db.RedisSetEndpoint(config.Identifier, "web3", "2", web3Endpoints[len(web3Endpoints)-2].URL)
		db.RedisSetEndpoint(config.Identifier, "web3", "3", web3Endpoints[len(web3Endpoints)-3].URL)
	}
}

func main() {
	for running {
		fmt.Println("Fetching network configs...")
		networkConfigs, err := resources.GetNetworkConfigs()
		if err != nil {
			// TODO: report to sentry?
			panic(err)
		}
		for _, v := range networkConfigs {
			go processNetwork(v)
			wg.Add(1)
		}
		wg.Wait()
	}
}
