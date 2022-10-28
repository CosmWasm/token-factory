package simulation

import (
	"math/rand"

	"github.com/CosmWasm/token-factory/app/params"
	"github.com/CosmWasm/token-factory/x/tokenfactory/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
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
	GetParams(ctx sdk.Context) (params types.Params)
	GetAuthorityMetadata(ctx sdk.Context, denom string) (types.DenomAuthorityMetadata, error)
	GetAllDenomsIterator(ctx sdk.Context) sdk.Iterator
	GetDenomsFromCreator(ctx sdk.Context, creator string) []string
}

type BankKeeper interface {
	simulation.BankKeeper
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
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
			SimulateMsgCreateDenom(
				tfKeeper,
				ak,
				bk,
			),
		),
	}
}

func SimulateMsgCreateDenom(tfKeeper TokenfactoryKeeper, ak types.AccountKeeper, bk BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand,
		app *baseapp.BaseApp,
		ctx sdk.Context,
		accs []simtypes.Account,
		chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		createFee := tfKeeper.GetParams(ctx).DenomCreationFee
		balances := bk.GetAllBalances(ctx, simAccount.Address)

		_, hasNeg := balances.SafeSub(createFee)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgCreateDenom{}.Type(), "Creator not enough creation fee"), nil, nil
		}

		msg := types.MsgCreateDenom{
			Sender:   simAccount.Address.String(),
			Subdenom: simtypes.RandStringOfLength(r, 10),
		}

		txCtx := BuildOperationInput(r, app, ctx, &msg, simAccount, ak, bk, createFee)
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// BuildOperationInput helper to build object
func BuildOperationInput(
	r *rand.Rand,
	app *baseapp.BaseApp,
	ctx sdk.Context,
	msg interface {
		sdk.Msg
		Type() string
	},
	simAccount simtypes.Account,
	ak types.AccountKeeper,
	bk BankKeeper,
	deposit sdk.Coins,
) simulation.OperationInput {
	return simulation.OperationInput{
		R:               r,
		App:             app,
		TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
		Cdc:             nil,
		Msg:             msg,
		MsgType:         msg.Type(),
		Context:         ctx,
		SimAccount:      simAccount,
		AccountKeeper:   ak,
		Bankkeeper:      bk,
		ModuleName:      types.ModuleName,
		CoinsSpentInMsg: deposit,
	}
}
