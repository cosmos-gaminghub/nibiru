#!/bin/bash

NETWORK=nibiru-3000
DAEMON=nibirud
HOME_DIR=~/.nibiru
CONFIG=~/.nibiru/config
TOKEN_DENOM=ugame

# upload
RES=$($DAEMON tx wasm store ./cw20_base.wasm --from=eg -y --keyring-backend=test --output=json --node=http://127.0.0.1:26657 --chain-id=$NETWORK --gas-prices=0.001$TOKEN_DENOM --gas=auto --gas-adjustment=1.3)
sleep 6
TX=$(echo $RES | jq -r '.txhash')
CODE_ID=$($DAEMON query tx $TX --output=json | jq -r '.logs[0].events[1].attributes[0].value')
echo "-> CODE_ID: $CODE_ID"
$DAEMON query wasm list-contract-by-code $CODE_ID --node=http://127.0.0.1:26657 --output=json

# init
MSG=$(jq -n "{\"name\":\"token-name\", \"symbol\":\"TKN\", \"decimals\": 18, \"initial_balances\": [], \"mint\": {\"minter\": \"$($DAEMON keys show eg -a --keyring-backend=test)\"}}")
$DAEMON tx wasm instantiate $CODE_ID "$MSG" --from=eg --label="contract1" --node=http://127.0.0.1:26657 --chain-id=$NETWORK --gas-prices=0.001$TOKEN_DENOM --gas=auto --gas-adjustment=1.3 -y --keyring-backend=test --admin $($DAEMON keys show eg --keyring-backend=test -a) --output=json
sleep 6
CONTRACT=$($DAEMON query wasm list-contract-by-code $CODE_ID --node=http://127.0.0.1:26657 --output=json | jq -r '.contracts[-1]')
echo "-> CONTRACT: $CONTRACT"

# execute
MSG=$(jq -n "{\"mint\":{\"recipient\": \"$($DAEMON keys show eg -a --keyring-backend=test)\", \"amount\": \"100\"}}")
$DAEMON tx wasm execute $CONTRACT "$MSG" --from=eg --node=http://127.0.0.1:26657 --chain-id=$NETWORK --gas-prices=0.001$TOKEN_DENOM --gas=auto --gas-adjustment=1.3 -y --keyring-backend=test --output=json
sleep 6

# query
MSG=$(jq -n "{\"balance\":{\"address\": \"$($DAEMON keys show eg -a --keyring-backend=test)\"}}")
$DAEMON query wasm contract-state smart $CONTRACT "$MSG"

# contract version
$DAEMON query wasm contract-state raw $CONTRACT 636F6E74726163745F696E666F --output=json | jq  -r .data | base64 -d | jq
