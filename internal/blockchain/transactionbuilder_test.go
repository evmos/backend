// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package blockchain

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/tharsis/dashboard-backend/internal/constants"
)

func TestCreateTransactionWithMessage(t *testing.T) {
	msgSendSdk, err := CreateMessageSend("evmos14uepnqnvkuyyvwe65wmncejq5g2f0tjft3wr65", "evmos1c7kfknyuamrkvgddg90g5tw2ncxy5g4wqlyxu6", sdk.NewInt(1), "aevmos", "evmos")
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
		"message",
		sdk.NewInt(3000000000000000),
		"aevmos",
		uint64(150000),
		pubKey,
		uint64(6),
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
		10, 137, 1, 10, 28, 47, 99, 111, 115, 109, 111, 115,
		46, 98, 97, 110, 107, 46, 118, 49, 98, 101, 116, 97,
		49, 46, 77, 115, 103, 83, 101, 110, 100, 18, 105, 10,
		44, 101, 118, 109, 111, 115, 49, 52, 117, 101, 112, 110,
		113, 110, 118, 107, 117, 121, 121, 118, 119, 101, 54, 53,
		119, 109, 110, 99, 101, 106, 113, 53, 103, 50, 102, 48,
		116, 106, 102, 116, 51, 119, 114, 54, 53, 18, 44, 101,
		118, 109, 111, 115, 49, 99, 55, 107, 102, 107, 110, 121,
		117, 97, 109, 114, 107, 118, 103, 100, 100, 103, 57, 48,
		103, 53, 116, 119, 50, 110, 99, 120, 121, 53, 103, 52,
		119, 113, 108, 121, 120, 117, 54, 26, 11, 10, 6, 97,
		101, 118, 109, 111, 115, 18, 1, 49, 18, 7, 109, 101,
		115, 115, 97, 103, 101,
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
		4, 10, 2, 8, 1, 24, 6, 18, 32, 10, 26, 10,
		6, 97, 101, 118, 109, 111, 115, 18, 16, 51, 48, 48,
		48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48,
		48, 16, 240, 147, 9,
	}

	for k, v := range authInfoBytesSignDirect {
		if v != authInfoBytesSignDirectFrontend[k] {
			t.Fatalf("Error Sign Direct: authInfoBytesSignDirect and authInfoBytesSignDirectFrontend are not equals")
		}
	}

	signBytesSignDirectFrontend := "R5Ok3UnAwRWqWRc0rCdSfPAT1hxv5OLa+5633xqHG+0="

	if tx.SignDirect.SignBytes != signBytesSignDirectFrontend {
		t.Fatalf("Error Sign Direct: tx.SignDirect.SignBytes and signBytesSignDirectFrontend are not equals")
	}

	bodyBytesLegacyAmino, err := proto.Marshal(&tx.LegacyAmino.Body)
	if err != nil {
		t.Fatalf("Error Legacy Amino: encoding body to []bytes from: %q", err)
	}

	bodyBtyesLegacyAminoFrontend := []byte{ //nolint:all
		10, 137, 1, 10, 28, 47, 99, 111, 115, 109, 111, 115,
		46, 98, 97, 110, 107, 46, 118, 49, 98, 101, 116, 97,
		49, 46, 77, 115, 103, 83, 101, 110, 100, 18, 105, 10,
		44, 101, 118, 109, 111, 115, 49, 52, 117, 101, 112, 110,
		113, 110, 118, 107, 117, 121, 121, 118, 119, 101, 54, 53,
		119, 109, 110, 99, 101, 106, 113, 53, 103, 50, 102, 48,
		116, 106, 102, 116, 51, 119, 114, 54, 53, 18, 44, 101,
		118, 109, 111, 115, 49, 99, 55, 107, 102, 107, 110, 121,
		117, 97, 109, 114, 107, 118, 103, 100, 100, 103, 57, 48,
		103, 53, 116, 119, 50, 110, 99, 120, 121, 53, 103, 52,
		119, 113, 108, 121, 120, 117, 54, 26, 11, 10, 6, 97,
		101, 118, 109, 111, 115, 18, 1, 49, 18, 7, 109, 101,
		115, 115, 97, 103, 101,
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
		4, 10, 2, 8, 127, 24, 6, 18, 32, 10, 26, 10,
		6, 97, 101, 118, 109, 111, 115, 18, 16, 51, 48, 48,
		48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48,
		48, 16, 240, 147, 9,
	}

	for k, v := range authInfoBytesLegacyAmino {
		if v != authInfoBytesLegacyAminoFrontend[k] {
			t.Fatalf("Error Legacy Amino: authInfoBytesLegacyAmino and authInfoBytesLegacyAminoFrontend are not equals")
		}
	}

	signBytesLegacyAminoFrontend := "4Ax/3aAn7dCKXSFuxM4vtNNN+zP+zXL5FSx5fh83A5U="

	if tx.LegacyAmino.SignBytes != signBytesLegacyAminoFrontend {
		t.Fatalf("Error Legacy Amino: tx.LegacyAmino.SignBytes and signBytesLegacyAminoFrontend are not equals")
	}

	// For the test I changed the chain_id in domain: 9001 for "0x2329"
	eipToSignFrontend := `{"types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"string"},{"name":"salt","type":"string"}],"Tx":[{"name":"account_number","type":"string"},{"name":"chain_id","type":"string"},{"name":"fee","type":"Fee"},{"name":"memo","type":"string"},{"name":"msgs","type":"Msg[]"},{"name":"sequence","type":"string"}],"Fee":[{"name":"feePayer","type":"string"},{"name":"amount","type":"Coin[]"},{"name":"gas","type":"string"}],"Coin":[{"name":"denom","type":"string"},{"name":"amount","type":"string"}],"Msg":[{"name":"type","type":"string"},{"name":"value","type":"MsgValue"}],"MsgValue":[{"name":"from_address","type":"string"},{"name":"to_address","type":"string"},{"name":"amount","type":"TypeAmount[]"}],"TypeAmount":[{"name":"denom","type":"string"},{"name":"amount","type":"string"}]},"primaryType":"Tx","domain":{"name":"Cosmos Web3","version":"1.0.0","chainId":"0x2329","verifyingContract":"cosmos","salt":"0"},"message":{"account_number":"2164290","chain_id":"evmos_9001-2","fee":{"amount":[{"amount":"3000000000000000","denom":"aevmos"}],"gas":"150000","feePayer":"evmos14uepnqnvkuyyvwe65wmncejq5g2f0tjft3wr65"},"memo":"message","msgs":[{"type":"cosmos-sdk/MsgSend","value":{"amount":[{"amount":"1","denom":"aevmos"}],"from_address":"evmos14uepnqnvkuyyvwe65wmncejq5g2f0tjft3wr65","to_address":"evmos1c7kfknyuamrkvgddg90g5tw2ncxy5g4wqlyxu6"}}],"sequence":"6"}}`

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

