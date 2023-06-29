// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/tharsis/dashboard-backend/internal/v1/constants"
)

type PubKeyAccount struct {
	Type string `json:"@type"`
	Key  string `json:"string"`
}

type BaseAccount struct {
	Address       string        `json:"address"`
	PubKey        PubKeyAccount `json:"pub_key"`
	AccountNumber string        `json:"account_number"`
	Sequence      string        `json:"sequence"`
}

type BaseVestingAccount struct {
	BaseAccount BaseAccount `json:"base_account"`
}
type BaseAccountDetails struct {
	Type               string             `json:"@type"`
	BaseAccount        BaseAccount        `json:"base_account"`
	CodeHash           string             `json:"code_hash"`
	BaseVestingAccount BaseVestingAccount `json:"base_vesting_account"`
}

type BaseAccountResponse struct {
	Account BaseAccountDetails `json:"account"`
}

type AccountDetails struct {
	Type               string             `json:"@type"`
	Address            string             `json:"address"`
	PubKey             PubKeyAccount      `json:"pub_key"`
	AccountNumber      string             `json:"account_number"`
	Sequence           string             `json:"sequence"`
	BaseVestingAccount BaseVestingAccount `json:"base_vesting_account"`
}

type AccountResponse struct {
	Account AccountDetails `json:"account"`
}

func GetAccountInfo(sender string, srcChain string) (uint64, uint64, error) {
	// EVMOS and OSMO have different struct for AccountInternal
	val, err := AccountInternal(sender, srcChain)
	if err != nil {
		return 0, 0, fmt.Errorf("error while getting account details, please try again")
	}

	if srcChain == constants.EVMOS {
		var accountDetails BaseAccountResponse
		err = json.Unmarshal([]byte(val), &accountDetails)
		if err != nil {
			return 0, 0, fmt.Errorf("error while parsing account details, please try again")
		}

		account := accountDetails.Account.BaseAccount
		if strings.Contains(strings.ToLower(accountDetails.Account.Type), "vesting") {
			account = accountDetails.Account.BaseVestingAccount.BaseAccount
		}

		accountNumber, err := strconv.ParseUint(account.AccountNumber, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("error while getting account number, please try again")
		}
		sequence, err := strconv.ParseUint(account.Sequence, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("error while getting account sequence, please try again")
		}
		return accountNumber, sequence, nil
	}
	var accountDetails AccountResponse
	err = json.Unmarshal([]byte(val), &accountDetails)
	if err != nil {
		return 0, 0, fmt.Errorf("error while parsing account details, please try again")
	}

	accountNumber := accountDetails.Account.AccountNumber
	accountSequence := accountDetails.Account.Sequence

	if strings.Contains(strings.ToLower(accountDetails.Account.Type), "vesting") {
		accountNumber = accountDetails.Account.BaseVestingAccount.BaseAccount.AccountNumber
		accountSequence = accountDetails.Account.BaseVestingAccount.BaseAccount.Sequence
	}

	number, err := strconv.ParseUint(accountNumber, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error while getting account number, please try again")
	}

	sequence, err := strconv.ParseUint(accountSequence, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error while getting account sequence, please try again")
	}
	return number, sequence, nil
}

func GetHeightInfo(m MessageSendIBCStruct) (uint64, uint64, error) {
	h, r, err := ChainHeightInternal(m.Message.DstChain)
	if err != nil {
		// y que no pueda enviarse numero negativo
		return 0, 0, fmt.Errorf("error while getting height chain info, please try again")
	}
	height, err := strconv.ParseUint(h, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error while getting height, please try again")
	}

	revision, err := strconv.ParseUint(r, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error while getting revision, please try again")
	}

	return height, revision, nil
}

func GetDenom(token string, srcChain string) (string, error) {
	if token == "EVMOS" { //nolint:all
		if srcChain == "EVMOS" {
			return "aevmos", nil
		}
		network, err := NetworkConfigByNameInternal(srcChain)
		if err != nil {
			return "", fmt.Errorf("invalid source chain, please try again")
		}

		var config NetworkByName
		err = json.Unmarshal([]byte(network), &config)
		if err != nil {
			return "", fmt.Errorf("invalid params for network configuration, please try again")
		}

		for _, v := range config.Values.Configurations {
			if v.ConfigurationType == constants.Mainnet {
				return v.Source.SourceIBCDenomToEvmos, nil
			}
		}

	} else if srcChain == "EVMOS" {
		token, err := ERC20TokensByNameInternal(token)
		if err != nil {
			return "", err
		}

		var tokensByName TokensByName
		err = json.Unmarshal([]byte(token), &tokensByName)
		if err != nil {
			return "", fmt.Errorf("error parsing token, please try again")
		}
		return tokensByName.Values.CosmosDenom, nil

	} else {
		token, err := ERC20TokensByNameInternal(token)
		if err != nil {
			return "", err
		}
		var tokensByName TokensByName
		err = json.Unmarshal([]byte(token), &tokensByName)
		if err != nil {
			return "", fmt.Errorf("error parsing token, please try again")
		}
		return tokensByName.Values.Ibc.SourceDenom, nil
	}
	return "", fmt.Errorf("invalid denom, please try again")
}

