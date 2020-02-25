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

**Add account with coins to the genesis file**

```sh
$ nbrd add-genesis-account $(nbrcli keys show jack -a) 100000000quark,100000000lepton
```

**Configure your CLI to eliminate need for chain-id flag**

```sh
$ nbrcli config chain-id testchain
$ nbrcli config output json
$ nbrcli config indent true
$ nbrcli config trust-node true
```

```sh
$ nbrd gentx --name jack
$ nbrd collect-gentxs
$ nbrd validate-genesis
```

**Now let's start!**
```sh
$ nbrd start
```
