package v2

import (
	"encoding/json"
	"fmt"

	"github.com/tharsis/dashboard-backend/internal/v2/encoding"
	"github.com/tharsis/dashboard-backend/internal/v2/node/rest"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/valyala/fasthttp"
)

// BroadcastTxParams represents the parameters for the POST /v2/tx/broadcast endpoint.
type BroadcastTxParams struct {
	// which network should the transaction be broadcasted to
	Network string `json:"network"`
	// the signed transaction to be broadcasted
	TxBytes []byte `json:"tx_bytes"`
}

type BroadcastTxResponse struct {
	Code   uint32 `json:"code"`
	TxHash string `json:"tx_hash"`
	RawLog string `json:"raw_log"`
}

// BroadcastTx handles POST /tx/broadcast.
// It broadcasts a signed transaction synchronously to the specified network.
// Returns:
//
//	{
//	  "txhash": "3CB7FCC9F5FB31E530CC15665F3FD655AE6CB56CDACAD58D1395C68EDD50D0BB",
//	  "code": 0,
//	  "raw_log": "[]",
//	}
func (h *Handler) BroadcastTx(ctx *fasthttp.RequestCtx) {
	reqParams := BroadcastTxParams{}
	if err := json.Unmarshal(ctx.PostBody(), &reqParams); err != nil {
		ctx.Logger().Printf("Error decoding request body: %s", err.Error())
		sendBadRequestResponse(ctx, "Invalid request body")
		return
	}

	txRequest := tx.BroadcastTxRequest{
		TxBytes: reqParams.TxBytes,
		Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
	}

	jsonTxRequest, err := json.Marshal(txRequest)
	if err != nil {
		ctx.Logger().Printf("Error marshaling txRequest: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	err = ValidateBroadcastTxParams(&reqParams)
	if err != nil {
		sendBadRequestResponse(ctx, err.Error())
		return
	}

	restClient, err := rest.NewClient(reqParams.Network)
	if err != nil {
		ctx.Logger().Printf("Error creating rest client: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	txResponse, err := restClient.BroadcastTx(jsonTxRequest)
	if err != nil {
		ctx.Logger().Printf("Error broadcasting tx: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	response := BroadcastTxResponse{
		Code:   txResponse.TxResponse.Code,
		TxHash: txResponse.TxResponse.TxHash,
		RawLog: txResponse.TxResponse.RawLog,
	}
	sendSuccessfulJSONResponse(ctx, &response)
}

func ValidateBroadcastTxParams(params *BroadcastTxParams) error {
	// TODO: validate network by checking if it's in the list of available networks
	if params.Network == "" {
		return fmt.Errorf("network cannot be empty")
	}
	if len(params.TxBytes) == 0 {
		return fmt.Errorf("tx_bytes cannot be empty")
	}
	return nil
}

// BroadcastTxParams represents the parameters for the POST /v2/tx/broadcast endpoint.
type BroadcastAminoTxParams struct {
	// which network should the transaction be broadcasted to
	Network   string                `json:"network"`
	Signed    legacytx.StdSignDoc   `json:"signed"`
	Signature legacytx.StdSignature `json:"signature"` //nolint:staticcheck
}

type BroadcastAminoTxResponse struct {
	Code   uint32 `json:"code"`
	TxHash string `json:"tx_hash"`
	RawLog string `json:"raw_log"`
}

// BroadcastAminoTx handles POST /tx/amino/broadcast.
// It broadcasts a signed transaction synchronously to the specified network.
// It receives StdSignDoc and StdSignature as input and builds a TxBuilder to generate
// the broadcast bytes.
// Returns:
//
//	{
//	  "txhash": "3CB7FCC9F5FB31E530CC15665F3FD655AE6CB56CDACAD58D1395C68EDD50D0BB",
//	  "code": 0,
//	  "raw_log": "[]",
//	}
func (h *Handler) BroadcastAminoTx(ctx *fasthttp.RequestCtx) {
	protoCfg := encoding.MakeEncodingConfig()
	aminoCodec := protoCfg.Amino

	reqParams := BroadcastAminoTxParams{}
	if err := aminoCodec.Amino.UnmarshalJSON(ctx.PostBody(), &reqParams); err != nil {
		ctx.Logger().Printf("Error decoding request body: %s", err.Error())
		sendBadRequestResponse(ctx, "Invalid request body")
		return
	}

	txBytes, err := EncodeTransaction(&protoCfg, reqParams.Signed, reqParams.Signature)
	if err != nil {
		ctx.Logger().Printf("Error generating tx bytes: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	txRequest := tx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
	}

	jsonTxRequest, err := json.Marshal(txRequest)
	if err != nil {
		ctx.Logger().Printf("Error marshaling txRequest: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	restClient, err := rest.NewClient(reqParams.Network)
	if err != nil {
		ctx.Logger().Printf("Error creating rest client: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	txResponse, err := restClient.BroadcastTx(jsonTxRequest)
	if err != nil {
		ctx.Logger().Printf("Error broadcasting tx: %s", err.Error())
		sendInternalErrorResponse(ctx)
		return
	}

	response := BroadcastTxResponse{
		Code:   txResponse.TxResponse.Code,
		TxHash: txResponse.TxResponse.TxHash,
		RawLog: txResponse.TxResponse.RawLog,
	}
	sendSuccessfulJSONResponse(ctx, &response)
}

// EncodeTransaction encodes the upcoming transaction using the provided configuration.
// It receives StdSignDoc and StdSignature as input and builds a TxBuilder to generate
// the broadcast bytes.
func EncodeTransaction(encConfig *params.EncodingConfig, signDoc legacytx.StdSignDoc, signature legacytx.StdSignature) ([]byte, error) { //nolint:staticcheck
	txBuilder := encConfig.TxConfig.NewTxBuilder()
	aminoCodec := encConfig.Amino
	var fees legacytx.StdFee
	if err := aminoCodec.UnmarshalJSON(signDoc.Fee, &fees); err != nil {
		return nil, err
	}

	// Validate payload messages
	msgs := make([]sdk.Msg, len(signDoc.Msgs))
	for i, jsonMsg := range signDoc.Msgs {
		var m sdk.Msg
		if err := aminoCodec.UnmarshalJSON(jsonMsg, &m); err != nil {
			return nil, err
		}
		msgs[i] = m
	}

	err := txBuilder.SetMsgs(msgs...)
	if err != nil {
		return nil, err
	}

	// Build transaction params
	txBuilder.SetMemo(signDoc.Memo)
	txBuilder.SetFeeAmount(fees.Amount)
	txBuilder.SetFeePayer(sdk.AccAddress(fees.Payer))
	txBuilder.SetFeeGranter(sdk.AccAddress(fees.Granter))
	txBuilder.SetGasLimit(fees.Gas)
	txBuilder.SetTimeoutHeight(signDoc.TimeoutHeight)

	sigV2, err := legacytx.StdSignatureToSignatureV2(aminoCodec, signature)
	if err != nil {
		return nil, err
	}
	sigV2.Sequence = signDoc.Sequence

	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		return nil, err
	}

	txBytes, err := encConfig.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}
	return txBytes, nil
}
