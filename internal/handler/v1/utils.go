// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/fasthttp/router"
	"github.com/tharsis/dashboard-backend/internal/db"
	"github.com/tharsis/dashboard-backend/internal/resources"
	"github.com/valyala/fasthttp"
)

type ChainHeightParams struct {
	Height   uint64 `json:"height"`
	Revision uint64 `jon:"revision"`
}

func ChainHeightInternal(chain string) (string, string, error) {
	cache, err := db.RedisGetChainHeight(chain)
	if err == nil {
		m := ChainHeightParams{}
		if err := json.Unmarshal([]byte(cache), &m); err != nil {
			return "", "", err
		}
		return strconv.FormatUint(m.Height, 10), strconv.FormatUint(m.Revision, 10), err
	}

	val, err := GetRequestJrpc(chain, "/status")
	if err != nil {
		return val, val, err
	}

	var m map[string]interface{}
	err = json.Unmarshal([]byte(val), &m)
	if err != nil {
		return val, val, err
	}

	var height string
	var revision string

	for k, v := range m {
		if k == "result" {
			for resK, resV := range v.(map[string]interface{}) {
				if resK == "sync_info" {
					for syncK, syncV := range resV.(map[string]interface{}) {
						if syncK == "latest_block_height" {
							height = syncV.(string)
						}
					}
				}

				if resK == "node_info" {
					for infoK, infoV := range resV.(map[string]interface{}) {
						if infoK == "network" {
							network := infoV.(string)
							revisionArray := strings.Split(network, "-")
							revision = revisionArray[len(revisionArray)-1]
						}
					}
				}
			}
			// Store the cache
			val = `{"height":` + height + `,"revision":` + revision + `}`
			db.RedisSetChainHeight(chain, val)
			return height, revision, nil
		}
	}

	return "", "", fmt.Errorf("invalid status response")
}

type TxStatusType string

const (
	TSUnconfirmed TxStatusType = "unconfirmed"
	TSConfirmed   TxStatusType = "confirmed"
	TSError       TxStatusType = "error"
)

func getTransaction(txHash string, chain string) (TxStatusType, map[string]interface{}) {
	endpoint := BuildTwoParamEndpoint("/tx?hash=0x", txHash)
	val, err := GetRequestJrpc(chain, endpoint)
	if err != nil {
		return TSError, nil
	}

	var txHashResult map[string]interface{}
	err = json.Unmarshal([]byte(val), &txHashResult)
	if err != nil {
		return TSError, nil
	}

	if val, ok := txHashResult["error"]; ok {
		if errorCode, ok := val.(map[string]interface{}); ok {
			if val2, ok2 := errorCode["code"]; ok2 {
				if code, ok3 := val2.(float64); ok3 {
					if code == -32603 {
						return TSUnconfirmed, txHashResult
					}
				}
			}
		}
	}

	if val, ok := txHashResult["result"]; ok {
		if result, ok := val.(map[string]interface{}); ok {
			if val2, ok2 := result["tx_result"]; ok2 {
				if txResult, ok := val2.(map[string]interface{}); ok {
					if val3, ok3 := txResult["code"]; ok3 {
						if code, ok4 := val3.(float64); ok4 {
							if code == 0 {
								return TSConfirmed, txHashResult
							}
						}
					}
				}
			}
		}
	}

	return TSError, txHashResult
}