func GetConfigInfo(m MessageSendIBCStruct) (string, string, string, string, string, error) {
	channel := ""
	clientID := ""
	chainID := ""
	prefix := ""
	explorerTxURL := ""

	networkSrcChain, err := NetworkConfigByNameInternal(m.Message.SrcChain)
	if err != nil {
		return "", "", "", "", "", err
	}

	var configSrcChain NetworkByName
	err = json.Unmarshal([]byte(networkSrcChain), &configSrcChain)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error parsing network, please try again")
	}

	prefix = configSrcChain.Values.Prefix

	for _, v := range configSrcChain.Values.Configurations {
		if v.ConfigurationType == constants.Mainnet {
			chainID = v.ChainID
			channel = v.Source.SourceChannel
			clientID = v.ClientID
			explorerTxURL = v.ExplorerTxURL
		}
	}

	if m.Message.SrcChain == "EVMOS" {
		// TODO: remove this after https://github.com/evmos/chain-token-registry/pull/29 is merged
		if m.Message.DstChain == "COSMOS" {
			m.Message.DstChain = "ATOM"
		}
		if m.Message.DstChain == "STARS" {
			m.Message.DstChain = "STARGAZE"
		}
		network, err := NetworkConfigByNameInternal(m.Message.DstChain)
		if err != nil {
			return "", "", "", "", "", err
		}

		var config NetworkByName
		err = json.Unmarshal([]byte(network), &config)
		if err != nil {
			return "", "", "", "", "", fmt.Errorf("error parsing network, please try again")
		}

		for _, v := range config.Values.Configurations {
			if v.ConfigurationType == constants.Mainnet {
				channel = v.Source.DestinationChannel
				clientID = v.ClientID
			}
		}

		if clientID == "" {
			return "", "", "", "", "", fmt.Errorf("client Id not registered")
		}

	}

	if channel == "" {
		return "", "", "", "", "", fmt.Errorf("invalid chain-channel combination")
	}

	if chainID == "" {
		return "", "", "", "", "", fmt.Errorf("chain Id not registered")
	}

	if prefix == "" {
		return "", "", "", "", "", fmt.Errorf("prefix not registered")
	}
	return channel, clientID, chainID, prefix, explorerTxURL, nil
}

type ClientStatus struct {
	Status string `json:"status"`
}

func IsIBCChannelActive(chain string, clientID string) error {
	val, err := IBCClientStatusInternal(chain, clientID)
	if err != nil {
		return err
	}
	var m ClientStatus
	_ = json.Unmarshal([]byte(val), &m)
	status := m.Status
	if status != "Active" {
		return fmt.Errorf(status)
	}

	return nil
}

func GetERC20Address(token string) (string, error) {
	token, err := ERC20TokensByNameInternal(token)
	if err != nil {
		return "", err
	}
	var tokensByName TokensByName
	err = json.Unmarshal([]byte(token), &tokensByName)
	if err != nil {
		return "", fmt.Errorf("error parsing token, please try again")
	}
	return tokensByName.Values.ERC20Address, nil
}

func GetSourceInfo(srcChain string) (string, string, string, error) {
	prefix := ""
	chainID := ""
	explorerTxURL := ""
	networkSrcChain, err := NetworkConfigByNameInternal(srcChain)
	if err != nil {
		return "", "", "", err
	}

	var configSrcChain NetworkByName
	err = json.Unmarshal([]byte(networkSrcChain), &configSrcChain)
	if err != nil {
		return "", "", "", fmt.Errorf("error parsing network, please try again")
	}

	prefix = configSrcChain.Values.Prefix

	for _, v := range configSrcChain.Values.Configurations {
		if v.ConfigurationType == constants.Mainnet {
			chainID = v.ChainID
			explorerTxURL = v.ExplorerTxURL
		}
	}

	if chainID == "" {
		return "", "", "", fmt.Errorf("chain Id not registered")
	}

	return prefix, chainID, explorerTxURL, nil
}
