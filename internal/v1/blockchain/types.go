// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package blockchain

type EipToSignTypesStruct struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type EipToSignTypes struct {
	Coin         []EipToSignTypesStruct `json:"Coin"`
	EIP712Domain []EipToSignTypesStruct `json:"EIP712Domain"`
	Fee          []EipToSignTypesStruct `json:"Fee"`
	Msg          []EipToSignTypesStruct `json:"Msg"`
	MsgValue     []EipToSignTypesStruct `json:"MsgValue"`
	Tx           []EipToSignTypesStruct `json:"Tx"`
	TypeAmount   []EipToSignTypesStruct `json:"TypeAmount"`
}

type EipToSignDomain struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	ChainID           string `json:"chainId"`
	VerifyingContract string `json:"verifyingContract"`
	Salt              string `json:"salt"`
}

type EipToSignAmountStruct struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

type EipToSignFee struct {
	Amount   []EipToSignAmountStruct `json:"amount"`
	FeePayer string                  `json:"feePayer"`
	Gas      string                  `json:"gas"`
}

type EipToSignMessage struct {
	AccountNumber string              `json:"account_number"`
	ChainID       string              `json:"chain_id"`
	Fee           EipToSignFee        `json:"fee"`
	Memo          string              `json:"memo"`
	Msgs          []EipToSignMessages `json:"msgs"`
	Sequence      string              `json:"sequence"`
}

type EipToSignMsgSend struct {
	Types       EipToSignTypes   `json:"types"`
	PrimaryType string           `json:"primaryType"`
	Domain      EipToSignDomain  `json:"domain"`
	Message     EipToSignMessage `json:"message"`
}

type EipToSignMessagesValue struct {
	Amount      []EipToSignAmountStruct `json:"amount"`
	FromAddress string                  `json:"from_address"`
	ToAddress   string                  `json:"to_address"`
}

type EipToSignMessages struct {
	Type  string                 `json:"type"`
	Value EipToSignMessagesValue `json:"value"`
}

type EipToSignIBC struct {
	Types       EipToSignTypes      `json:"types"`
	PrimaryType string              `json:"primaryType"`
	Domain      EipToSignDomain     `json:"domain"`
	Message     EipToSignMessageIBC `json:"message"`
}

type EipToSignMessageIBC struct {
	AccountNumber string                 `json:"account_number"`
	ChainID       string                 `json:"chain_id"`
	Fee           EipToSignFee           `json:"fee"`
	Memo          string                 `json:"memo"`
	Msgs          []EipToSignMessagesIBC `json:"msgs"`
	Sequence      string                 `json:"sequence"`
}

type EipToSignMessagesValueIBC struct {
	Receiver         string                `json:"receiver"`
	Sender           string                `json:"sender"`
	SrcChannel       string                `json:"source_channel"`
	SrcPort          string                `json:"source_port"`
	TimeoutHeight    TimeoutHeightIBC      `json:"timeout_height"`
	TimeoutTimestamp string                `json:"timeout_timestamp"`
	Token            EipToSignAmountStruct `json:"token"`
	Sequence         string                `json:"string"`
}

type TimeoutHeightIBC struct {
	RevisionHeight string `json:"revision_height"`
	RevisionNumber string `json:"revision_number"`
}

type EipToSignMessagesIBC struct {
	Type  string                    `json:"type"`
	Value EipToSignMessagesValueIBC `json:"value"`
}

type EipToSignMessagesValueConvertERC20 struct {
	Amount          string `json:"amount"`
	FromAddress     string `json:"from_address"`
	ToAddress       string `json:"to_address"`
	ContractAddress string `json:"contract_address"`
}

type EipToSignMessagesConvertERC20 struct {
	Type  string                             `json:"type"`
	Value EipToSignMessagesValueConvertERC20 `json:"value"`
}

type EipToSignMessageConvertERC20 struct {
	AccountNumber string                          `json:"account_number"`
	ChainID       string                          `json:"chain_id"`
	Fee           EipToSignFee                    `json:"fee"`
	Memo          string                          `json:"memo"`
	Msgs          []EipToSignMessagesConvertERC20 `json:"msgs"`
	Sequence      string                          `json:"sequence"`
}

type EipToSignMsgConvertERC20 struct {
	Types       EipToSignTypes               `json:"types"`
	PrimaryType string                       `json:"primaryType"`
	Domain      EipToSignDomain              `json:"domain"`
	Message     EipToSignMessageConvertERC20 `json:"message"`
}

type EipToSignMessagesValueDelegate struct {
	Amount           EipToSignAmountStruct `json:"amount"`
	DelegatorAddress string                `json:"delegator_address"`
	ValidatorAddress string                `json:"validator_address"`
}

type EipToSignMsgDelegate struct {
	Types       EipToSignTypes               `json:"types"`
	PrimaryType string                       `json:"primaryType"`
	Domain      EipToSignDomain              `json:"domain"`
	Message     EipToSignMessageConvertERC20 `json:"message"`
}

type EipToSignMessageDelegate struct {
	AccountNumber string                           `json:"account_number"`
	ChainID       string                           `json:"chain_id"`
	Fee           EipToSignFee                     `json:"fee"`
	Memo          string                           `json:"memo"`
	Msgs          []EipToSignMessagesValueDelegate `json:"msgs"`
	Sequence      string                           `json:"sequence"`
}
