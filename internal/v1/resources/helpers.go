// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package resources

import (
	"encoding/json"

	"github.com/tharsis/dashboard-backend/internal/v1/constants"
	"github.com/tharsis/dashboard-backend/internal/v1/db"
	"github.com/tharsis/dashboard-backend/internal/v1/requester"
)

func GetERC20Tokens() ([]CoinConfig, error) {
	var erc20tokens []CoinConfig
	if redisVal, err := db.RedisGetERC20TokensDirectory(); err == nil && redisVal != "null" {
		err = json.Unmarshal([]byte(redisVal), &erc20tokens)
		if err != nil {
			return nil, err
		}
	} else {
		gitRes, err := requester.GetERC20TokensDirectory()
		if err != nil {
			return nil, err
		}

		for _, v := range gitRes {
			var coinConfig CoinConfig
			err := json.Unmarshal([]byte(v.Content), &coinConfig)
			if err != nil {
				return nil, err
			}
			erc20tokens = append(erc20tokens, coinConfig)
		}

		stringRes, err := json.Marshal(erc20tokens)
		if err != nil {
			return nil, err
		}

		db.RedisSetERC20TokensDirectory(string(stringRes))
	}

	return erc20tokens, nil
}

func GetNetworkConfigs() ([]NetworkConfig, error) {
	var networkConfigs []NetworkConfig
	if redisVal, err := db.RedisGetNetworkConfig(); err == nil && redisVal != "null" {
		err = json.Unmarshal([]byte(redisVal), &networkConfigs)
		if err != nil {
			return nil, err
		}
	} else {
		gitRes, err := requester.GetNetworkConfig()
		if err != nil {
			return nil, err
		}

		for _, v := range gitRes {
			var networkConfig NetworkConfig
			err := json.Unmarshal([]byte(v.Content), &networkConfig)
			if err != nil {
				return nil, err
			}
			networkConfigs = append(networkConfigs, networkConfig)
		}

		stringRes, err := json.Marshal(networkConfigs)
		if err != nil {
			return nil, err
		}

		db.RedisSetNetworkConfig(string(stringRes))
	}

	return networkConfigs, nil
}

func GetMainnetConfig(nc NetworkConfig) ConfigurationEntry {
	for _, nc := range nc.Configurations {
		if nc.ConfigurationType == constants.Mainnet {
			return nc
		}
	}
	return ConfigurationEntry{}
}
