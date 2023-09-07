// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package v1

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/gogo/protobuf/proto"

	"github.com/tharsis/dashboard-backend/internal/v1/blockchain"
	"github.com/tharsis/dashboard-backend/internal/v1/constants"
	"github.com/valyala/fasthttp"
)

type MessageSendStruct struct {
	Transaction blockchain.TransactionParams `json:"transaction"`
	Message     blockchain.MessageSendParams `json:"message"`
}

type SignDataString struct {
	Body      string `json:"body"`
	AuthInfo  string `json:"authInfo"`
	SignBytes string `json:"signBytes"`
}

type TransactionString struct {
	LegacyAmino      SignDataString `json:"legacyAmino"`
	SignDirect       SignDataString `json:"signDirect"`
	EipToSign        string         `json:"eipToSign"`
	AccountNumber    string         `json:"accountNumber"`
	ChainID          string         `json:"chainId"`
	ExplorerTxURL    string         `json:"explorerTxUrl"`
	DataSigningAmino string         `json:"dataSigningAmino"`
}

type BroadcastMetamaskParams struct {
	Chain       uint64 `json:"chainId"`
	FeePayer    string `json:"feePayer"`
	FeePayerSig string `json:"feePayerSig"`
	Body        string `json:"body"`
	AuthInfo    string `json:"authInfo"`
}

func BroadcastMetamask(ctx *fasthttp.RequestCtx) {
	m := BroadcastMetamaskParams{}
	if err := json.Unmarshal(ctx.PostBody(), &m); err != nil {
		sendResponse(buildErrorBroadcast("Error while parsing broadcast, please try again"), nil, ctx)
		return
	}

	var bodyProto tx.TxBody
	bodyBytes, err := base64.StdEncoding.DecodeString(m.Body)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	err = bodyProto.Unmarshal(bodyBytes)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	var authInfoProto tx.AuthInfo
	authInfoBytes, err := base64.StdEncoding.DecodeString(m.AuthInfo)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	err = authInfoProto.Unmarshal(authInfoBytes)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	txRaw, err := blockchain.JoinEipWithSignature(m.Chain, m.FeePayer, m.FeePayerSig, &bodyProto, authInfoProto)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	bytesTxRaw, err := proto.Marshal(&txRaw)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	val, err := broadcastInternal(bytesTxRaw, "EVMOS")
	if err != nil {
		sendResponse("", err, ctx)
		return
	}
	sendResponse(val, err, ctx)
}

func GetTransactionBytes(tx blockchain.Transaction, accountNumber uint64, chainID string, explorerTxURL string) ([]byte, error) {
	bodyBytesSignDirect, err := proto.Marshal(&tx.SignDirect.Body)
	if err != nil {
		return []byte{}, err
	}
	bodySignDirect := base64.StdEncoding.EncodeToString(bodyBytesSignDirect)

	authInfoBytesSignDirect, err := proto.Marshal(&tx.SignDirect.AuthInfo)
	if err != nil {
		return []byte{}, err
	}
	authInfoSignDirect := base64.StdEncoding.EncodeToString(authInfoBytesSignDirect)

	bodyBytesLegacyAmino, err := proto.Marshal(&tx.LegacyAmino.Body)
	if err != nil {
		return []byte{}, err
	}
	bodyLegacyAmino := base64.StdEncoding.EncodeToString(bodyBytesLegacyAmino)

	authInfoBytesLegacyAmino, err := proto.Marshal(&tx.LegacyAmino.AuthInfo)
	if err != nil {
		return []byte{}, err
	}
	authInfoLegacyAmino := base64.StdEncoding.EncodeToString(authInfoBytesLegacyAmino)

	eipToSign := ""

	if tx.EipToSign != "" {
		err = json.Unmarshal([]byte(tx.EipToSign), &tx.MessagingEncoding)
		if err != nil {
			return []byte{}, err
		}
		resEipToSignBytes, err := json.Marshal(tx.MessagingEncoding)
		if err != nil {
			return []byte{}, err
		}
		eipToSign = base64.StdEncoding.EncodeToString(resEipToSignBytes)
	}

	account := strconv.FormatUint(accountNumber, 10)

	txString := TransactionString{
		LegacyAmino: SignDataString{
			Body:      bodyLegacyAmino,
			AuthInfo:  authInfoLegacyAmino,
			SignBytes: tx.LegacyAmino.SignBytes,
		},
		SignDirect: SignDataString{
			Body:      bodySignDirect,
			AuthInfo:  authInfoSignDirect,
			SignBytes: tx.SignDirect.SignBytes,
		},
		EipToSign:        eipToSign,
		AccountNumber:    account,
		ChainID:          chainID,
		ExplorerTxURL:    explorerTxURL,
		DataSigningAmino: tx.DataSigningAmino,
	}

	return json.Marshal(txString)
}