func TestSignature(t *testing.T) {
	signature := "0xb3a91fcc428c4bc4faecdcc4cf8e1dc012c8da5a4cc69d4d9898ea1978dbba8840d8437b74448b06e19f2f05a9170d8165e1c562f787d26accc7b9377733ebfb1b"
	expected := "b3a91fcc428c4bc4faecdcc4cf8e1dc012c8da5a4cc69d4d9898ea1978dbba8840d8437b74448b06e19f2f05a9170d8165e1c562f787d26accc7b9377733ebfb1b"
	result := Remove0xFromHex(signature)
	if result != expected {
		t.Fatalf("Error: couldn't remove 0x from signature")
	}

	result2 := Remove0xFromHex(expected)
	if result2 != expected {
		t.Fatalf("Error: signature without 0x not working")
	}
}

func TestCreateTransactionWithIBCTransfer(t *testing.T) {
	val, ok := constants.IBCChannels["EVMOS"]
	if !ok {
		t.Fatalf("Error trying to get ibc channel for evmos")
		return
	}

	srcChannel := ""
	for k, v := range val {
		if v == "OSMOSIS" {
			srcChannel = k
		}
	}
	if srcChannel == "" {
		t.Fatalf("Error trying to get ibc channel for evmos-osmosis")
		return
	}

	msg := CreateMsgTransfer("transfer", srcChannel, sdk.NewInt(1), "aevmos", "evmos14uepnqnvkuyyvwe65wmncejq5g2f0tjft3wr65", "osmo1j30xhsxcqss0n662wrma0vqw4zcx285munun8a", 1, 6641130, 9223372036854775808)

	pubKey, err := base64.StdEncoding.DecodeString("Ak8wUTcElcOofCZZJM97pduO+Aw3w4wzClrJgN2VzTVQ")
	if err != nil {
		t.Fatalf("Error decoding string represented by base64 to bytes")
	}

	tx, err := CreateTransactionWithMessage(
		EipToSignIBC{},
		[]sdk.Msg{msg},
		"message",
		sdk.NewInt(3000000000000000),
		"aevmos",
		uint64(150000),
		pubKey,
		uint64(9),
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
	//nolint:all
	bodyBytesSignDirectFrontend := []byte{
		10, 191, 1, 10, 41, 47, 105, 98, 99, 46, 97, 112,
		112, 108, 105, 99, 97, 116, 105, 111, 110, 115, 46, 116,
		114, 97, 110, 115, 102, 101, 114, 46, 118, 49, 46, 77,
		115, 103, 84, 114, 97, 110, 115, 102, 101, 114, 18, 145,
		1, 10, 8, 116, 114, 97, 110, 115, 102, 101, 114, 18,
		9, 99, 104, 97, 110, 110, 101, 108, 45, 48, 26, 11,
		10, 6, 97, 101, 118, 109, 111, 115, 18, 1, 49, 34,
		44, 101, 118, 109, 111, 115, 49, 52, 117, 101, 112, 110,
		113, 110, 118, 107, 117, 121, 121, 118, 119, 101, 54, 53,
		119, 109, 110, 99, 101, 106, 113, 53, 103, 50, 102, 48,
		116, 106, 102, 116, 51, 119, 114, 54, 53, 42, 43, 111,
		115, 109, 111, 49, 106, 51, 48, 120, 104, 115, 120, 99,
		113, 115, 115, 48, 110, 54, 54, 50, 119, 114, 109, 97,
		48, 118, 113, 119, 52, 122, 99, 120, 50, 56, 53, 109,
		117, 110, 117, 110, 56, 97, 50, 7, 8, 1, 16, 234,
		171, 149, 3, 56, 128, 128, 128, 128, 128, 128, 128, 128,
		128, 1, 18, 7, 109, 101, 115, 115, 97, 103, 101,
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
		4, 10, 2, 8, 1, 24, 9, 18, 32, 10, 26, 10,
		6, 97, 101, 118, 109, 111, 115, 18, 16, 51, 48, 48,
		48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48,
		48, 16, 240, 147, 9,
	}

	for k, v := range authInfoBytesSignDirect {
		if v != authInfoBytesSignDirectFrontend[k] {
			t.Fatalf("Error Sign Direct: authInfoBytesSignDirect and authInfoBytesSignDirectFrontend are not equals")
		}
	}

	signBytesSignDirectFrontend := "zQAD+0jF/QEcwj+wrXO9VUhOJOsFkaMgMXzDJvJBIWs="

	if tx.SignDirect.SignBytes != signBytesSignDirectFrontend {
		t.Fatalf("Error Sign Direct: tx.SignDirect.SignBytes and signBytesSignDirectFrontend are not equals")
	}

	bodyBytesLegacyAmino, err := proto.Marshal(&tx.LegacyAmino.Body)
	if err != nil {
		t.Fatalf("Error Legacy Amino: encoding body to []bytes from: %q", err)
	}

	bodyBtyesLegacyAminoFrontend := []byte{ //nolint:all
		10, 191, 1, 10, 41, 47, 105, 98, 99, 46, 97, 112,
		112, 108, 105, 99, 97, 116, 105, 111, 110, 115, 46, 116,
		114, 97, 110, 115, 102, 101, 114, 46, 118, 49, 46, 77,
		115, 103, 84, 114, 97, 110, 115, 102, 101, 114, 18, 145,
		1, 10, 8, 116, 114, 97, 110, 115, 102, 101, 114, 18,
		9, 99, 104, 97, 110, 110, 101, 108, 45, 48, 26, 11,
		10, 6, 97, 101, 118, 109, 111, 115, 18, 1, 49, 34,
		44, 101, 118, 109, 111, 115, 49, 52, 117, 101, 112, 110,
		113, 110, 118, 107, 117, 121, 121, 118, 119, 101, 54, 53,
		119, 109, 110, 99, 101, 106, 113, 53, 103, 50, 102, 48,
		116, 106, 102, 116, 51, 119, 114, 54, 53, 42, 43, 111,
		115, 109, 111, 49, 106, 51, 48, 120, 104, 115, 120, 99,
		113, 115, 115, 48, 110, 54, 54, 50, 119, 114, 109, 97,
		48, 118, 113, 119, 52, 122, 99, 120, 50, 56, 53, 109,
		117, 110, 117, 110, 56, 97, 50, 7, 8, 1, 16, 234,
		171, 149, 3, 56, 128, 128, 128, 128, 128, 128, 128, 128,
		128, 1, 18, 7, 109, 101, 115, 115, 97, 103, 101,
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
		4, 10, 2, 8, 127, 24, 9, 18, 32, 10, 26, 10,
		6, 97, 101, 118, 109, 111, 115, 18, 16, 51, 48, 48,
		48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48,
		48, 16, 240, 147, 9,
	}

	for k, v := range authInfoBytesLegacyAmino {
		if v != authInfoBytesLegacyAminoFrontend[k] {
			t.Fatalf("Error Legacy Amino: authInfoBytesLegacyAmino and authInfoBytesLegacyAminoFrontend are not equals")
		}
	}

	signBytesLegacyAminoFrontend := "IAd1jaU2AClmSW6rMwlhsFnZcH4ypxM5q237GtyMtio="

	if tx.LegacyAmino.SignBytes != signBytesLegacyAminoFrontend {
		t.Fatalf("Error Legacy Amino: tx.LegacyAmino.SignBytes and signBytesLegacyAminoFrontend are not equals")
	}

	// // For the test I changed the chain_id in domain: 9001 for "0x2329"
	eipToSignFrontend := `{"types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"string"},{"name":"salt","type":"string"}],"Tx":[{"name":"account_number","type":"string"},{"name":"chain_id","type":"string"},{"name":"fee","type":"Fee"},{"name":"memo","type":"string"},{"name":"msgs","type":"Msg[]"},{"name":"sequence","type":"string"}],"Fee":[{"name":"feePayer","type":"string"},{"name":"amount","type":"Coin[]"},{"name":"gas","type":"string"}],"Coin":[{"name":"denom","type":"string"},{"name":"amount","type":"string"}],"Msg":[{"name":"type","type":"string"},{"name":"value","type":"MsgValue"}],"MsgValue":[{"name":"source_port","type":"string"},{"name":"source_channel","type":"string"},{"name":"token","type":"TypeToken"},{"name":"sender","type":"string"},{"name":"receiver","type":"string"},{"name":"timeout_height","type":"TypeTimeoutHeight"},{"name":"timeout_timestamp","type":"uint64"}],"TypeToken":[{"name":"denom","type":"string"},{"name":"amount","type":"string"}],"TypeTimeoutHeight":[{"name":"revision_number","type":"uint64"},{"name":"revision_height","type":"uint64"}]},"primaryType":"Tx","domain":{"name":"Cosmos Web3","version":"1.0.0","chainId":"0x2329","verifyingContract":"cosmos","salt":"0"},"message":{"account_number":"2164290","chain_id":"evmos_9001-2","fee":{"amount":[{"amount":"3000000000000000","denom":"aevmos"}],"gas":"150000","feePayer":"evmos14uepnqnvkuyyvwe65wmncejq5g2f0tjft3wr65"},"memo":"message","msgs":[{"type":"cosmos-sdk/MsgTransfer","value":{"receiver":"osmo1j30xhsxcqss0n662wrma0vqw4zcx285munun8a","sender":"evmos14uepnqnvkuyyvwe65wmncejq5g2f0tjft3wr65","source_channel":"channel-0","source_port":"transfer","timeout_height":{"revision_height":"6641130","revision_number":"1"},"timeout_timestamp":"9223372036854775808","token":{"amount":"1","denom":"aevmos"}}}],"sequence":"9"}}`
	var resEipToSign EipToSignIBC
	err = json.Unmarshal([]byte(tx.EipToSign), &resEipToSign)
	if err != nil {
		t.Fatalf("Error: couldn't unmarshal tx.EipToSign: %q", err)
	}

	var resEipToSignFrontend EipToSignIBC
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
