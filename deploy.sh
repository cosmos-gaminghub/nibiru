#!/bin/bash

export APP_HOME=./.chaindata/
export RPC=http://127.0.0.1:26657
export CHAIN_ID=localnet
export NODE=(--node $RPC)
export TXFLAG=($NODE --chain-id $CHAIN_ID --gas-prices 0.001ugtn --gas auto --gas-adjustment 1.3)

// upload
nibirud query wasm list-code $NODE
export RES=$(nibirud tx wasm store .artifacts/nft.wasm --from alice $TXFLAG -y --home ${APP_HOME})
# export CODE_ID=$(echo $RES | jq -r '.logs[0].events[1].attributes[0].value')
# nibirud query wasm list-contract-by-code $CODE_ID $NODE --output json

# // init
# nibirud tx wasm instantiate $CODE_ID "{}" --from alice --label "nft 1" $TXFLAG -y --home ${APP_HOME}
# nibirud query wasm list-contract-by-code $CODE_ID $NODE --output json
# export CONTRACT=$(nibirud query wasm list-contract-by-code $CODE_ID $NODE --output json | jq -r '.contracts[-1]')
# nibirud query wasm contract $CONTRACT $NODE
# nibirud query bank balances $CONTRACT $NODE
# nibirud query wasm contract-state all $CONTRACT $NODE
# nibirud query wasm contract-state raw $CONTRACT 0006636f6e666967 $NODE --hex

# // execute
# export ISSUE=$(jq -n '{"issue_denom":{"denom_id":"test-denom-id","name":"test-name","schema":"test-schema"}}')
# nibirud tx wasm execute $CONTRACT "$ISSUE" \
#     --from alice $TXFLAG -y --home ${APP_HOME}