func getIBCData(txJSON map[string]interface{}) (string, string, string, error) {
	var rawLog string
	found := false

	if val, ok := txJSON["result"]; ok {
		if result, ok := val.(map[string]interface{}); ok {
			if val2, ok2 := result["tx_result"]; ok2 {
				if txResult, ok := val2.(map[string]interface{}); ok {
					if val3, ok3 := txResult["log"]; ok3 {
						if log, ok4 := val3.(string); ok4 {
							rawLog = log
							found = true
						}
					}
				}
			}
		}
	}

	if !found {
		return "", "", "", fmt.Errorf("error getting logs")
	}

	var jsonLogs []interface{}
	_ = json.Unmarshal([]byte(rawLog), &jsonLogs)

	if len(jsonLogs) == 0 {
		return "", "", "", fmt.Errorf("error getting logs")
	}

	sequence := ""
	srcChannel := ""
	dstChannel := ""

	// TODO: make an object to unmarshal this without so many ifs
	if firstLogs, ok := jsonLogs[0].(map[string]interface{}); ok {
		if eventsRaw, ok := firstLogs["events"]; ok {
			if events, ok := eventsRaw.([]interface{}); ok {
				for _, e := range events {
					if eventRaw, ok := e.(map[string]interface{}); ok {
						if eventTypeRaw, ok := eventRaw["type"]; ok {
							if eventType, ok := eventTypeRaw.(string); ok {
								if eventType == "send_packet" {
									if attributesRaw, ok := eventRaw["attributes"]; ok {
										if attributes, ok := attributesRaw.([]interface{}); ok {
											for _, attributeRaw := range attributes {
												if attribute, ok := attributeRaw.(map[string]interface{}); ok {
													// Get the key
													if attributeKeyRaw, ok := attribute["key"]; ok {
														if attributeKey, ok := attributeKeyRaw.(string); ok {
															// Get the value
															if attributeValueRaw, ok := attribute["value"]; ok {
																if attributeValue, ok := attributeValueRaw.(string); ok {
																	if attributeKey == "packet_sequence" {
																		sequence = attributeValue
																	}
																	if attributeKey == "packet_src_channel" {
																		srcChannel = attributeValue
																	}
																	if attributeKey == "packet_dst_channel" {
																		dstChannel = attributeValue
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	err := fmt.Errorf("error getting logs")
	if sequence != "" {
		err = nil
	}

	return sequence, srcChannel, dstChannel, err
}

func isIBCExecutedResponse(executed bool, msg string) string {
	var sb strings.Builder
	sb.WriteString(`{"executed":`)
	if executed {
		sb.WriteString(`true`)
	} else {
		sb.WriteString(`false`)
	}
	sb.WriteString(`,"msg":"`)
	sb.WriteString(msg)
	sb.WriteString(`"}`)
	return sb.String()
}

func isIBCExecuted(ctx *fasthttp.RequestCtx) {
	chain := getChain(ctx)
	status, txJSON := getTransaction(paramToString("tx_hash", ctx), chain)
	if status != TSConfirmed {
		// Transaction not confirmed
		sendResponse(isIBCExecutedResponse(false, "Transaction not confirmed"), nil, ctx)
		return
	}

	sequence, srcChannel, dstChannel, err := getIBCData(txJSON)
	if err != nil {
		// Error reading ibc info
		sendResponse(isIBCExecutedResponse(false, "Error getting IBC info"), nil, ctx)
		return
	}
	ibcChannels, err := resources.GetIBCChannels()
	if err != nil {
		return
	}
	val, ok := ibcChannels[chain]
	if !ok {
		sendResponse(isIBCExecutedResponse(false, "Invalid chain-channel combination"), nil, ctx)
		return
	}
	dstChain, ok := val[srcChannel]

	if !ok {
		sendResponse(isIBCExecutedResponse(false, "Invalid chain-channel combination"), nil, ctx)
		return
	}

	endpoint := BuildFourParamEndpoint("/ibc/core/channel/v1/channels/", dstChannel, "/ports/transfer/packet_acks/", sequence)
	dstChainRequest, err := getRequestRest(dstChain, endpoint)
	if err != nil {
		sendResponse(isIBCExecutedResponse(false, "ACK not found"), nil, ctx)
		return
	}

	var ackResponse map[string]interface{}
	_ = json.Unmarshal([]byte(dstChainRequest), &ackResponse)

	if _, ok := ackResponse["code"]; ok {
		sendResponse(isIBCExecutedResponse(false, "IBC not executed yet"), nil, ctx)
		return
	}

	sendResponse(isIBCExecutedResponse(true, "IBC ack ready"), nil, ctx)
}
