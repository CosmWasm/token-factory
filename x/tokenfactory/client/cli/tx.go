package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/token-factory/x/tokenfactory/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

var DeniedDenoms = [4]string{"ujuno", "juno", "ibc/", "factory/"}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewCreateDenomCmd(),
		NewMintCmd(),
		NewBurnCmd(),
		// NewForceTransferCmd(),
		NewChangeAdminCmd(),
		NewModifyDenomMetadataCmd(),
	)

	return cmd
}

// NewCreateDenomCmd broadcast MsgCreateDenom
func NewCreateDenomCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-denom [subdenom] [flags]",
		Short: "create a new denom from an account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			msg := types.NewMsgCreateDenom(
				clientCtx.GetFromAddress().String(),
				args[0],
			)

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewMintCmd broadcast MsgMint
func NewMintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [amount] [flags]",
		Short: "Mint a denom to an address. Must have admin authority to do so.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgMint(
				clientCtx.GetFromAddress().String(),
				amount,
			)

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewBurnCmd broadcast MsgBurn
func NewBurnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [amount] [flags]",
		Short: "Burn tokens from an address. Must have admin authority to do so.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgBurn(
				clientCtx.GetFromAddress().String(),
				amount,
			)

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// // NewForceTransferCmd broadcast MsgForceTransfer
// func NewForceTransferCmd() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "force-transfer [amount] [transfer-from-address] [transfer-to-address] [flags]",
// 		Short: "Force transfer tokens from one address to another address. Must have admin authority to do so.",
// 		Args:  cobra.ExactArgs(3),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx, err := client.GetClientTxContext(cmd)
// 			if err != nil {
// 				return err
// 			}

// 			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

// 			amount, err := sdk.ParseCoinNormalized(args[0])
// 			if err != nil {
// 				return err
// 			}

// 			msg := types.NewMsgForceTransfer(
// 				clientCtx.GetFromAddress().String(),
// 				amount,
// 				args[1],
// 				args[2],
// 			)

// 			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
// 		},
// 	}

// 	flags.AddTxFlagsToCmd(cmd)
// 	return cmd
// }

// NewChangeAdminCmd broadcast MsgChangeAdmin
func NewChangeAdminCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-admin [denom] [new-admin-address] [flags]",
		Short: "Changes the admin address for a factory-created denom. Must have admin authority to do so.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			msg := types.NewMsgChangeAdmin(
				clientCtx.GetFromAddress().String(),
				args[0],
				args[1],
			)

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewModifyDenomMetadataCmd broadcast a Bank Metadata modification transaction
func NewModifyDenomMetadataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify-metadata [denom] [pretty-name] [ticker-symbol] [description] [exponent] [flags]",
		Short: "Changes the base data for frontends to query the data of.",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			fullDenom, prettyName, ticker, desc := args[0], args[1], args[2], args[3]

			for _, prefix := range DeniedDenoms {
				if strings.HasPrefix(strings.ToLower(ticker), prefix) || strings.HasPrefix(strings.ToLower(prettyName), prefix) {
					return fmt.Errorf("ticker or prefix symbol which starts with: %s is not allowed", prefix)
				}
			}

			exponent, err := strconv.ParseUint(args[4], 10, 32)
			if err != nil {
				return err
			}

			bankMetadata := banktypes.Metadata{
				Description: desc,
				Display:     prettyName,
				Symbol:      ticker,
				Name:        prettyName,
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    fullDenom,
						Exponent: 0, // must be 0 for the base denom
						Aliases:  []string{prettyName},
					},
					{
						Denom:    prettyName,
						Exponent: uint32(exponent),
						Aliases:  []string{fullDenom},
					},
				},
				Base: fullDenom,
			}

			msg := types.NewMsgSetDenomMetadata(
				clientCtx.GetFromAddress().String(),
				bankMetadata,
			)

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
