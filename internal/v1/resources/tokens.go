// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package resources

import (
	"strconv"
)

func GetERC20ModuleCoins() (map[string]ERC20ModuleCoin, error) {
	erc20Tokens, err := GetERC20Tokens()
	if err != nil {
		return nil, err
	}

	erc20Coins := make(map[string]ERC20ModuleCoin)

	for _, token := range erc20Tokens {

		decimals, err := strconv.Atoi(token.Exponent)
		if err != nil {
			decimals = 18
		}
		erc20Token := ERC20ModuleCoin{
			Denom:               token.CosmosDenom,
			Erc20:               token.Erc20Address,
			TokenName:           token.CoinDenom,
			Description:         token.Description,
			CoingeckoID:         token.CoingeckoID,
			Name:                token.Name,
			TokenRepresentation: token.TokenRepresentation,
			Symbol:              token.CoinDenom,
			Decimals:            decimals,
			ChainPrefix:         token.CoinSourcePrefix,
			HandledByExternalUI: token.HandledByExternalUI,
			PngSrc:              token.PngSrc,
			Prefix:              token.CoinSourcePrefix,
		}
		erc20Coins[token.CoinDenom] = erc20Token
	}

	return erc20Coins, nil
}
