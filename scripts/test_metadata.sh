export KEY="tf1" # wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g
export CHAIN_ID=${CHAIN_ID:-"tf-1"}
export KEYRING=${KEYRING:-"test"}
export HOME_DIR=$(eval echo "${HOME_DIR:-"~/.tf1/"}")
export BINARY=${BINARY:-toked}

alias BINARY="$BINARY --home=$HOME_DIR"

FLAGS="--from $KEY --keyring-backend=test --gas=1000000 --node http://localhost:26657 --chain-id tf-1 --yes --broadcast-mode=block"
FLAGS2="--from feeacc --keyring-backend=test --gas=1000000 --node http://localhost:26657 --chain-id tf-1 --yes --broadcast-mode=block"

# toked keys list --keyring-backend test --home $HOME_DIR
DENOM=factory/wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g/test
DENOM2=factory/wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g/test2
DENOM_NOT_MINE=factory/wasm1efd63aw40lxf3n4mhf7dzhjkr453axursysrvp/notmine

BINARY tx tokenfactory create-denom test $FLAGS
BINARY tx tokenfactory create-denom test2 $FLAGS
BINARY tx tokenfactory create-denom notmine $FLAGS2

BINARY query tokenfactory denoms-from-creator wasm1hj5fveer5cjtn4wd6wstzugjfdxzl0xpvsr89g --node http://localhost:26657

# set the metadata for it
BINARY tx tokenfactory modify-metadata $DENOM "ticker" "some desc https://www.com" 6 $FLAGS
BINARY tx tokenfactory modify-metadata $DENOM "ticker" "" 18 $FLAGS
BINARY q bank denom-metadata --denom $DENOM --node http://localhost:26657

# fails
BINARY tx tokenfactory modify-metadata $DENOM2 "JUNO" "desc" 6 $FLAGS
BINARY tx tokenfactory modify-metadata $DENOM2 "!" "desc" 6 $FLAGS
BINARY tx tokenfactory modify-metadata $DENOM2 "LJUNO" "desc" 6 $FLAGS
BINARY tx tokenfactory modify-metadata $DENOM2 "LUJUNO" "desc" 6 $FLAGS
BINARY tx tokenfactory modify-metadata $DENOM2 "sla/sh" "desc" 6 $FLAGS
BINARY tx tokenfactory modify-metadata $DENOM2 "thelengthiswaytoolonghere" "desc" 1 $FLAGS
BINARY tx tokenfactory modify-metadata $DENOM2 "" "desc" 1 $FLAGS
BINARY tx tokenfactory modify-metadata $DENOM2 "expont" "desc" 19 $FLAGS

# can't edit a token we don't own (tries to edit feeacc's address)
BINARY tx tokenfactory modify-metadata $DENOM_NOT_MINE "notmy" "desc" 1 $FLAGS

# query data

BINARY q bank denom-metadata --denom $DENOM2 --node http://localhost:26657

BINARY tx tokenfactory mint 100$DENOM $FLAGS