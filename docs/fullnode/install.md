# Install and Setup Nibiru fullnode

## Install go

:::tip Required
**Go 1.13.0+** is required for the Cosmos SDK.
:::

Firstly, install `golang` from [the official golang donwload page](https://golang.org/dl/).
Be sure to set your `$GOPATH`, `$GOBIN`, and `$PATH` environment variables, for example:

```sh
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bashrc
echo "export GOBIN=$GOPATH/bin" >> ~/.bashrc
echo "export PATH=$PATH:$GOBIN" >> ~/.bashrc
echo "export GO111MODULE=on" >> ~/.bashrc
source ~/.bashrc
```

Verify that `golang` has been installed successfully.

```sh
$ go version
go version go1.13.7 linux/amd64
```


## Install Nibiru
With `golang`, you can compile and run `nibiru`.

```sh
git clone github.com/cosmos-gaminghub/nibiru.git
cd nibiru && git checkout master
make install
```

Try `nbrcli version` and `nbrd version` to verify that everything is fine.

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
