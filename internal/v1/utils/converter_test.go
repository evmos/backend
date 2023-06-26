// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package utils_test

import (
	"fmt"
	"testing"

	decimal "github.com/cosmos/cosmos-sdk/types"
	"github.com/tharsis/dashboard-backend/internal/v1/utils"
)

func TestNumberToBiggerDenom(t *testing.T) {
	// Invalid string
	value := "evmos"
	exponent := uint64(18)

	_, err := utils.NumberToBiggerDenom(value, exponent)

	if err == nil {
		t.Fatalf("Error converting invalid string")
	}

	// EVMOS
	value = "1000000000000000000"
	exponent = uint64(18)
	expected := decimal.NewDec(1)

	res, err := utils.NumberToBiggerDenom(value, exponent)

	if err != nil || !res.Equal(expected) {
		t.Fatalf("Error converting evmos %v", res)
	}
	// Osmosis
	value = "1000000"
	exponent = uint64(6)
	expected = decimal.NewDec(1)

	res, err = utils.NumberToBiggerDenom(value, exponent)

	if err != nil || !res.Equal(expected) {
		t.Fatalf("Error converting osmosis %v", res)
	}
}

func TestNumberToLowerDenom(t *testing.T) {
	// Invalid string
	value := "evmos"
	exponent := uint64(18)

	_, err := utils.NumberToLowerDenom(value, exponent)

	if err == nil {
		t.Fatalf("Error converting invalid string")
	}

	// EVMOS
	value = "1.23"
	exponent = uint64(18)
	expected := decimal.NewDec(1230000000000000000)

	res, err := utils.NumberToLowerDenom(value, exponent)
	if err != nil || !res.Equal(expected) {
		fmt.Println(res)
		t.Fatalf("Error converting evmos %v", res)
	}

	// Osmosis
	value = "1.23"
	exponent = uint64(6)
	expected = decimal.NewDec(1230000)

	res, err = utils.NumberToLowerDenom(value, exponent)
	if err != nil || !res.Equal(expected) {
		fmt.Println(res)
		t.Fatalf("Error converting osmosis %v", res)
	}
}
