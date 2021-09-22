### Sample CLI Commands
```sh
// issue nft denom
nibirud tx nft issue nbr --name nibiru_coin --home ./.chaindata/ --chain-id mchain --from alice --keyring-backend test

// show all denoms
nibirud query nft denoms

// mint new nft
nibirud tx nft mint nbr --uri http://nibiru.com --home ./.chaindata/ --chain-id mchain --from alice --keyring-backend test

// show minted token
nibirud query nft token nbr 1
```
