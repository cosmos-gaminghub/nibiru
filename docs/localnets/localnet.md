# Start your own localnet

:::tip Required
[Install Nibiru](../install/install.md)
:::

## Initialize configuration files and genesis file


```sh
$ nibirud init <your_moniker> --chain-id testchain
```

**Copy the `Address` output here and save it for later use**

```sh
$ nibirud keys add eg
```


::: warning
If you get the following message, you need to specify `--keyring-backend=file`.

`No such interface “org.freedesktop.DBus.Properties” on object at path /`
:::


**Add account with tokens to the genesis file**.

```sh
$ nibirud add-genesis-account $(nibirud keys show eg -a) 100000000000000ugame
```

Default denom is `stake`, so if you want to customize the denom(ex:ugame), you have to edit `genesis.json` like below command.

```sh
$ sed -i "s/\"stake\"/\"ugame\"/g" ~/.nibiru/config/genesis.json
```

:::tip 
if you use Mac, then the command should be like this.

```sh
$ sed -i "" "s/\"stake\"/\"ugame\"/g" ~/.nibiru/config/genesis.json
```
:::

**Configure your CLI to eliminate need for chain-id flag**


Set defalut config value like `chain-id` to skip putting flag when broadcasting the transactions.
```sh
$ nibirud config chain-id testchain
$ nibirud config output json
```

Ready for start

```sh
# gentx is the create-validator command from genesis state, deciding how much token is self-delegated at the first place.
$ nibirud gentx eg 50000000000000ugame --chain-id=testchain --commission-max-change-rate=0.1 --commission-max-rate=1 --commission-rate=0.1 --moniker=eg-validator
$ nibirud collect-gentxs
$ nibirud validate-genesis
```

**Now let's start!**
```sh
$ nibirud start
```
