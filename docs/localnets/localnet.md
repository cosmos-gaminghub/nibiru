# Start your own localnet

:::tip Required
[Install Nibiru](/install/install.md)
:::

## Initialize configuration files and genesis file


```sh
$ nibirud init <your_moniker> --chain-id testchain
```

**Copy the `Address` output here and save it for later use**

```sh
$ nibirud keys add jack
```


::: warning
If you get the following message, you need to specify `--keyring-backend=file`.

`No such interface “org.freedesktop.DBus.Properties” on object at path /`
:::




**Add account with tokens to the genesis file**.

```sh
$ nibirud add-genesis-account $(nibirud keys show jack -a) 100000000000000game
```

Default denom is `stake`, so if you want to customize the denom(ex:game), you have to edit `genesis.json` like below command.

```
$ sed -i "s/\"stake\"/\"game\"/g" ~/.nibiru/config/genesis.json
```

if you use Mac, then the command should be like this.

```
$ sed -i "" "s/\"stake\"/\"game\"/g" ~/.nibiru/config/genesis.json
```

**Configure your CLI to eliminate need for chain-id flag**

::: warning
**v0.3** is not ready this command
:::


```sh
$ nibirud config chain-id testchain
$ nibirud config output json
$ nibirud config indent true
$ nibirud config trust-node true
```

Ready for start

```sh
# gentx is the create-validator command from genesis state, deciding how much token is self-delegated at the first place.
$ nibirud gentx jack 50000000000000game --chain-id=testchain --commission-max-change-rate=0.1 --commission-max-rate=1 --commission-rate=0.1 --moniker=jack-validator
$ nibirud collect-gentxs
$ nibirud validate-genesis
```

**Now let's start!**
```sh
$ nibirud start
```
