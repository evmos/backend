// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package blockchain

import (
	"testing"
)

func TestDeriveCosmosAddress(t *testing.T) {
	pubkey := "AyOw8+o7OXFMu/UccBbUfTB6yQ0jy6EPEF3UEiRvmyNJ"
	walletOsmosis := "osmo1pmk2r32ssqwps42y3c9d4clqlca403yd05x9ye"
	walletAkash := "akash1pmk2r32ssqwps42y3c9d4clqlca403yd25cjt3"

	derivedOsmosis, err := DeriveCosmosAddress(pubkey, "osmo")
	if err != nil {
		t.Fatalf("Error getting osmosis wallet, %s", err)
	}

	if derivedOsmosis != walletOsmosis {
		t.Fatalf("The generated wallet (%s) is not equal to (%s)", derivedOsmosis, walletOsmosis)
	}

	derivedAkash, err := DeriveCosmosAddress(pubkey, "akash")
	if err != nil {
		t.Fatalf("Error getting akash wallet, %s", err)
	}

	if derivedAkash != walletAkash {
		t.Fatalf("The generated wallet (%s) is not equal to (%s)", derivedAkash, walletAkash)
	}

	errorString, err := DeriveCosmosAddress("", "osmo")
	if errorString != "" {
		t.Fatalf("Error addresses must be empty")
	}

	if err == nil {
		t.Fatalf("Derive empty pubkey key must return error")
	}
}
