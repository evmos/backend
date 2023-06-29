// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package metrics

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

func init() {
	err := sentry.Init(sentry.ClientOptions{
		// Sentry mainnet project endpoint
		// Dsn:              "https://536f4f89dcbb4a668f2c90466c0a26ca@o1345442.ingest.sentry.io/6660507",

		// Sentry gotest project endpoint
		Dsn:              "https://fc22a54305d3459aaeafbdc8c9b7f2a7@o1345442.ingest.sentry.io/4503964622454784",
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}

func Send(msg string) *sentry.EventID {
	return sentry.CaptureMessage(msg)
}

func Flush() {
	sentry.Flush(2 * time.Second)
}
