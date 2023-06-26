// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package utils

import (
	"fmt"

	decimal "github.com/cosmos/cosmos-sdk/types"
)

func NumberToBiggerDenom(value string, exponent uint64) (decimal.Dec, error) {
	dec, err := decimal.NewDecFromStr(value)
	if err != nil {
		return decimal.Dec{}, fmt.Errorf("invalid string")
	}

	exp := decimal.NewDec(10)
	exp = exp.Power(exponent)
	dec = dec.Quo(exp)

	return dec, nil
}

func NumberToLowerDenom(value string, exponent uint64) (decimal.Dec, error) {
	dec, err := decimal.NewDecFromStr(value)
	if err != nil {
		return decimal.Dec{}, fmt.Errorf("invalid string")
	}

	exp := decimal.NewDec(10)
	exp = exp.Power(exponent)
	dec = dec.Mul(exp)

	return dec, nil
}
