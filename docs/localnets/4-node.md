# Start 4 Node

:::tip Required
- Docker is required to run 4 Node in your localnet.
:::

```sh
# Build the linux binary in ./build
$ make build-linux

# Build nbrchain/nbrdnode image
$ make build-docker
```


## Run your localnet with 4 Node
To start a 4 node testnet run:

```sh
make localnet-start
```

This command creates a 4-node network using the nbrdnode image. The ports for each node are found in this table:
| Node ID | P2P Port | RPC Port |
| --------|-------|------|
| `nbrnode0` | `26656` | `26657` |
| `nbrnode1` | `26659` | `26660` |
| `nbrnode2` | `26661` | `26662` |
| `nbrnode3` | `26663` | `26664` |

To update the binary, just rebuild it and restart the nodes:

```sh
$ make build-linux localnet-start
```

### Keys & Accounts

To interact with `nbrcli` and start querying state or creating txs, you use the
`nbrcli` directory of any given node as your `home`, for example:

```shell
$ nbrcli keys list --home ./build/node0/nbrcli
```
