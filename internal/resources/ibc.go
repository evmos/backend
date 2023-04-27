// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package resources

import (
	"strings"

	"github.com/tharsis/dashboard-backend/internal/constants"
)

func GetIBCChannels() (map[string]map[string]string, error) {
	networkConfigs, err := GetNetworkConfigs()
	if err != nil {
		return nil, err
	}
	ibcChannels := make(map[string]map[string]string)
	destinationChannels := make(map[string]string)
	for _, networkConfig := range networkConfigs {
		for _, c := range networkConfig.Configurations {
			identifier := strings.ToUpper(c.Identifier)
			if c.ConfigurationType == constants.Mainnet && identifier != constants.EVMOS {
				destinationChannels[c.Source.DestinationChannel] = identifier
			}
			ibcChannel := make(map[string]string)
			ibcChannel[c.Source.SourceChannel] = constants.EVMOS
			ibcChannels[identifier] = ibcChannel
		}
	}
	ibcChannels["EVMOS"] = destinationChannels
	return ibcChannels, nil
}

func GetIBCCoins() (map[string]map[string]string, error) {
	erc20tokens, err := GetERC20Tokens()
	if err != nil {
		return nil, err
	}

	ibcCoins := make(map[string]map[string]string)

	for _, token := range erc20tokens {

		coinDenom := token.CoinDenom
		ibcCoin := make(map[string]string)
		ibcCoin["sourceDenom"] = token.Ibc.SourceDenom
		ibcCoin["source"] = strings.ToUpper(token.Ibc.Source)
		ibcCoins[coinDenom] = ibcCoin

	}

	return ibcCoins, nil
}
