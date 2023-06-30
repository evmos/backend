// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package blockchain

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"

	sdkmath "cosmossdk.io/math"

	tx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"

	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	ethsecp256k1 "github.com/evmos/evmos/v12/crypto/ethsecp256k1"
	"github.com/evmos/evmos/v12/ethereum/eip712"
	"github.com/evmos/evmos/v12/types"
	"github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	signing "github.com/cosmos/cosmos-sdk/types/tx/signing"
	crypto "github.com/ethereum/go-ethereum/crypto"
	cryptocodec "github.com/evmos/evmos/v12/crypto/codec"
)

type Transaction struct {
	LegacyAmino       SignData `json:"legacyAmino"`
	SignDirect        SignData `json:"signDirect"`
	EipToSign         string   `json:"eipToSign"`
	MessagingEncoding interface{}
	DataSigningAmino  string `json:"dataSigningAmino"`
}

type SignData struct {
	Body      tx.TxBody   `json:"body"`
	AuthInfo  tx.AuthInfo `json:"authInfo"`
	SignBytes string      `json:"signBytes"`
}

type TransactionIBCParams struct {
	PubKey []uint8 `json:"pubKey"`
	Sender string  `json:"sender"`
	Gas    float64 `json:"gas"`
}

type TransactionParams struct {
	Memo          string  `json:"memo"`
	Fee           string  `json:"fee"`
	Denom         string  `json:"denom"`
	GasLimit      uint64  `json:"gasLimit"`
	PubKey        []uint8 `json:"pubKey"`
	Sequence      uint64  `json:"sequence"`
	AccountNumber uint64  `json:"accountNumber"`
	ChainID       string  `json:"chainId"`
	Sender        string  `json:"sender"`
}

func NewFee(fee sdkmath.Int, denom string, gasLimit uint64) tx.Fee {
	return tx.Fee{Amount: SdkIntToCoins(fee, denom), GasLimit: gasLimit}
}

func NewPubKeyEvmos(key []byte) proto.Message {
	return &ethsecp256k1.PubKey{Key: key}
}

func NewPubKey(key []byte) proto.Message {
	return &secp256k1.PubKey{Key: key}
}

func CreateSignerInfo(publicKey []uint8, sequence uint64, mode signing.SignMode, chainID string) (tx.SignerInfo, error) {
	pubkey := NewPubKey(publicKey)
	if strings.Contains(chainID, "evmos") {
		pubkey = NewPubKeyEvmos(publicKey)
	}

	anyWithValue, err := codectypes.NewAnyWithValue(pubkey)
	if err != nil {
		return tx.SignerInfo{}, err
	}

	single := tx.ModeInfo_Single{Mode: mode}
	modeInfo := tx.ModeInfo{Sum: &tx.ModeInfo_Single_{Single: &single}}
	return tx.SignerInfo{PublicKey: anyWithValue, ModeInfo: &modeInfo, Sequence: sequence}, nil
}

func CreateAuthInfo(signerInfo []*tx.SignerInfo, fee *tx.Fee) tx.AuthInfo {
	return tx.AuthInfo{SignerInfos: signerInfo, Fee: fee}
}

func CreateSignDoc(bodyBytes []byte, authInfoBytes []byte, chainID string, accountNumber uint64) tx.SignDoc {
	return tx.SignDoc{BodyBytes: bodyBytes, AuthInfoBytes: authInfoBytes, ChainId: chainID, AccountNumber: accountNumber}
}

func CreateBodyWithMessage(messages []sdk.Msg, memo string) (tx.TxBody, error) {
	messagesEncoded := make([]*codectypes.Any, len(messages))
	for k, v := range messages {
		anyWithValue, err := codectypes.NewAnyWithValue(v)
		if err != nil {
			return tx.TxBody{}, err
		}
		messagesEncoded[k] = anyWithValue
	}

	return tx.TxBody{Messages: messagesEncoded, Memo: memo}, nil
}

