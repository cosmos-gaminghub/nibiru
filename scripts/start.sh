rm -rf ~/.nibiru
# Initialize configuration files and genesis file
nibirud init eguegu --chain-id testchain

# Copy the `Address` output here and save it for later use
nibirud keys add jack

# Copy the `Address` output here and save it for later use
nibirud keys add alice

# Add both accounts, with coins to the genesis file
nibirud add-genesis-account $(nibirud keys show jack -a) 100000000000000nbr,100000000000000quark
nibirud add-genesis-account $(nibirud keys show alice -a) 100000000000000nbr,100000000000000quark

# Configure your CLI to eliminate need for chain-id flag
#nibirud config chain-id testchain
#nibirud config output json
#nibirud config indent true
#nibirud config trust-node true

nibirud gentx jack 50000000000000quark --chain-id=testchain --commission-max-change-rate=0.1 --commission-max-rate=1 --commission-rate=0.1 --moniker=jack-validator

nibirud collect-gentxs

nibirud validate-genesis
