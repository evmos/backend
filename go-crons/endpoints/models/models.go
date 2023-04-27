// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package models

type Endpoint struct {
	URL     string  `json:"url"`
	Height  int     `json:"height"`
	Latency float64 `json:"latency"`
}

type RestResponse struct {
	Block struct {
		Header struct {
			Height string `json:"height"`
		} `json:"header"`
	} `json:"block"`
}

type JrpcTransactionErrorResponse struct {
	Error struct {
		Data string `json:"data"`
	} `json:"error"`
}

type JrpcStatusResponse struct {
	Result struct {
		SyncInfo struct {
			LatestBlockHeight string `json:"latest_block_height"`
		} `json:"sync_info"`
		NodeInfo struct {
			Other struct {
				TxIndex string `json:"tx_index"`
			} `json:"other"`
		} `json:"node_info"`
	} `json:"result"`
}

type Web3Response struct {
	Result string `json:"result"`
}
