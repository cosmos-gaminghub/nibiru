#!/bin/bash

NETWORK=upgrade-1
DAEMON=nibirud
HOME_DIR=~/.nibiru
CONFIG=~/.nibiru/config
TOKEN_DENOM=ugame

rm -rf $HOME_DIR

$DAEMON init $NETWORK --chain-id $NETWORK

$DAEMON config chain-id $NETWORK
$DAEMON config keyring-backend test

$DAEMON keys add eg --keyring-backend test
$DAEMON add-genesis-account $($DAEMON keys show eg -a --keyring-backend test) 100000000000000$TOKEN_DENOM

sed -i "s/\"stake\"/\"$TOKEN_DENOM\"/g" $CONFIG/genesis.json
jq '.app_state.gov.voting_params.voting_period = "60s"' $CONFIG/genesis.json > tmp.json && mv tmp.json $CONFIG/genesis.json

$DAEMON gentx eg 50000000000000$TOKEN_DENOM --commission-max-change-rate=0.1 --commission-max-rate=1 --commission-rate=0.1 --moniker=eg-validator --keyring-backend test --chain-id $NETWORK

$DAEMON collect-gentxs

$DAEMON validate-genesis

#timeout 10s $DAEMON start || ( [[ $? -eq 124 ]] && \
#echo "WARNING: Timeout reached, but that's OK" )
