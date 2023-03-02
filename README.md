# Token Factory
The token-factory is a [Cosmos-Sdk](https://github.com/cosmos/cosmos-sdk) [module](https://docs.cosmos.network/main/building-modules/intro) and extension to the [wasmd](https://github.com/CosmWasm/wasmd) smart contract platform module.

The code in this repository is copied from https://github.com/osmosis-labs/osmosis/tree/main/x/tokenfactory. 
Credits and big thank you go to the original authors!

### Development
The project is setup as a Go [multi-module workspace](https://go.dev/doc/tutorial/workspaces).
* `x/` is the cosmos module code
* `demo/` contains a demo app to showcase the integration and used for e2e testing
* `e2e/` end-to-end tests


`go work sync` syncronize workspace dependencies