type IBCTransferParams struct {
	SrcChain      string `json:"srcChain"`
	DstChain      string `json:"dstChain"`
	Sender        string `json:"sender"`
	Receiver      string `json:"receiver"`
	Amount        string `json:"amount"`
	Token         string `json:"token"`
	UseERC20Denom bool   `json:"useERC20Denom"`
}

type MessageSendIBCStruct struct {
	Transaction blockchain.TransactionIBCParams `json:"transaction"`
	Message     IBCTransferParams               `json:"message"`
}

type IBCChannelActive struct {
	Status string `json:"status"`
}

type ParamsFeeUnSignMarket struct {
	MinGasPrice string `json:"min_gas_price"`
}

type FeeUnsignMarketStruct struct {
	Params ParamsFeeUnSignMarket `json:"params"`
}

type GasPriceStruct struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Result  string `json:"result"`
}

const (
	DefaultGas float64 = 350000
	// osmosis had an upgrade V15, and it seems they got rid of 0 fees and put in minimum fees
	// such minimum fees are requiring more gas, from the expected 200000 to 250000
	DefaultGasIbcTransfer          float64 = 300000
	feeDenom                       string  = "aevmos"
	EvmosTxFeeConvertAssetGas      float64 = 10500000
	DefaultGasIbcTransferUint      uint64  = 200000
	MarginMultiplicatorCoefficient float64 = 1.15
)

func GenerateFeeGasPrice(gas float64) (int64, error) {
	// this function uses the eth_gasPrice to solve the fee issues that we are facing
	val, err := EthGasPriceInternal()
	if err != nil {
		return 0, err
	}
	var m GasPriceStruct
	_ = json.Unmarshal([]byte(val), &m)

	hexInt, err := strconv.ParseInt(strings.TrimPrefix(m.Result, "0x"), 16, 64)
	if err != nil {
		return 0, err
	}
	floatValue := float64(hexInt)

	// TODO: convert everything to sdkInt
	return int64(floatValue * gas * MarginMultiplicatorCoefficient), nil
}

func GenerateFeeUnSignMarket(chain string, gas float64) (int64, error) {
	val, err := FeeMarketParamsInternal(chain)
	if err != nil {
		return 0, err
	}

	var m FeeUnsignMarketStruct
	_ = json.Unmarshal([]byte(val), &m)
	minGasPrice, err := strconv.ParseFloat(m.Params.MinGasPrice, 64)
	if err != nil {
		return 0, err
	}
	// TODO: convert everything to sdkInt
	return int64(minGasPrice * gas * MarginMultiplicatorCoefficient), nil
}

