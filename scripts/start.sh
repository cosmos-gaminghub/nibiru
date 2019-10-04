rm -rf ~/.nbrd
rm -rf ~/.nbrcli
# Initialize configuration files and genesis file
nbrd init eguegu --chain-id testchain

# Copy the `Address` output here and save it for later use
nbrcli keys add jack

# Copy the `Address` output here and save it for later use
nbrcli keys add alice

# Add both accounts, with coins to the genesis file
nbrd add-genesis-account $(nbrcli keys show jack -a) 100000000nbr,100000000stake
nbrd add-genesis-account $(nbrcli keys show alice -a) 100000000nbr,100000000stake

# Configure your CLI to eliminate need for chain-id flag
nbrcli config chain-id testchain
nbrcli config output json
nbrcli config indent true
nbrcli config trust-node true

nbrd gentx --name jack

nbrd collect-gentxs

nbrd validate-genesis
