// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package blockchain

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

func TestCreateMsgConvertERC20(t *testing.T) {
	_, err := CreateMsgConvertERC20(
		sdk.NewInt(1),
		"evmos14uepnqnvkuyyvwe65wmncejq5g2f0tjft3wr65",
		"0x3452e23F9c4cC62c70B7ADAd699B264AF3549C19",
		"0xaF3219826Cb708463B3AA3B73c6640A21497AE49",
		"evmos")
	if err != nil {
		t.Fatalf("Error creating msgSend")
	}
}

func TestTxWithConvertERC20(t *testing.T) {
	msgSendSdk, err := CreateMsgConvertERC20(
		sdk.NewInt(1),
		"evmos14uepnqnvkuyyvwe65wmncejq5g2f0tjft3wr65",
		"0x3452e23F9c4cC62c70B7ADAd699B264AF3549C19",
		"0xaF3219826Cb708463B3AA3B73c6640A21497AE49",
		"evmos")
	if err != nil {
		t.Fatalf("Error creating msgSend")
	}

	pubKey, err := base64.StdEncoding.DecodeString("Ak8wUTcElcOofCZZJM97pduO+Aw3w4wzClrJgN2VzTVQ")
	if err != nil {
		t.Fatalf("Error decoding string represented by base64 to bytes")
	}

	tx, err := CreateTransactionWithMessage(
		EipToSignMsgSend{},
		[]sdk.Msg{msgSendSdk},
		"",
		sdk.NewInt(241499999999999968),
		"aevmos",
		uint64(10500000),
		pubKey,
		uint64(43),
		uint64(2164290),
		"evmos_9001-2",
		"evmos14uepnqnvkuyyvwe65wmncejq5g2f0tjft3wr65",
		"evmos",
	)
	if err != nil {
		t.Fatalf("Error creating transaction with message")
	}

	bodyBytes, err := proto.Marshal(&tx.SignDirect.Body)
	if err != nil {
		t.Fatalf("Error Sign Direct: encoding body to []bytes: %q", err)
	}

	bodyBytesSignDirectFrontend := []byte{ //nolint:all
		10, 173, 1, 10, 31, 47, 101, 118, 109, 111, 115, 46,
		101, 114, 99, 50, 48, 46, 118, 49, 46, 77, 115, 103,
		67, 111, 110, 118, 101, 114, 116, 69, 82, 67, 50, 48,
		18, 137, 1, 10, 42, 48, 120, 51, 52, 53, 50, 101,
		50, 51, 70, 57, 99, 52, 99, 67, 54, 50, 99, 55,
		48, 66, 55, 65, 68, 65, 100, 54, 57, 57, 66, 50,
		54, 52, 65, 70, 51, 53, 52, 57, 67, 49, 57, 18,
		1, 49, 26, 44, 101, 118, 109, 111, 115, 49, 52, 117,
		101, 112, 110, 113, 110, 118, 107, 117, 121, 121, 118, 119,
		101, 54, 53, 119, 109, 110, 99, 101, 106, 113, 53, 103,
		50, 102, 48, 116, 106, 102, 116, 51, 119, 114, 54, 53,
		34, 42, 48, 120, 97, 70, 51, 50, 49, 57, 56, 50,
		54, 67, 98, 55, 48, 56, 52, 54, 51, 66, 51, 65,
		65, 51, 66, 55, 51, 99, 54, 54, 52, 48, 65, 50,
		49, 52, 57, 55, 65, 69, 52, 57,
	}

	for k, v := range bodyBytes {
		if v != bodyBytesSignDirectFrontend[k] {
			t.Fatalf("Error Sign Direct: bodyBytes and bodyBytesSignDirectFrontend are not equals")
		}
	}

	authInfoBytesSignDirect, err := proto.Marshal(&tx.SignDirect.AuthInfo)
	if err != nil {
		t.Fatalf("Error Sign Direct: encoding authInfo to []bytes: %q", err)
	}

	authInfoBytesSignDirectFrontend := []byte{
		10, 89, 10, 79, 10, 40, 47, 101, 116, 104, 101, 114,
		109, 105, 110, 116, 46, 99, 114, 121, 112, 116, 111, 46,
		118, 49, 46, 101, 116, 104, 115, 101, 99, 112, 50, 53,
		54, 107, 49, 46, 80, 117, 98, 75, 101, 121, 18, 35,
		10, 33, 2, 79, 48, 81, 55, 4, 149, 195, 168, 124,
		38, 89, 36, 207, 123, 165, 219, 142, 248, 12, 55, 195,
		140, 51, 10, 90, 201, 128, 221, 149, 205, 53, 80, 18,
		4, 10, 2, 8, 1, 24, 43, 18, 35, 10, 28, 10,
		6, 97, 101, 118, 109, 111, 115, 18, 18, 50, 52, 49,
		52, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57,
		57, 54, 56, 16, 160, 239, 128, 5,
	}

	for k, v := range authInfoBytesSignDirect {
		if v != authInfoBytesSignDirectFrontend[k] {
			t.Fatalf("Error Sign Direct: authInfoBytesSignDirect and authInfoBytesSignDirectFrontend are not equals")
		}
	}

	signBytesSignDirectFrontend := "azX1aDLI0/6akqWnNy6Ne/Z8OMOXADbqFFA06jmVx1w="

	if tx.SignDirect.SignBytes != signBytesSignDirectFrontend {
		t.Fatalf("Error Sign Direct: tx.SignDirect.SignBytes and signBytesSignDirectFrontend are not equals")
	}

	bodyBytesLegacyAmino, err := proto.Marshal(&tx.LegacyAmino.Body)
	if err != nil {
		t.Fatalf("Error Legacy Amino: encoding body to []bytes from: %q", err)
	}

	bodyBtyesLegacyAminoFrontend := []byte{ //nolint:all
		10, 173, 1, 10, 31, 47, 101, 118, 109, 111, 115, 46,
		101, 114, 99, 50, 48, 46, 118, 49, 46, 77, 115, 103,
		67, 111, 110, 118, 101, 114, 116, 69, 82, 67, 50, 48,
		18, 137, 1, 10, 42, 48, 120, 51, 52, 53, 50, 101,
		50, 51, 70, 57, 99, 52, 99, 67, 54, 50, 99, 55,
		48, 66, 55, 65, 68, 65, 100, 54, 57, 57, 66, 50,
		54, 52, 65, 70, 51, 53, 52, 57, 67, 49, 57, 18,
		1, 49, 26, 44, 101, 118, 109, 111, 115, 49, 52, 117,
		101, 112, 110, 113, 110, 118, 107, 117, 121, 121, 118, 119,
		101, 54, 53, 119, 109, 110, 99, 101, 106, 113, 53, 103,
		50, 102, 48, 116, 106, 102, 116, 51, 119, 114, 54, 53,
		34, 42, 48, 120, 97, 70, 51, 50, 49, 57, 56, 50,
		54, 67, 98, 55, 48, 56, 52, 54, 51, 66, 51, 65,
		65, 51, 66, 55, 51, 99, 54, 54, 52, 48, 65, 50,
		49, 52, 57, 55, 65, 69, 52, 57,
	}

	for k, v := range bodyBytesLegacyAmino {
		if v != bodyBtyesLegacyAminoFrontend[k] {
			t.Fatalf("Error Legacy Amino: bodyBytesLegacyAmino and bodyBtyesLegacyAminoFrontend are not equals")
		}
	}

	authInfoBytesLegacyAmino, err := proto.Marshal(&tx.LegacyAmino.AuthInfo)
	if err != nil {
		t.Fatalf("Error Legacy Amino: encoding authInfo to []bytes: %q", err)
	}

	authInfoBytesLegacyAminoFrontend := []byte{
		10, 89, 10, 79, 10, 40, 47, 101, 116, 104, 101, 114,
		109, 105, 110, 116, 46, 99, 114, 121, 112, 116, 111, 46,
		118, 49, 46, 101, 116, 104, 115, 101, 99, 112, 50, 53,
		54, 107, 49, 46, 80, 117, 98, 75, 101, 121, 18, 35,
		10, 33, 2, 79, 48, 81, 55, 4, 149, 195, 168, 124,
		38, 89, 36, 207, 123, 165, 219, 142, 248, 12, 55, 195,
		140, 51, 10, 90, 201, 128, 221, 149, 205, 53, 80, 18,
		4, 10, 2, 8, 127, 24, 43, 18, 35, 10, 28, 10,
		6, 97, 101, 118, 109, 111, 115, 18, 18, 50, 52, 49,
		52, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57,
		57, 54, 56, 16, 160, 239, 128, 5,
	}

	for k, v := range authInfoBytesLegacyAmino {
		if v != authInfoBytesLegacyAminoFrontend[k] {
			t.Fatalf("Error Legacy Amino: authInfoBytesLegacyAmino and authInfoBytesLegacyAminoFrontend are not equals")
		}
	}

	signBytesLegacyAminoFrontend := "ksrlM/uvpuPVT8MIgMHLcp/4gvaljSWsKtDU3q/5svY="

	if tx.LegacyAmino.SignBytes != signBytesLegacyAminoFrontend {
		t.Fatalf("Error Legacy Amino: tx.LegacyAmino.SignBytes and signBytesLegacyAminoFrontend are not equals")
	}

	// 	// 	// For the test I changed the chain_id in domain: 9001 for "0x2329"
	eipToSignFrontend := `{"types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"string"},{"name":"salt","type":"string"}],"Tx":[{"name":"account_number","type":"string"},{"name":"chain_id","type":"string"},{"name":"fee","type":"Fee"},{"name":"memo","type":"string"},{"name":"msgs","type":"Msg[]"},{"name":"sequence","type":"string"}],"Fee":[{"name":"feePayer","type":"string"},{"name":"amount","type":"Coin[]"},{"name":"gas","type":"string"}],"Coin":[{"name":"denom","type":"string"},{"name":"amount","type":"string"}],"Msg":[{"name":"type","type":"string"},{"name":"value","type":"MsgValue"}],"MsgValue":[{"name":"contract_address","type":"string"},{"name":"amount","type":"string"},{"name":"receiver","type":"string"},{"name":"sender","type":"string"}]},"primaryType":"Tx","domain":{"name":"Cosmos Web3","version":"1.0.0","chainId":"0x2329","verifyingContract":"cosmos","salt":"0"},"message":{"account_number":"2164290","chain_id":"evmos_9001-2","fee":{"amount":[{"amount":"241499999999999968","denom":"aevmos"}],"gas":"10500000","feePayer":"evmos14uepnqnvkuyyvwe65wmncejq5g2f0tjft3wr65"},"memo":"","msgs":[{"type":"evmos/MsgConvertERC20","value":{"contract_address":"0x3452e23F9c4cC62c70B7ADAd699B264AF3549C19","amount":"1","receiver":"evmos14uepnqnvkuyyvwe65wmncejq5g2f0tjft3wr65","sender":"0xaF3219826Cb708463B3AA3B73c6640A21497AE49"}}],"sequence":"43"}}`
	var resEipToSign EipToSignMsgConvertERC20
	err = json.Unmarshal([]byte(tx.EipToSign), &resEipToSign)
	if err != nil {
		t.Fatalf("Error: couldn't unmarshal tx.EipToSign: %q", err)
	}

	var resEipToSignFrontend EipToSignMsgConvertERC20
	err = json.Unmarshal([]byte(eipToSignFrontend), &resEipToSignFrontend)
	if err != nil {
		t.Fatalf("Error decoding string to []byte: %q", err)
	}

	resEipToSignBytes, err := json.Marshal(resEipToSign)
	if err != nil {
		t.Fatalf("Error: couldn't marshal resEipToSign: %q", err)
	}

	resEipToSignBytesFrontend, err := json.Marshal(resEipToSignFrontend)
	if err != nil {
		t.Fatalf("Error: couldn't marshal resEipToSignFrontend: %q", err)
	}

	for k, v := range resEipToSignBytes {
		if v != resEipToSignBytesFrontend[k] {
			t.Fatalf("Error EipToSign: resEipToSignBytes and resEipToSignBytesFrontend are not equals")
		}
	}
}

