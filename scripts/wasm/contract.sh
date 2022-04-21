#!/bin/bash

NETWORK=nibiru-3000
DAEMON=nibirud
HOME_DIR=~/.nibiru
CONFIG=~/.nibiru/config
TOKEN_DENOM=ugame

# upload
RES=$($DAEMON tx wasm store ./queue.wasm --from=eg -y --keyring-backend=test --output=json --node=http://127.0.0.1:26657 --chain-id=$NETWORK --gas-prices=0.001$TOKEN_DENOM --gas=auto --gas-adjustment=1.3)
sleep 5
TX=$(echo $RES | jq -r '.txhash')
CODE_ID=$($DAEMON query tx $TX --output=json | jq -r '.logs[0].events[1].attributes[0].value')
echo "-> CODE_ID: $CODE_ID"
$DAEMON query wasm list-contract-by-code $CODE_ID --node=http://127.0.0.1:26657 --output=json

# init
$DAEMON tx wasm instantiate $CODE_ID "{}" --from=eg --label="contract1" --node=http://127.0.0.1:26657 --chain-id=$NETWORK --gas-prices=0.001$TOKEN_DENOM --gas=auto --gas-adjustment=1.3 -y --keyring-backend=test --admin $($DAEMON keys show eg --keyring-backend=test -a) --output=json
sleep 5
CONTRACT=$($DAEMON query wasm list-contract-by-code $CODE_ID --node=http://127.0.0.1:26657 --output=json | jq -r '.contracts[-1]')
echo "-> CONTRACT: $CONTRACT"

# execute
MSG=$(jq -n '{"enqueue":{"value":123}}')
$DAEMON tx wasm execute $CONTRACT "$MSG" --from=eg --node=http://127.0.0.1:26657 --chain-id=$NETWORK --gas-prices=0.001$TOKEN_DENOM --gas=auto --gas-adjustment=1.3 -y --keyring-backend=test --output=json
sleep 5

# query
MSG=$(jq -n '{"count":{}}')
$DAEMON query wasm contract-state smart $CONTRACT "$MSG"
MSG=$(jq -n '{"sum":{}}')
$DAEMON query wasm contract-state smart $CONTRACT "$MSG"
