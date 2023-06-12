// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package blockchain

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/evmos/evmos/v12/x/erc20/types"

	"github.com/tharsis/dashboard-backend/internal/db"
	"github.com/tharsis/dashboard-backend/internal/requester"
)

func GetERC20Balance(contract string, wallet string) (string, error) {
	cache, err := db.RedisGetERC20Balance(contract, wallet)
	if err == nil {
		return cache, nil
	}

	walletSplitted := strings.Split(wallet, "0x")
	if len(walletSplitted) != 2 {
		return "", fmt.Errorf("invalid wallet")
	}
	var sb strings.Builder
	sb.WriteString(`{"method":"eth_call", "params":[{"to": "`)
	sb.WriteString(contract)
	sb.WriteString(`", "data": "0x70a08231000000000000000000000000`)
	sb.WriteString(walletSplitted[1])
	sb.WriteString(`"}, "latest"], "id":1,"jsonrpc":"2.0"}`)
	jsonBody := []byte(sb.String())

	val, err := requester.MakePostRequest("EVMOS", "web3", "/", jsonBody)
	if err != nil {
		return "", err
	}

	var m map[string]interface{}
	err = json.Unmarshal([]byte(val), &m)
	if err != nil {
		return "", err
	}

	for k, v := range m {
		if k == "result" {
			m := new(big.Int)
			m.SetString(v.(string), 0)
			db.RedisSetERC20Balance(contract, wallet, m.String())
			return m.String(), nil
		}
	}

	return "", err
}

func CreateMsgConvertCoin(amount sdkmath.Int, token string, receiver string, sender string, prefix string) (sdk.Msg, error) {
	// to erc20
	from, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return &types.MsgConvertCoin{}, fmt.Errorf("error creating from address: %q", err)
	}

	msgConvert := types.NewMsgConvertCoin(sdk.Coin{Denom: token, Amount: amount}, common.HexToAddress(receiver), from)
	return msgConvert, nil
}

func CreateMsgConvertERC20(amount sdkmath.Int, receiver string, contract string, sender string, prefix string) (sdk.Msg, error) {
	// to ibc
	to, err := sdk.AccAddressFromBech32(receiver)
	if err != nil {
		return &types.MsgConvertERC20{}, fmt.Errorf("error creating to address: %q", err)
	}

	msgConvert := types.NewMsgConvertERC20(amount, to, common.HexToAddress(contract), common.HexToAddress(sender))
	return msgConvert, nil
}
