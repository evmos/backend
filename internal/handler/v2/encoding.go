package v2

import (
	amino "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"

	"github.com/evmos/evmos/v12/app"
	enccodec "github.com/evmos/evmos/v12/encoding/codec"
)

// MakeConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	mb := app.ModuleBasics
	cdc := amino.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	codec := amino.NewProtoCodec(interfaceRegistry)

	encodingConfig := params.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          tx.NewTxConfig(codec, tx.DefaultSignModes),
		Amino:             cdc,
	}

	enccodec.RegisterLegacyAminoCodec(encodingConfig.Amino)
	mb.RegisterLegacyAminoCodec(encodingConfig.Amino)
	enccodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	mb.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
