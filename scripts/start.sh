#!/bin/bash

NETWORK=nibiru-3000
DAEMON=nibirud
HOME_DIR=~/.nibiru
CONFIG=~/.nibiru/config
TOKEN_DENOM=ugame

rm -rf $HOME_DIR

$DAEMON init $NETWORK --chain-id $NETWORK
$DAEMON keys add eg --keyring-backend=test
$DAEMON add-genesis-account $($DAEMON keys show eg -a --keyring-backend=test) 100000000000000$TOKEN_DENOM

sed -i "s/\"stake\"/\"$TOKEN_DENOM\"/g" $HOME_DIR/config/genesis.json

$DAEMON gentx eg 50000000000000$TOKEN_DENOM --chain-id=$NETWORK --commission-max-change-rate=0.1 --commission-max-rate=1 --commission-rate=0.1 --moniker=eg-validator --keyring-backend=test


$DAEMON collect-gentxs

$DAEMON validate-genesis

timeout 10s $DAEMON start || ( [[ $? -eq 124 ]] && \
echo "WARNING: Timeout reached, but that's OK" )
