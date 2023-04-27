// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tharsis/dashboard-backend/internal/endpoints"
	"github.com/tharsis/dashboard-backend/internal/metrics"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
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

	r := endpoints.CreateRouter()

	ln, err := reuseport.Listen("tcp4", "0.0.0.0:8081")
	if err != nil {
		fmt.Printf("error in reuseport listener: %v\n", err)
		return
	}

	if err = fasthttp.Serve(ln, r.Handler); err != nil {
		fmt.Printf("error in fasthttp Server: %v\n", err)
	}
}
