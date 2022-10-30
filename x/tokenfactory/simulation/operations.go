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
		weightMsgMint        int
		// weightMsgBurn             int
		// weightMsgChangeAdmin      int
		// weightMsgSetDenomMetadata int
	)

	simstate.AppParams.GetOrGenerate(simstate.Cdc, OpWeightMsgCreateDenom, &weightMsgCreateDenom, nil,
		func(_ *rand.Rand) {
			weightMsgCreateDenom = params.DefaultWeightMsgCreateDenom
		},
	)
	simstate.AppParams.GetOrGenerate(simstate.Cdc, OpWeightMsgMint, &weightMsgMint, nil,
		func(_ *rand.Rand) {
			weightMsgMint = params.DefaultWeightMsgMint
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
		simulation.NewWeightedOperation(
			weightMsgMint,
			SimulateMsgMint(
				tfKeeper,
				ak,
				bk,
				DefaultSimulationDenomSelector,
			),
		),
	}
}

type MsgMintDenomSelector = func(*rand.Rand, sdk.Context, TokenfactoryKeeper, string) (string, bool)

func DefaultSimulationDenomSelector(r *rand.Rand, ctx sdk.Context, tfKeeper TokenfactoryKeeper, creator string) (string, bool) {
	denoms := tfKeeper.GetDenomsFromCreator(ctx, creator)
	if len(denoms) == 0 {
		return "", false
	}
	randPos := r.Intn(len(denoms))

	return denoms[randPos], true
}

// Simulate msg mint denom
func SimulateMsgMint(
	tfKeeper TokenfactoryKeeper,
	ak types.AccountKeeper,
	bk BankKeeper,
	denomSelector MsgMintDenomSelector,
) simtypes.Operation {
	return func(
		r *rand.Rand,
		app *baseapp.BaseApp,
		ctx sdk.Context,
		accs []simtypes.Account,
		chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Get sims account
		simAccount, _ := simtypes.RandomAcc(r, accs)

		// Get demon created from sims account
		denom, hasDenom := denomSelector(r, ctx, tfKeeper, simAccount.Address.String())
		if !hasDenom {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgMint{}.Type(), "sim account have no denom created"), nil, nil
		}

		// Rand mint amount
		mintAmount, _ := simtypes.RandPositiveInt(r, sdk.NewIntFromUint64(100_000_000))

		// Create msg mint
		msg := types.MsgMint{
			Sender: simAccount.Address.String(),
			Amount: sdk.NewCoin(denom, mintAmount),
		}

		txCtx := BuildOperationInput(r, app, ctx, &msg, simAccount, ak, bk, nil)
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// Simulate msg create denom
func SimulateMsgCreateDenom(tfKeeper TokenfactoryKeeper, ak types.AccountKeeper, bk BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand,
		app *baseapp.BaseApp,
		ctx sdk.Context,
		accs []simtypes.Account,
		chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Get sims account
		simAccount, _ := simtypes.RandomAcc(r, accs)

		// Check if sims account enough create fee
		createFee := tfKeeper.GetParams(ctx).DenomCreationFee
		balances := bk.GetAllBalances(ctx, simAccount.Address)
		_, hasNeg := balances.SafeSub(createFee)
		if hasNeg {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgCreateDenom{}.Type(), "Creator not enough creation fee"), nil, nil
		}

		// Create msg create denom
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
