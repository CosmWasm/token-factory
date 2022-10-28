package simulation

import (
	"math/rand"

	"github.com/CosmWasm/token-factory/app/params"
	"github.com/CosmWasm/token-factory/x/tokenfactory/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// Simulation operation weights constants
//
//nolint:gosec
const (
	OpWeightMsgCreateDenom      = "op_weight_msg_create_denom"
	OpWeightMsgMint             = "op_weight_msg_mint"
	OpWeightMsgBurn             = "op_weight_msg_burn"
	OpWeightMsgChangeAdmin      = "op_weight_msg_change_admin"
	OpWeightMsgSetDenomMetadata = "op_weight_msg_set_denom_metadata"
)

type TokenfactoryKeeper interface {
}

type BankKeeper interface {
}

func WeightedOperations(
	simstate *module.SimulationState,
	tfKeeper TokenfactoryKeeper,
	ak types.AccountKeeper,
	bk BankKeeper,
) simulation.WeightedOperations {

	var (
		weightMsgCreateDenom int
		// weightMsgMint             int
		// weightMsgBurn             int
		// weightMsgChangeAdmin      int
		// weightMsgSetDenomMetadata int
	)

	simstate.AppParams.GetOrGenerate(simstate.Cdc, OpWeightMsgCreateDenom, &weightMsgCreateDenom, nil,
		func(_ *rand.Rand) {
			weightMsgCreateDenom = params.DefaultWeightMsgCreateDenom
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreateDenom,
			SimulateMsgCreateDenom(),
		),
	}
}

func SimulateMsgCreateDenom() simtypes.Operation {
	return func(
		r *rand.Rand,
		app *baseapp.BaseApp,
		ctx sdk.Context,
		accs []simtypes.Account,
		chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

	}
}
