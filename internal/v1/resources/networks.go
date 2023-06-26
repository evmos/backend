// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package resources

import (
	"strings"
)

func GetNetworks() ([]string, error) {
	networkConfigs, err := GetNetworkConfigs()
	if err != nil {
		return nil, err
	}
	networks := make([]string, 0, len(networkConfigs))
	for _, networkConfig := range networkConfigs {
		prefix := strings.ToUpper(networkConfig.Prefix)

		networks = append(networks, prefix)

	}

	return networks, nil
}