func IBCTransfer(ctx *fasthttp.RequestCtx) {
	m := MessageSendIBCStruct{}
	if err := json.Unmarshal(ctx.PostBody(), &m); err != nil {
		sendResponse(buildErrorResponse("Error parsing IBC Transfer, please try again"), nil, ctx)
		return
	}

	accountNumber, sequence, err := GetAccountInfo(m.Message.Sender, m.Message.SrcChain)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}
	gas := DefaultGasIbcTransfer
	if m.Transaction.Gas != 0 {
		gas = m.Transaction.Gas
	}

	fee, err := GenerateFeeGasPrice(gas)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
	}

	if m.Message.SrcChain != "EVMOS" && m.Message.DstChain != "EVMOS" {
		sendResponse(buildErrorResponse("Source or destination has to be EVMOS"), nil, ctx)
		return
	}

	height, revision, err := GetHeightInfo(m)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}
	// We timeout the ibc after 500 blocks
	height += 500

	channel, clientID, chainID, prefix, explorerTxURL, err := GetConfigInfo(m)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}

	if err := IsIBCChannelActive("EVMOS", clientID); err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}

	denom := ""

	if m.Message.UseERC20Denom {
		denom, err = GetERC20Address(m.Message.Token)
		if err != nil {
			sendResponse(buildErrorResponse(err.Error()), nil, ctx)
			return
		}
		denom = "erc20/" + denom

	} else {
		denom, err = GetDenom(m.Message.Token, m.Message.SrcChain)
		if err != nil {
			sendResponse(buildErrorResponse(err.Error()), nil, ctx)
			return
		}
	}

	timestampOffset := uint64(time.Now().Add(10 * time.Minute).UnixNano())

	amountInt, ok := sdk.NewIntFromString(m.Message.Amount)
	if !ok {
		sendResponse(buildErrorResponse("Invalid amount"), nil, ctx)
		return
	}

	msgSend := blockchain.CreateMsgTransfer("transfer", channel, amountInt, denom, m.Message.Sender, m.Message.Receiver, revision, height, timestampOffset)

	var eipEncoding blockchain.EipToSignIBC
	tx, err := blockchain.CreateTransactionWithMessage(
		eipEncoding,
		[]sdk.Msg{msgSend},
		"",
		sdk.NewInt(fee),
		feeDenom,
		uint64(gas),
		m.Transaction.PubKey,
		sequence,
		accountNumber,
		chainID,
		m.Transaction.Sender,
		prefix,
	)
	if err != nil {
		sendResponse(buildErrorResponse("Could not create tx, please try again"), nil, ctx)
		return
	}
	resultBytes, err := GetTransactionBytes(tx, accountNumber, chainID, explorerTxURL)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}
	var tx1 TransactionString
	err = json.Unmarshal(resultBytes, &tx1)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	sendResponse(string(resultBytes), err, ctx)
}

type MsgTransferAminoToken struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type MsgTransferAminoTimeoutHeight struct {
	RevisionNumber string `json:"revisionNumber"`
	RevisionHeight string `json:"revisionHeight"`
}
type MsgTransferAmino struct {
	SourcePort       string                        `json:"sourcePort"`
	SourceChannel    string                        `json:"sourceChannel"`
	Sender           string                        `json:"sender"`
	Receiver         string                        `json:"receiver"`
	Token            MsgTransferAminoToken         `json:"token"`
	TimeoutHeight    MsgTransferAminoTimeoutHeight `json:"timeoutHeight"`
	TimeoutTimestamp string                        `json:"timeoutTimestamp"`
}

