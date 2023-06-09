// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tharsis/dashboard-backend/api"
	"github.com/tharsis/dashboard-backend/api/config"

	"github.com/tharsis/dashboard-backend/internal/v1/metrics"
)

func main() {
	// Flush metrics if we are killing the process
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		metrics.Flush()
		os.Exit(1)
	}()

	// Load the configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	rpcserver := api.NewServer(cfg)
	if err = rpcserver.Start(); err != nil {
		log.Printf("Error starting RPC server: %v\n", err)
	}
}