func CreateTransactionWithMessage(
	eipEncoding interface{},
	sdkMessages []sdk.Msg,
	memo string,
	fee sdkmath.Int,
	denom string,
	gasLimit uint64,
	pubKey []uint8,
	sequence uint64,
	accountNumber uint64,
	chainID string,
	sender string,
	prefix string,
) (Transaction, error) {
	body, err := CreateBodyWithMessage(sdkMessages, memo)
	if err != nil {
		return Transaction{}, err
	}

	feeMessage := NewFee(fee, denom, gasLimit)

	// AMINO
	signInfoAmino, err := CreateSignerInfo(pubKey, sequence, signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON, chainID)
	if err != nil {
		return Transaction{}, err
	}

	signInfoAminoArray := []*tx.SignerInfo{&signInfoAmino}
	authInfoAmino := CreateAuthInfo(signInfoAminoArray, &feeMessage)

	messageBytes, err := proto.Marshal(&body)
	if err != nil {
		return Transaction{}, err
	}

	authInfoAminoBytes, err := proto.Marshal(&authInfoAmino)
	if err != nil {
		return Transaction{}, err
	}

	signDocAmino := CreateSignDoc(messageBytes, authInfoAminoBytes, chainID, accountNumber)

	signDocAminoBytes, err := proto.Marshal(&signDocAmino)
	if err != nil {
		return Transaction{}, err
	}

	toSignAmino := crypto.Keccak256(signDocAminoBytes)
	encodeSignAmino := base64.StdEncoding.EncodeToString(toSignAmino)

	// SignDirect
	signInfoDirect, err := CreateSignerInfo(pubKey, sequence, signing.SignMode_SIGN_MODE_DIRECT, chainID)
	if err != nil {
		return Transaction{}, err
	}

	signInfoDirectArray := []*tx.SignerInfo{&signInfoDirect}

	authInfoDirect := CreateAuthInfo(signInfoDirectArray, &feeMessage)

	authInfoDirectBytes, err := proto.Marshal(&authInfoDirect)
	if err != nil {
		return Transaction{}, err
	}

	signDocDirect := CreateSignDoc(messageBytes, authInfoDirectBytes, chainID, accountNumber)

	signDocDirectBytes, err := proto.Marshal(&signDocDirect)
	if err != nil {
		return Transaction{}, err
	}

	toSignDirect := crypto.Keccak256(signDocDirectBytes)
	encodeSignDirect := base64.StdEncoding.EncodeToString(toSignDirect)

	bytes := []byte{}

	gasAmount := sdk.NewCoins(sdk.NewCoin(denom, fee))

	//nolint:staticcheck
	feeSdk := legacytx.NewStdFee(gasLimit, gasAmount)
	// TODO: use AuxTxBuilder
	dataAmino := legacytx.StdSignBytes(chainID, accountNumber, sequence, 0, feeSdk, sdkMessages, memo, nil)

	if strings.Contains(chainID, "evmos") {
		// GenerateTypedData TypedData
		var ethermintCodec codec.ProtoCodecMarshaler
		registry := codectypes.NewInterfaceRegistry()
		types.RegisterInterfaces(registry)
		ethermintCodec = codec.NewProtoCodec(registry)
		cryptocodec.RegisterInterfaces(registry)

		from, err := Bech32StringToAddress(sender, prefix)
		if err != nil {
			return Transaction{}, err
		}

		ethChainID, err := types.ParseChainID(chainID)
		if err != nil {
			return Transaction{}, err
		}

		// NOTE: eip712 only works for more than 1 message if all the messages in the array are the same type
		typedData, err := eip712.LegacyWrapTxToTypedData(ethermintCodec, ethChainID.Uint64(), sdkMessages[0], dataAmino, &eip712.FeeDelegationOptions{
			FeePayer: from,
		})
		if err != nil {
			return Transaction{}, err
		}

		bytes, err = json.Marshal(typedData)
		if err != nil {
			return Transaction{}, err
		}
	}

	return Transaction{
		LegacyAmino: SignData{
			Body:      body,
			AuthInfo:  authInfoAmino,
			SignBytes: encodeSignAmino,
		},
		SignDirect: SignData{
			Body:      body,
			AuthInfo:  authInfoDirect,
			SignBytes: encodeSignDirect,
		},
		EipToSign:         string(bytes),
		MessagingEncoding: eipEncoding,
		DataSigningAmino:  string(dataAmino),
	}, nil
}

func Remove0xFromHex(signature string) string {
	signature = strings.TrimPrefix(signature, "0x")
	return signature
}

func SignatureToWeb3Extension(chain uint64, sender string, signature string) (types.ExtensionOptionsWeb3Tx, error) {
	sign, err := hex.DecodeString(Remove0xFromHex(signature))
	if err != nil {
		return types.ExtensionOptionsWeb3Tx{}, err
	}

	return types.ExtensionOptionsWeb3Tx{TypedDataChainID: chain, FeePayer: sender, FeePayerSig: sign}, nil
}

func CreateTxRawEIP712(body *tx.TxBody, authInfo tx.AuthInfo, extension types.ExtensionOptionsWeb3Tx) (tx.TxRaw, error) {
	anyWithValue, err := codectypes.NewAnyWithValue(&extension)
	if err != nil {
		return tx.TxRaw{}, err
	}

	body.ExtensionOptions = []*codectypes.Any{anyWithValue}

	bodyBytes, err := proto.Marshal(body)
	if err != nil {
		return tx.TxRaw{}, err
	}

	authInfoBytes, err := proto.Marshal(&authInfo)
	if err != nil {
		return tx.TxRaw{}, err
	}

	return tx.TxRaw{BodyBytes: bodyBytes, AuthInfoBytes: authInfoBytes, Signatures: make([][]byte, 1)}, nil
}

func JoinEipWithSignature(chain uint64, feePayer string, feePayerSig string, body *tx.TxBody, authInfo tx.AuthInfo) (tx.TxRaw, error) {
	extension, err := SignatureToWeb3Extension(chain, feePayer, feePayerSig)
	if err != nil {
		return tx.TxRaw{}, err
	}

	txRaw, err := CreateTxRawEIP712(body, authInfo, extension)
	if err != nil {
		return tx.TxRaw{}, err
	}

	return txRaw, nil
}