type ConvertParams struct {
	SrcChain string `json:"srcChain"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
	Token    string `json:"token"`
}

type TxConvertStruct struct {
	Transaction blockchain.TransactionIBCParams `json:"transaction"`
	Message     ConvertParams                   `json:"message"`
}

type DelegateLikeParam struct {
	Amount           string `json:"amount"`
	ValidatorAddress string `json:"validatorAddress"`
}

type TxDelegateLikeStruct struct {
	Transaction blockchain.TransactionIBCParams `json:"transaction"`
	Message     DelegateLikeParam               `json:"message"`
}

func ConvertCoin(ctx *fasthttp.RequestCtx) {
	// convert to erc20
	m := TxConvertStruct{}
	if err := json.Unmarshal(ctx.PostBody(), &m); err != nil {
		sendResponse("", err, ctx)
		return
	}

	accountNumber, sequence, err := GetAccountInfo(m.Transaction.Sender, m.Message.SrcChain)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}

	prefix, chainID, explorerTxURL, err := GetSourceInfo(m.Message.SrcChain)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}

	denom, err := GetDenom(m.Message.Token, m.Message.SrcChain)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}

	amountInt, ok := sdk.NewIntFromString(m.Message.Amount)
	if !ok {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}
	msgSendSdk, err := blockchain.CreateMsgConvertCoin(
		amountInt,
		denom,
		m.Message.Receiver,
		m.Message.Sender,
		prefix,
	)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	fee, err := GenerateFeeGasPrice(EvmosTxFeeConvertAssetGas)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}

	var eipEncoding blockchain.EipToSignMsgSend
	tx, err := blockchain.CreateTransactionWithMessage(
		eipEncoding,
		[]sdk.Msg{msgSendSdk},
		"",
		sdk.NewInt(fee),
		feeDenom,
		uint64(EvmosTxFeeConvertAssetGas),
		m.Transaction.PubKey,
		sequence,
		accountNumber,
		chainID,
		m.Transaction.Sender,
		prefix,
	)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	resultBytes, err := GetTransactionBytes(tx, accountNumber, chainID, explorerTxURL)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}
	sendResponse(string(resultBytes), err, ctx)
}

func ConvertERC20(ctx *fasthttp.RequestCtx) {
	// convert to ibc
	m := TxConvertStruct{}
	if err := json.Unmarshal(ctx.PostBody(), &m); err != nil {
		sendResponse("", err, ctx)
		return
	}

	accountNumber, sequence, err := GetAccountInfo(m.Transaction.Sender, m.Message.SrcChain)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}

	prefix, chainID, explorerTxURL, err := GetSourceInfo(m.Message.SrcChain)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}

	contract, err := GetERC20Address(m.Message.Token)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}

	amountInt, ok := sdk.NewIntFromString(m.Message.Amount)
	if !ok {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}
	msgSendSdk, err := blockchain.CreateMsgConvertERC20(
		amountInt,
		m.Message.Receiver,
		contract,
		m.Message.Sender,
		prefix,
	)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	fee, err := GenerateFeeGasPrice(EvmosTxFeeConvertAssetGas)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}

	var eipEncoding blockchain.EipToSignMsgConvertERC20
	tx, err := blockchain.CreateTransactionWithMessage(
		eipEncoding,
		[]sdk.Msg{msgSendSdk},
		"",
		sdk.NewInt(fee),
		feeDenom,
		uint64(EvmosTxFeeConvertAssetGas),
		m.Transaction.PubKey,
		sequence,
		accountNumber,
		chainID,
		m.Transaction.Sender,
		prefix,
	)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	resultBytes, err := GetTransactionBytes(tx, accountNumber, chainID, explorerTxURL)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	sendResponse(string(resultBytes), err, ctx)
}

func delegateLikeParams(ctx *fasthttp.RequestCtx) (TxDelegateLikeStruct, sdkmath.Int, string, error) {
	m := TxDelegateLikeStruct{}
	if err := json.Unmarshal(ctx.PostBody(), &m); err != nil {
		sendResponse("", err, ctx)
		return m, sdkmath.Int{}, "", fmt.Errorf("invalid params")
	}

	amountInt, ok := sdk.NewIntFromString(m.Message.Amount)
	if !ok {
		sendResponse(buildErrorResponse("Invalid amount"), nil, ctx)
		return m, sdkmath.Int{}, "", fmt.Errorf("invalid amount")
	}

	denom, err := GetDenom(constants.EVMOS, constants.EVMOS)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return m, sdkmath.Int{}, "", fmt.Errorf("invalid denom")
	}
	return m, amountInt, denom, nil
}

func createDelegateLikeTransaction(ctx *fasthttp.RequestCtx, txParams blockchain.TransactionIBCParams, msgs []sdk.Msg, gas float64) {
	accountNumber, sequence, err := GetAccountInfo(txParams.Sender, constants.EVMOS)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	fee, err := GenerateFeeGasPrice(gas)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}

	prefix, chainID, explorerTxURL, err := GetSourceInfo(constants.EVMOS)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}

	var eipEncoding interface{}
	tx, err := blockchain.CreateTransactionWithMessage(
		eipEncoding,
		msgs,
		"",
		sdk.NewInt(fee),
		feeDenom,
		uint64(gas),
		txParams.PubKey,
		sequence,
		accountNumber,
		chainID,
		txParams.Sender,
		prefix,
	)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	resultBytes, err := GetTransactionBytes(tx, accountNumber, chainID, explorerTxURL)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}
	sendResponse(string(resultBytes), err, ctx)
}

func Undelegate(ctx *fasthttp.RequestCtx) {
	m, amountInt, denom, err := delegateLikeParams(ctx)
	if err != nil {
		// The previous function sends the response to the client
		return
	}

	undelegateSDKMsg, err := blockchain.CreateMsgUndelegate(amountInt, m.Transaction.Sender, m.Message.ValidatorAddress, denom)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	createDelegateLikeTransaction(ctx, m.Transaction, []sdk.Msg{undelegateSDKMsg}, DefaultGas)
}

func Delegate(ctx *fasthttp.RequestCtx) {
	m, amountInt, denom, err := delegateLikeParams(ctx)
	if err != nil {
		// The previous function sends the response to the client
		return
	}

	delegateSDKMsg, err := blockchain.CreateMsgDelegate(amountInt, m.Transaction.Sender, m.Message.ValidatorAddress, denom)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	createDelegateLikeTransaction(ctx, m.Transaction, []sdk.Msg{delegateSDKMsg}, DefaultGas)
}

type RedelegateParams struct {
	Amount              string `json:"amount"`
	ValidatorAddress    string `json:"validatorAddress"`
	ValidatorDstAddress string `json:"validatorDstAddress"`
}

type TxRedelegateStruct struct {
	Transaction blockchain.TransactionIBCParams `json:"transaction"`
	Message     RedelegateParams                `json:"message"`
}

func Redelegate(ctx *fasthttp.RequestCtx) {
	m := TxRedelegateStruct{}
	if err := json.Unmarshal(ctx.PostBody(), &m); err != nil {
		sendResponse("", err, ctx)
		return
	}

	amountInt, ok := sdk.NewIntFromString(m.Message.Amount)
	if !ok {
		sendResponse(buildErrorResponse("Invalid amount"), nil, ctx)
		return
	}

	denom, err := GetDenom(constants.EVMOS, constants.EVMOS)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}

	redelegateSDKMsg, err := blockchain.CreateMsgRedelegate(amountInt, m.Transaction.Sender, m.Message.ValidatorAddress, m.Message.ValidatorDstAddress, denom)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	createDelegateLikeTransaction(ctx, m.Transaction, []sdk.Msg{redelegateSDKMsg}, DefaultGas)
}

type RewardsParams struct {
	Validators []string `json:"validators"`
}

type TxRewardsStruct struct {
	Transaction blockchain.TransactionIBCParams `json:"transaction"`
}

func sortRewards(array []StakingReward) []StakingReward {
	sort.Slice(array, func(i, j int) bool {
		a := big.NewInt(0)
		for _, v := range array[i].Reward {
			if v.Denom == feeDenom {
				w := strings.Split(v.Amount, ".")
				_, ok := a.SetString(w[0], 10)
				if !ok {
					a = big.NewInt(0)
				}
				break
			}
		}

		b := big.NewInt(0)
		for _, v := range array[j].Reward {
			if v.Denom == feeDenom {
				w := strings.Split(v.Amount, ".")
				_, ok := b.SetString(w[0], 10)
				if !ok {
					b = big.NewInt(0)
				}
				break
			}
		}

		return a.Cmp(b) == 1
	})

	return array
}

func Rewards(ctx *fasthttp.RequestCtx) {
	m := TxRewardsStruct{}
	if err := json.Unmarshal(ctx.PostBody(), &m); err != nil {
		sendResponse("", err, ctx)
		return
	}

	endpoint := buildThreeParamEndpoint("/cosmos/distribution/v1beta1/delegators/", m.Transaction.Sender, "/rewards")
	rewardsRes, err := getRequestRest("EVMOS", endpoint)
	if err != nil {
		sendResponse("unable to get rewards data", err, ctx)
		return
	}
	var rewards StakingRewardsResponse
	err = json.Unmarshal([]byte(rewardsRes), &rewards)
	if err != nil {
		sendResponse("unable to get rewards data", err, ctx)
		return
	}

	sortedRewards := sortRewards(rewards.Rewards)

	length := 7
	if len(sortedRewards) < 7 {
		length = len(sortedRewards)
	}

	msgs := make([]sdk.Msg, length)

	for k, v := range sortedRewards {
		if k > 6 {
			break
		}

		msg, err := blockchain.CreateMsgRewards(m.Transaction.Sender, v.ValidatorAddress)
		if err != nil {
			sendResponse("", err, ctx)
			return
		}
		msgs[k] = msg
	}

	createDelegateLikeTransaction(ctx, m.Transaction, msgs, DefaultGas*2)
}

type VoteParams struct {
	Option     int `json:"option"`
	ProposalID int `json:"proposalId"`
}

type TxVoteStruct struct {
	Transaction blockchain.TransactionIBCParams `json:"transaction"`
	Message     VoteParams                      `json:"message"`
}

type CancelUndelegationParams struct {
	ValidatorAddress string `json:"validatorAddress"`
	CreationHeight   string `json:"creationHeight"`
	Amount           string `json:"amount"`
}

type TxCancelUndelegationStruct struct {
	Transaction blockchain.TransactionIBCParams `json:"transaction"`
	Message     CancelUndelegationParams        `json:"message"`
}

func Vote(ctx *fasthttp.RequestCtx) {
	m := TxVoteStruct{}
	if err := json.Unmarshal(ctx.PostBody(), &m); err != nil {
		sendResponse(buildErrorResponse("Error parsing Vote, please try again"), nil, ctx)
		return
	}

	msg, err := blockchain.CreateMsgVote(m.Message.ProposalID, m.Message.Option, m.Transaction.Sender)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	gas := DefaultGasIbcTransfer
	if m.Transaction.Gas != 0 {
		gas = m.Transaction.Gas
	}

	createDelegateLikeTransaction(ctx, m.Transaction, []sdk.Msg{msg}, gas)
}

func CancelUndelegation(ctx *fasthttp.RequestCtx) {
	m := TxCancelUndelegationStruct{}
	if err := json.Unmarshal(ctx.PostBody(), &m); err != nil {
		sendResponse("", err, ctx)
		return
	}

	amountInt, ok := sdk.NewIntFromString(m.Message.Amount)
	if !ok {
		sendResponse(buildErrorResponse("Invalid amount"), nil, ctx)
		return
	}

	denom, err := GetDenom(constants.EVMOS, constants.EVMOS)
	if err != nil {
		sendResponse(buildErrorResponse(err.Error()), nil, ctx)
		return
	}

	creationHeight, err := strconv.ParseInt(m.Message.CreationHeight, 10, 64)
	if err != nil {
		sendResponse(buildErrorResponse("error while getting creation height, please try again"), nil, ctx)
		return
	}

	cancelUndelegationSDKMsg, err := blockchain.CreateMsgCancelUndelegations(amountInt, m.Transaction.Sender, m.Message.ValidatorAddress, denom, creationHeight)
	if err != nil {
		sendResponse("", err, ctx)
		return
	}

	createDelegateLikeTransaction(ctx, m.Transaction, []sdk.Msg{cancelUndelegationSDKMsg}, DefaultGas)
}
