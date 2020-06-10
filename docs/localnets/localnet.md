# Start your own localnet

:::tip Required
[Install Nibiru](/install/install.md)
:::

## Initialize configuration files and genesis file


```sh
$ nbrd init <your_moniker> --chain-id testchain
```

**Copy the `Address` output here and save it for later use**

```sh
$ nbrcli keys add jack
```

**Add account with tokens to the genesis file**.

```sh
$ nbrd add-genesis-account $(nbrcli keys show jack -a) 100000000quark,100000000lepton
```

Default denom is `stake`, so if you want to customize the denom(ex:quark, lepton), you have to edit `genesis.json` like below command.

```
$ sed -i "s/\"stake\"/\"quark\"/g" ~/.nbrd/config/genesis.json
```

if you use Mac, then the command should be like this.

```
$ sed -i "" "s/\"stake\"/\"quark\"/g" ~/.nbrd/config/genesis.json
```

**Configure your CLI to eliminate need for chain-id flag**

```sh
$ nbrcli config chain-id testchain
$ nbrcli config output json
$ nbrcli config indent true
$ nbrcli config trust-node true
```

```sh
# gentx is the create-validator command from genesis state, deciding how much token is self-delegated at the first place.
$ nbrd gentx --amount 100000000quark --name jack
$ nbrd collect-gentxs
$ nbrd validate-genesis
```

**Now let's start!**
```sh
$ nbrd start
```