func TestCreateMsgConvertCoin(t *testing.T) {
	_, err := CreateMsgConvertCoin(
		sdk.NewInt(1),
		"ibc/448C1061CE97D86CC5E86374CD914870FB8EBA16C58661B5F1D3F46729A2422D", // JUNO
		"0xaF3219826Cb708463B3AA3B73c6640A21497AE49",
		"evmos14uepnqnvkuyyvwe65wmncejq5g2f0tjft3wr65",
		"evmos")
	if err != nil {
		t.Fatalf("Error creating msgSend")
	}
}

func TestTxWithConvertCoin(t *testing.T) {
	msgSendSdk, err := CreateMsgConvertCoin(
		sdk.NewInt(1),
		"ibc/448C1061CE97D86CC5E86374CD914870FB8EBA16C58661B5F1D3F46729A2422D", // JUNO
		"0xaF3219826Cb708463B3AA3B73c6640A21497AE49",
		"evmos14uepnqnvkuyyvwe65wmncejq5g2f0tjft3wr65",
		"evmos")
	if err != nil {
		t.Fatalf("Error creating msgSend")
	}

	pubKey, err := base64.StdEncoding.DecodeString("Ak8wUTcElcOofCZZJM97pduO+Aw3w4wzClrJgN2VzTVQ")
	if err != nil {
		t.Fatalf("Error decoding string represented by base64 to bytes")
	}

	tx, err := CreateTransactionWithMessage(
		EipToSignMsgSend{},
		[]sdk.Msg{msgSendSdk},
		"",
		sdk.NewInt(241499999999999968),
		"aevmos",
		uint64(10500000),
		pubKey,
		uint64(43),
		uint64(2164290),
		"evmos_9001-2",
		"evmos14uepnqnvkuyyvwe65wmncejq5g2f0tjft3wr65",
		"evmos",
	)
	if err != nil {
		t.Fatalf("Error creating transaction with message")
	}

	bodyBytes, err := proto.Marshal(&tx.SignDirect.Body)
	if err != nil {
		t.Fatalf("Error Sign Direct: encoding body to []bytes: %q", err)
	}

	bodyBytesSignDirectFrontend := []byte{ //nolint:all
		10, 200, 1, 10, 30, 47, 101, 118, 109, 111, 115, 46,
		101, 114, 99, 50, 48, 46, 118, 49, 46, 77, 115, 103,
		67, 111, 110, 118, 101, 114, 116, 67, 111, 105, 110, 18,
		165, 1, 10, 73, 10, 68, 105, 98, 99, 47, 52, 52,
		56, 67, 49, 48, 54, 49, 67, 69, 57, 55, 68, 56,
		54, 67, 67, 53, 69, 56, 54, 51, 55, 52, 67, 68,
		57, 49, 52, 56, 55, 48, 70, 66, 56, 69, 66, 65,
		49, 54, 67, 53, 56, 54, 54, 49, 66, 53, 70, 49,
		68, 51, 70, 52, 54, 55, 50, 57, 65, 50, 52, 50,
		50, 68, 18, 1, 49, 18, 42, 48, 120, 97, 70, 51,
		50, 49, 57, 56, 50, 54, 67, 98, 55, 48, 56, 52,
		54, 51, 66, 51, 65, 65, 51, 66, 55, 51, 99, 54,
		54, 52, 48, 65, 50, 49, 52, 57, 55, 65, 69, 52,
		57, 26, 44, 101, 118, 109, 111, 115, 49, 52, 117, 101,
		112, 110, 113, 110, 118, 107, 117, 121, 121, 118, 119, 101,
		54, 53, 119, 109, 110, 99, 101, 106, 113, 53, 103, 50,
		102, 48, 116, 106, 102, 116, 51, 119, 114, 54, 53,
	}

	for k, v := range bodyBytes {
		if v != bodyBytesSignDirectFrontend[k] {
			t.Fatalf("Error Sign Direct: bodyBytes and bodyBytesSignDirectFrontend are not equals")
		}
	}

	authInfoBytesSignDirect, err := proto.Marshal(&tx.SignDirect.AuthInfo)
	if err != nil {
		t.Fatalf("Error Sign Direct: encoding authInfo to []bytes: %q", err)
	}

	authInfoBytesSignDirectFrontend := []byte{
		10, 89, 10, 79, 10, 40, 47, 101, 116, 104, 101, 114,
		109, 105, 110, 116, 46, 99, 114, 121, 112, 116, 111, 46,
		118, 49, 46, 101, 116, 104, 115, 101, 99, 112, 50, 53,
		54, 107, 49, 46, 80, 117, 98, 75, 101, 121, 18, 35,
		10, 33, 2, 79, 48, 81, 55, 4, 149, 195, 168, 124,
		38, 89, 36, 207, 123, 165, 219, 142, 248, 12, 55, 195,
		140, 51, 10, 90, 201, 128, 221, 149, 205, 53, 80, 18,
		4, 10, 2, 8, 1, 24, 43, 18, 35, 10, 28, 10,
		6, 97, 101, 118, 109, 111, 115, 18, 18, 50, 52, 49,
		52, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57,
		57, 54, 56, 16, 160, 239, 128, 5,
	}

	for k, v := range authInfoBytesSignDirect {
		if v != authInfoBytesSignDirectFrontend[k] {
			t.Fatalf("Error Sign Direct: authInfoBytesSignDirect and authInfoBytesSignDirectFrontend are not equals")
		}
	}
	signBytesSignDirectFrontend := "rr5mUTjXd5yrxDwaPfPN3e6pYOJGZnccI04vH3vZP/8="

	if tx.SignDirect.SignBytes != signBytesSignDirectFrontend {
		t.Fatalf("Error Sign Direct: tx.SignDirect.SignBytes and signBytesSignDirectFrontend are not equals")
	}

	bodyBytesLegacyAmino, err := proto.Marshal(&tx.LegacyAmino.Body)
	if err != nil {
		t.Fatalf("Error Legacy Amino: encoding body to []bytes from: %q", err)
	}

	bodyBtyesLegacyAminoFrontend := []byte{ //nolint:all
		10, 200, 1, 10, 30, 47, 101, 118, 109, 111, 115, 46,
		101, 114, 99, 50, 48, 46, 118, 49, 46, 77, 115, 103,
		67, 111, 110, 118, 101, 114, 116, 67, 111, 105, 110, 18,
		165, 1, 10, 73, 10, 68, 105, 98, 99, 47, 52, 52,
		56, 67, 49, 48, 54, 49, 67, 69, 57, 55, 68, 56,
		54, 67, 67, 53, 69, 56, 54, 51, 55, 52, 67, 68,
		57, 49, 52, 56, 55, 48, 70, 66, 56, 69, 66, 65,
		49, 54, 67, 53, 56, 54, 54, 49, 66, 53, 70, 49,
		68, 51, 70, 52, 54, 55, 50, 57, 65, 50, 52, 50,
		50, 68, 18, 1, 49, 18, 42, 48, 120, 97, 70, 51,
		50, 49, 57, 56, 50, 54, 67, 98, 55, 48, 56, 52,
		54, 51, 66, 51, 65, 65, 51, 66, 55, 51, 99, 54,
		54, 52, 48, 65, 50, 49, 52, 57, 55, 65, 69, 52,
		57, 26, 44, 101, 118, 109, 111, 115, 49, 52, 117, 101,
		112, 110, 113, 110, 118, 107, 117, 121, 121, 118, 119, 101,
		54, 53, 119, 109, 110, 99, 101, 106, 113, 53, 103, 50,
		102, 48, 116, 106, 102, 116, 51, 119, 114, 54, 53,
	}

	for k, v := range bodyBytesLegacyAmino {
		if v != bodyBtyesLegacyAminoFrontend[k] {
			t.Fatalf("Error Legacy Amino: bodyBytesLegacyAmino and bodyBtyesLegacyAminoFrontend are not equals")
		}
	}

	authInfoBytesLegacyAmino, err := proto.Marshal(&tx.LegacyAmino.AuthInfo)
	if err != nil {
		t.Fatalf("Error Legacy Amino: encoding authInfo to []bytes: %q", err)
	}

	authInfoBytesLegacyAminoFrontend := []byte{
		10, 89, 10, 79, 10, 40, 47, 101, 116, 104, 101, 114,
		109, 105, 110, 116, 46, 99, 114, 121, 112, 116, 111, 46,
		118, 49, 46, 101, 116, 104, 115, 101, 99, 112, 50, 53,
		54, 107, 49, 46, 80, 117, 98, 75, 101, 121, 18, 35,
		10, 33, 2, 79, 48, 81, 55, 4, 149, 195, 168, 124,
		38, 89, 36, 207, 123, 165, 219, 142, 248, 12, 55, 195,
		140, 51, 10, 90, 201, 128, 221, 149, 205, 53, 80, 18,
		4, 10, 2, 8, 127, 24, 43, 18, 35, 10, 28, 10,
		6, 97, 101, 118, 109, 111, 115, 18, 18, 50, 52, 49,
		52, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57,
		57, 54, 56, 16, 160, 239, 128, 5,
	}

	for k, v := range authInfoBytesLegacyAmino {
		if v != authInfoBytesLegacyAminoFrontend[k] {
			t.Fatalf("Error Legacy Amino: authInfoBytesLegacyAmino and authInfoBytesLegacyAminoFrontend are not equals")
		}
	}

	signBytesLegacyAminoFrontend := "a7ZY5t/5U2yEwZdG0yWhIjUGLr8JeU87EtYd6m6uhts="

	if tx.LegacyAmino.SignBytes != signBytesLegacyAminoFrontend {
		t.Fatalf("Error Legacy Amino: tx.LegacyAmino.SignBytes and signBytesLegacyAminoFrontend are not equals")
	}

	// 	// For the test I changed the chain_id in domain: 9001 for "0x2329"
	eipToSignFrontend := `{"types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"string"},{"name":"salt","type":"string"}],"Tx":[{"name":"account_number","type":"string"},{"name":"chain_id","type":"string"},{"name":"fee","type":"Fee"},{"name":"memo","type":"string"},{"name":"msgs","type":"Msg[]"},{"name":"sequence","type":"string"}],"Fee":[{"name":"feePayer","type":"string"},{"name":"amount","type":"Coin[]"},{"name":"gas","type":"string"}],"Coin":[{"name":"denom","type":"string"},{"name":"amount","type":"string"}],"Msg":[{"name":"type","type":"string"},{"name":"value","type":"MsgValue"}],"MsgValue":[{"name":"coin","type":"TypeCoin"},{"name":"receiver","type":"string"},{"name":"sender","type":"string"}],"TypeCoin":[{"name":"denom","type":"string"},{"name":"amount","type":"string"}]},"primaryType":"Tx","domain":{"name":"Cosmos Web3","version":"1.0.0","chainId":"0x2329","verifyingContract":"cosmos","salt":"0"},"message":{"account_number":"2164290","chain_id":"evmos_9001-2","fee":{"amount":[{"amount":"241499999999999968","denom":"aevmos"}],"gas":"10500000","feePayer":"evmos14uepnqnvkuyyvwe65wmncejq5g2f0tjft3wr65"},"memo":"","msgs":[{"type":"evmos/MsgConvertCoin","value":{"coin":{"denom":"ibc/448C1061CE97D86CC5E86374CD914870FB8EBA16C58661B5F1D3F46729A2422D","amount":"1"},"receiver":"0xaF3219826Cb708463B3AA3B73c6640A21497AE49","sender":"evmos14uepnqnvkuyyvwe65wmncejq5g2f0tjft3wr65"}}],"sequence":"43"}}`

	var resEipToSign EipToSignMsgSend
	err = json.Unmarshal([]byte(tx.EipToSign), &resEipToSign)
	if err != nil {
		t.Fatalf("Error: couldn't unmarshal tx.EipToSign: %q", err)
	}

	var resEipToSignFrontend EipToSignMsgSend
	err = json.Unmarshal([]byte(eipToSignFrontend), &resEipToSignFrontend)
	if err != nil {
		t.Fatalf("Error decoding string to []byte: %q", err)
	}

	resEipToSignBytes, err := json.Marshal(resEipToSign)
	if err != nil {
		t.Fatalf("Error: couldn't marshal resEipToSign: %q", err)
	}

	resEipToSignBytesFrontend, err := json.Marshal(resEipToSignFrontend)
	if err != nil {
		t.Fatalf("Error: couldn't marshal resEipToSignFrontend: %q", err)
	}

	for k, v := range resEipToSignBytes {
		if v != resEipToSignBytesFrontend[k] {
			t.Fatalf("Error EipToSign: resEipToSignBytes and resEipToSignBytesFrontend are not equals")
		}
	}
}
