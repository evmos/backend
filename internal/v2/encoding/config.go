package encoding

import (
	"github.com/cosmos/cosmos-sdk/simapp/params"

	evmosapp "github.com/evmos/evmos/v12/app"
	"github.com/evmos/evmos/v12/encoding"
)

// MakeConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	mb := evmosapp.ModuleBasics
	return encoding.MakeConfig(mb)
}
