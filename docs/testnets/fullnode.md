# Run Nibiru Full Node

:::tip Required
[Install Nibiru](../install/install.md)

:::

## Join from the genesis (gentx submission)
1. Initialize the nibiru directory and create the local file with the correct chain-id

```
nibirud init <moniker> --chain-id=nibiru-3000
```

2. Create a local key pair in the keybase
```
nibirud keys add <your key name>
```

3. Add tour account to your local genesis file with a given amount and key you just created.
```
nibirud add-genesis-account $(nibirud keys show <your key name> -a) 100000000000ugame
```

4. Create the gentx
```
nibirud gentx <your key name> 100000000000ugame --commission-rate=0.1 --commission-max-rate=1 --commission-max-change-rate=0.1 --pubkey $(nibirud tendermint show-validator) --chain-id=nibiru-3000
```

5. Create Pull Request to [the testnet repo](https://github.com/cosmos-gaminghub/testnets).


## Join from the middle of the testnet

### Initialize node

```sh
nibirud init <your_moniker> --chain-id=nibiru-3000
```

After that command, you can confirm that `.nibiru` folder is created in your home directory.

### Get Genesis file
Nibiru testnets genesis file is in [testnets repo](https://github.com/cosmos-gaminghub/testnets).
Download the latest genesis file by running the following command.
```sh
curl -o $HOME/.nibiru/config/genesis.json https://raw.githubusercontent.com/cosmos-gaminghub/testnets/master/latest/genesis.json
```


### Setup config
To connect other nodes running in the network, you have to set the seed nodes infomation in `config.toml`.
Open the `config.toml` file with vim editor, for example.

```sh
vim $HOME/.nibiru/config/config.toml
```

```config.toml
persistent_peers = "<node_id>@<node_ip_address>:<port>,<node_id>@<node_ip_address>:<port>"
```

### Run Full Node
```sh
nibirud start
```

It will take some time to sync with other node.
So be patient until your node find other available connections.


You can check sync status with the following command.

```sh
nibirud status
```
