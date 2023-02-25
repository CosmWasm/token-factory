export KEY="tf1" # wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g
export CHAIN_ID=${CHAIN_ID:-"tf-1"}
export KEYRING=${KEYRING:-"test"}
export HOME_DIR=$(eval echo "${HOME_DIR:-"~/.tf1/"}")
export BINARY=${BINARY:-toked}

alias BINARY="$BINARY --home=$HOME_DIR"

FLAGS="--from $KEY --keyring-backend=test --gas=1000000 --node http://localhost:26657 --chain-id tf-1 --yes --broadcast-mode=block"

DENOM=test
DENOM2=test2
BINARY tx tokenfactory create-denom $DENOM $FLAGS
BINARY tx tokenfactory create-denom $DENOM2 $FLAGS

BINARY query tokenfactory denoms-from-creator wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g --node http://localhost:26657

# set the metadata for it
BINARY tx tokenfactory modify-metadata factory/wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g/$DENOM "udenom" "TICKER" "some desc https://www.com" 6 $FLAGS

# fails
BINARY tx tokenfactory modify-metadata factory/wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g/$DENOM2 "ujuno" "TICKER" "desc" 6 $FLAGS
BINARY tx tokenfactory modify-metadata factory/wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g/$DENOM2 "juno" "TICKER" "desc" 6 $FLAGS
BINARY tx tokenfactory modify-metadata factory/wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g/$DENOM2 "ibc/test" "TICKER" "desc" 1 $FLAGS
BINARY tx tokenfactory modify-metadata factory/wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g/$DENOM2 "factory/" "TICKER" "desc" 1 $FLAGS

# query data
BINARY q bank denom-metadata --denom factory/wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g/$DENOM --node http://localhost:26657
BINARY q bank denom-metadata --denom factory/wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g/$DENOM2 --node http://localhost:26657

BINARY tx tokenfactory mint 100factory/wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g/$DENOM $FLAGS