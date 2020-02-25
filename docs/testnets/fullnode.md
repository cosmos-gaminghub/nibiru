
# Run Nibiru Full Node

:::tip Required
[Install Nibiru](/install/install)
:::

## Initialize node

```sh
nbrd init <your_moniker>
```

After that command, you can confirm that `.nbrd` folder is created in your home directory.

## Genesis file
Nibiru testnets genesis file is in [testnets repo](https://github.com/cosmos-gaminghub/testnets).
Download the latest genesis file by running the following command.
```sh
curl -o $HOME/.nbrd/config/genesis.json https://raw.githubusercontent.com/cosmos-gaminghub/testnets/master/latest/genesis.json
```


## Setup config
To connect other nodes running in the network, you have to set the seed nodes infomation in `config.toml`.
Open the `config.toml` file with vim editor, for example.

```sh
vim $HOME/.nbrd/config/config.toml
```

```config.toml
persistent_peers = "<node_id>@<node_ip_address>:<port>,<node_id>@<node_ip_address>:<port>"
```

## Run Full Node
```sh
nbrd start
```

It will take some time to sync with other node.
So be patient until your node find other available connections.


You can check sync status with the following command.

```sh
nbrcli status
```
