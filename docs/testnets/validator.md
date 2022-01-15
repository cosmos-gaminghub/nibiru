# Become A Validator

:::tip Required
[Run Nibiru Full Node](../testnets/fullnode.md)
:::

## Create a Wallet
Firstly, you neeed to create wallet or import existing wallet, then get some $GAME from the faucet.

```
nibirud keys add <wallet_name>
```

:::warning Back up seed words
Be sure to back up the seed phrase in a safe way.
You need seeds when you recover your account when you lose your account password.
:::

## Check your node is sync
```sh
nibirud status |jq .sync_info
```
Check if your node has the same `latest_block_height`.


## Create Validator
Check your `nibiruvalconspub` address which is used to create a new validator.

```sh
nibirud tendermint show-validator
```

If your node is sync fully, then you can run the following command to upgrade your node to be a validator.

:::tip testnet faucet
You need some $GAME as Gus to send a tx.
Please ask in [the Discord validator channel](discord.gg/VfvTCP7Rm2) for faucet token with your address starting `nibiru1`.
:::

```sh
nibirud tx staking create-validator \
  --amount=1000game \
  --pubkey=$(nibirud tendermint show-validator) \
  --moniker=<your_validator_name> \
  --chain-id=<chain_id> \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas="auto" \
  --from=<wallet_name>
```

:::tip Current testnet chain-id
`chain-id` is `nibiru-2000`
:::
:::tip Insufficient Gas
If you got an error that implys insufficient gas, then you can modify `--gas="auto"` to some appropriate value.
:::


### Block Explorer
Check [GAME Explorer](https://nibiru-2000.game-explorer.io/) to see if your node in validator set correctly.


## Edit Validator
You can edit your validator metadata with the following command.


```sh
nibirud tx staking edit-validator \
  --moniker=<your_validator_name> \
  --website=<your_website> \
  --identity=<your_keybase_identity> \
  --details=<some_description> \
  --chain-id=<chain_id> \
  --from=<wallet_name> \
  --commission-rate="0.11"
```

Param| Description
--------- | ---------
moniker | your full node moniker is default
website | your website
identity | you can get 16 degit string from [keybase.io](https://keybase.io/)
details| some description
chain-id| you can check latest testnets in [testnets repo](https://github.com/cosmos-gaminghub/testnets)
from| your wallet name
commission-rate| refer to the tip below

:::tip Commission Rate
`commission-rate` can be changed within `commission-max-change-rate` per day and upper limit is `commission-max-rate`.
:::
