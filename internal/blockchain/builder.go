// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package blockchain

import (
	"encoding/base64"
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DeriveCosmosAddress(pubkey string, prefix string) (string, error) {
	bz, err := base64.StdEncoding.DecodeString(pubkey)
	if err == nil {
		if len(bz) == secp256k1.PubKeySize {
			pk := &secp256k1.PubKey{Key: bz}
			return sdk.Bech32ifyAddressBytes(prefix, pk.Address().Bytes())
		}
	}

	return "", fmt.Errorf("invalid pubkey")
}
