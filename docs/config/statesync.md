# State Sync

Nibiru has the state sync feature which allows a new node to join a network by fetching a snapshot of the application state at a recent height. This can reduce the time needed to sync with the network from days to minutes.

::: tip
Check the detail about state sync protocol in [the Cosmos Blog post](https://blog.cosmos.network/cosmos-sdk-state-sync-guide-99e4cf43be2f).
:::

::: warning
A state synced node will only restore the application state for the height the snapshot was taken at, and will not contain historical data nor historical blocks.
:::

## Snapshot Node
There should be a snaphot node in the network. That node should start process with the command below.

```sh
nibirud start --state-sync.snapshot-interval 100 --state-sync.snapshot-keep-recent 2
```

## State Sync Node
If you want to use state sync feature, then you need to set some parameters in `config.toml` as below.

```toml
#######################################################
###           P2P Configuration Options             ###
#######################################################
[p2p]
seeds = "node-id@rpcendpoint-1:26656"

...

#######################################################
###         State Sync Configuration Options        ###
#######################################################
[statesync]
enable = true
rpc_servers = "rpcendpoint-1:26657,rpcendpoint-2:26657"
trust_height = 1819
trust_hash = "2F7B6108627F7BF888E1D3C4ED649337D43DF8C1C8D957808469845E7475A995"
trust_period = "336h"
```

To get `trust_height` and `trust_hash`, you can use the command below.

```sh
curl -s http://rpcendpoint-1:26657/block | jq -r '.result.block.header.height + "\n" + .result.block_id.hash'
```

After the config is set, you can start node with `nibirud start` as nomally. Then your node will find snapshot data from the network and use it for state sync.

```sh
INF Discovering snapshots for 15s module=statesync
INF Discovered new snapshot format=1 hash="f|~R\u009b \x1e�a<�>���7���G�\toƯ���\x1fNA" height=1800 module=statesync
INF VerifyHeader hash=FC9E9C22045C56E0C9464F7AA08CA05BDAA50BA7F64348A3C1688779EAF38CB9 height=1801 module=light
INF VerifyHeader hash=5B33B4D3AFF4DA0E3F4731D983EE88C07EFF0807FA8979EB910D5FF0F2437798 height=1802 module=light
INF Offering snapshot to ABCI app format=1 hash="f|~R\u009b \x1e�a<�>���7���G�\toƯ���\x1fNA" height=1800 module=statesync
INF Snapshot accepted, restoring format=1 hash="f|~R\u009b \x1e�a<�>���7���G�\toƯ���\x1fNA" height=1800 module=statesync
INF Fetching snapshot chunk chunk=0 format=1 height=1800 module=statesync total=5
INF Fetching snapshot chunk chunk=2 format=1 height=1800 module=statesync total=5
INF Fetching snapshot chunk chunk=3 format=1 height=1800 module=statesync total=5
INF Fetching snapshot chunk chunk=1 format=1 height=1800 module=statesync total=5
INF VerifyHeader hash=5F4572EEC0142EE781303970E4F856EC4A2113A0CE0E236CFC0CCE26FFED57DD height=1800 module=light
INF Header has already been verified hash=FC9E9C22045C56E0C9464F7AA08CA05BDAA50BA7F64348A3C1688779EAF38CB9 height=1801 module=light
INF Header has already been verified hash=5B33B4D3AFF4DA0E3F4731D983EE88C07EFF0807FA8979EB910D5FF0F2437798 height=1802 module=light
INF Header has already been verified hash=FC9E9C22045C56E0C9464F7AA08CA05BDAA50BA7F64348A3C1688779EAF38CB9 height=1801 module=light
INF Header has already been verified hash=5F4572EEC0142EE781303970E4F856EC4A2113A0CE0E236CFC0CCE26FFED57DD height=1800 module=light
INF Fetching snapshot chunk chunk=4 format=1 height=1800 module=statesync total=5
INF Applied snapshot chunk to ABCI app chunk=0 format=1 height=1800 module=statesync total=5
INF Applied snapshot chunk to ABCI app chunk=1 format=1 height=1800 module=statesync total=5
INF Applied snapshot chunk to ABCI app chunk=2 format=1 height=1800 module=statesync total=5
INF Applied snapshot chunk to ABCI app chunk=3 format=1 height=1800 module=statesync total=5
INF Applied snapshot chunk to ABCI app chunk=4 format=1 height=1800 module=statesync total=5
INF Verified ABCI app appHash="'\f4���XB\uf159�\x03}��\u007f\x12�\x1a\x06���/�+\"�\x19q\x1d" height=1800 module=statesync
INF Snapshot restored format=1 hash="f|~R\u009b \x1e�a<�>���7���G�\toƯ���\x1fNA" height=1800 module=statesync
INF Starting BlockPool service impl=BlockPool module=blockchain
INF minted coins from module account amount=4283880ugame from=mint module=x/bank
INF executed block height=1801 module=state num_invalid_txs=0 num_valid_txs=12
INF commit synced commit=436F6D6D697449447B5B3439203231392037352031383420363220313533203136382031323120313131203934203131203220323333203230372031373920313231203230352031353320313732203620313337203932203339203933203320313520313332203830203232382031353720313239203134305D3A3730397D
INF committed state app_hash=31DB4BB83E99A8796F5E0B02E9CFB379CD99AC06895C275D030F8450E49D818C height=1801 module=state num_txs=12
```


:::tip
If you fail to use state sync feature, check if the snapshot node and state sync node are same version.
Also if `trust_period` is too short, It showing error message "Can't verify ..."
:::
