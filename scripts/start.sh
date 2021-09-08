#!/bin/bash
NETWORK=nibiru-2000
DAEMON=nibirud
HOME_DIR=~/.nibiru
CONFIG=~/.nibiru/config
TOKEN_DENOM=game

rm $CONFIG/genesis.json
rm -rf $CONFIG/gentx && mkdir $CONFIG/gentx

# Initialize configuration files and genesis file
$DAEMON init $NETWORK --chain-id $NETWORK

# fix denom in genesis.json
sed -i "s/\"stake\"/\"$TOKEN_DENOM\"/g" $HOME_DIR/config/genesis.json

# Copy the `Address` output here and save it for later use
$DAEMON keys add jack

# Add both accounts, with coins to the genesis file
$DAEMON add-genesis-account $($DAEMON keys show jack -a) 100000000000000$TOKEN_DENOM

$DAEMON gentx jack 10000000000000$TOKEN_DENOM --commission-rate=0.1 --commission-max-rate=1 --commission-max-change-rate=0.1 --pubkey $($DAEMON tendermint show-validator) --chain-id=$NETWORK

$DAEMON collect-gentxs

$DAEMON validate-genesis

timeout 10s $DAEMON start || ( [[ $? -eq 124 ]] && \
echo "WARNING: Timeout reached, but that's OK" )
