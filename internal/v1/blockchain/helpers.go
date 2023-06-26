// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package blockchain

import (
	"strings"
)

func IsEVMAddress(address string) bool {
	walletSplitted := strings.Split(address, "0x")
	if len(walletSplitted) != 2 {
		return false
	}
	return false
}
