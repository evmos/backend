package numia

type HeightResponse struct {
	LatestBlockHash   string `json:"latestBlockHash"`
	LatestBlockHeight string `json:"latestBlockHeight"`
	LatestBlockTime   string `json:"latestBlockTime"`
}

// QueryHeight queries the height of the latest block on the EVMOS blockchain.
// URL: "https://evmos.numia.xyz/height"
func (c *RPCClient) QueryHeight() (*HeightResponse, error) {

	// Unmarshal response into struct
	var data HeightResponse
	if err := c.get("/height", &data); err != nil {
		return nil, err
	}

	return &data, nil
}

type DelegationResponse struct {
	ValidatorAddress string `json:"validatorAddress"`
	Delegated        struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"delegated"`
	Unclaimed struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"unclaimed"`
}

// QueryDelegations queries the delegations of the requested address.
// It handles both Hex and Bech32 addresses.
// URL: "https://evmos.numia.xyz/evmos/delegations"
func (c *RPCClient) QueryDelegations(address string) ([]DelegationResponse, error) {
	var data []DelegationResponse
	if err := c.get("/evmos/delegations/"+address, &data); err != nil {
		return nil, err
	}

	return data, nil
}
