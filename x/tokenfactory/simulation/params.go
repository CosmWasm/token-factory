package simulation

import (
	"fmt"
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/CosmWasm/token-factory/x/tokenfactory/types"
)

// Simulation parameter constants
const (
	DefaultWeightMsgCreateDenom      int = 100
	DefaultWeightMsgMint             int = 100
	DefaultWeightMsgBurn             int = 100
	DefaultWeightMsgChangeAdmin      int = 100
	DefaultWeightMsgSetDenomMetadata int = 100
)

func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(
			types.ModuleName,
			string(types.KeyDenomCreationFee),
			func(r *rand.Rand) string {
				amount := RandDenomCreationFeeParam(r)
				return fmt.Sprintf("[{\"denom\":\"%v\",\"amount\":\"%v\"}]", amount[0].Denom, amount[0].Amount)
			},
		),
	}
}
